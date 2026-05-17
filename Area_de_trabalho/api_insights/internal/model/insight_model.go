package model

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

type InsightModel struct {
	repo InsightRepository
}

func NewInsightModel(repo InsightRepository) *InsightModel {
	return &InsightModel{repo: repo}
}

func (m *InsightModel) GenerateReport(ctx context.Context) (InsightReport, error) {
	data, err := m.loadAndNormalize(ctx)
	if err != nil {
		return InsightReport{}, err
	}

	var wg sync.WaitGroup
	salesCh := make(chan salesInsight, 1)
	suggestionsCh := make(chan map[string]int, 1)
	complaintsCh := make(chan ComplaintInsight, 1)
	usersCh := make(chan UserInsight, 1)

	wg.Add(4)
	go run(&wg, salesCh, func() salesInsight { return summarizeSales(data.sales) })
	go run(&wg, suggestionsCh, func() map[string]int { return summarizeSuggestions(data.suggestions) })
	go run(&wg, complaintsCh, func() ComplaintInsight { return summarizeComplaints(data.complaints) })
	go run(&wg, usersCh, func() UserInsight { return summarizeUsers(data.userCaptures) })

	wg.Wait()
	close(salesCh)
	close(suggestionsCh)
	close(complaintsCh)
	close(usersCh)

	sales := <-salesCh
	suggestions := <-suggestionsCh
	complaints := <-complaintsCh
	users := <-usersCh

	return InsightReport{
		GeneratedAt:         time.Now().UTC(),
		TotalRevenue:        sales.totalRevenue,
		TotalSales:          sales.totalSales,
		AverageTicket:       sales.averageTicket,
		RevenueByChannel:    sales.revenueByChannel,
		TopProducts:         sales.topProducts,
		SuggestionThemes:    suggestions,
		ComplaintSummary:    complaints,
		UserCaptureSummary:  users,
		RecommendedActions:  recommendActions(sales, suggestions, complaints, users),
		NormalizationIssues: data.issues,
	}, nil
}

func (m *InsightModel) loadAndNormalize(ctx context.Context) (normalizedData, error) {
	type result[T any] struct {
		items []T
		err   error
	}

	salesCh := make(chan result[Sale], 1)
	suggestionsCh := make(chan result[Suggestion], 1)
	complaintsCh := make(chan result[Complaint], 1)
	usersCh := make(chan result[UserCapture], 1)

	go func() {
		items, err := m.repo.ListSales(ctx)
		salesCh <- result[Sale]{items: items, err: err}
	}()
	go func() {
		items, err := m.repo.ListSuggestions(ctx)
		suggestionsCh <- result[Suggestion]{items: items, err: err}
	}()
	go func() {
		items, err := m.repo.ListComplaints(ctx)
		complaintsCh <- result[Complaint]{items: items, err: err}
	}()
	go func() {
		items, err := m.repo.ListUserCaptures(ctx)
		usersCh <- result[UserCapture]{items: items, err: err}
	}()

	salesResult := <-salesCh
	suggestionsResult := <-suggestionsCh
	complaintsResult := <-complaintsCh
	usersResult := <-usersCh

	if err := errors.Join(salesResult.err, suggestionsResult.err, complaintsResult.err, usersResult.err); err != nil {
		return normalizedData{}, err
	}

	return normalizeData(salesResult.items, suggestionsResult.items, complaintsResult.items, usersResult.items), nil
}

func normalizeData(
	sales []Sale,
	suggestions []Suggestion,
	complaints []Complaint,
	userCaptures []UserCapture,
) normalizedData {
	data := normalizedData{
		sales:        make([]Sale, 0, len(sales)),
		suggestions:  make([]Suggestion, 0, len(suggestions)),
		complaints:   make([]Complaint, 0, len(complaints)),
		userCaptures: make([]UserCapture, 0, len(userCaptures)),
	}

	for _, sale := range sales {
		normalized, issue := normalizeSale(sale)
		data.sales = append(data.sales, normalized)
		if issue != "" {
			data.issues = append(data.issues, issue)
		}
	}
	for _, suggestion := range suggestions {
		data.suggestions = append(data.suggestions, normalizeSuggestion(suggestion))
	}
	for _, complaint := range complaints {
		normalized, issue := normalizeComplaint(complaint)
		data.complaints = append(data.complaints, normalized)
		if issue != "" {
			data.issues = append(data.issues, issue)
		}
	}
	for _, capture := range userCaptures {
		data.userCaptures = append(data.userCaptures, normalizeUserCapture(capture))
	}

	return data
}

func run[T any](wg *sync.WaitGroup, out chan<- T, fn func() T) {
	defer wg.Done()
	out <- fn()
}

type salesInsight struct {
	totalRevenue     float64
	totalSales       int
	averageTicket    float64
	revenueByChannel map[string]float64
	topProducts      []ProductInsight
}

func summarizeSales(sales []Sale) salesInsight {
	byProduct := make(map[string]ProductInsight)
	byChannel := make(map[string]float64)
	var total float64

	for _, sale := range sales {
		total += sale.Amount
		byChannel[sale.Channel] += sale.Amount

		product := byProduct[sale.Product]
		product.Product = sale.Product
		product.Revenue += sale.Amount
		product.Sales++
		byProduct[sale.Product] = product
	}

	topProducts := make([]ProductInsight, 0, len(byProduct))
	for _, product := range byProduct {
		topProducts = append(topProducts, product)
	}
	sort.Slice(topProducts, func(i, j int) bool {
		return topProducts[i].Revenue > topProducts[j].Revenue
	})
	if len(topProducts) > 3 {
		topProducts = topProducts[:3]
	}

	average := 0.0
	if len(sales) > 0 {
		average = total / float64(len(sales))
	}

	return salesInsight{
		totalRevenue:     roundMoney(total),
		totalSales:       len(sales),
		averageTicket:    roundMoney(average),
		revenueByChannel: roundMoneyMap(byChannel),
		topProducts:      roundProductRevenue(topProducts),
	}
}

func summarizeSuggestions(suggestions []Suggestion) map[string]int {
	themes := make(map[string]int)
	for _, suggestion := range suggestions {
		themes[suggestion.Category]++
	}
	return themes
}

func summarizeComplaints(complaints []Complaint) ComplaintInsight {
	var totalSeverity int
	var highSeverity int
	for _, complaint := range complaints {
		totalSeverity += complaint.Severity
		if complaint.Severity >= 4 {
			highSeverity++
		}
	}

	average := 0.0
	if len(complaints) > 0 {
		average = float64(totalSeverity) / float64(len(complaints))
	}

	return ComplaintInsight{
		Total:           len(complaints),
		AverageSeverity: roundMoney(average),
		HighSeverity:    highSeverity,
	}
}

func summarizeUsers(captures []UserCapture) UserInsight {
	bySource := make(map[string]int)
	bySegment := make(map[string]int)
	totalWithConsent := 0

	for _, capture := range captures {
		if !capture.Consent {
			continue
		}
		totalWithConsent++
		bySource[capture.Source]++
		bySegment[capture.Segment]++
	}

	return UserInsight{
		TotalWithConsent: totalWithConsent,
		BySource:         bySource,
		BySegment:        bySegment,
	}
}

func recommendActions(
	sales salesInsight,
	suggestions map[string]int,
	complaints ComplaintInsight,
	users UserInsight,
) []string {
	actions := make([]string, 0, 4)

	if len(sales.topProducts) > 0 {
		actions = append(actions, "priorizar campanhas para "+sales.topProducts[0].Product)
	}
	if suggestions["produto"] >= 2 {
		actions = append(actions, "avaliar melhorias de produto mais solicitadas pelos usuarios")
	}
	if complaints.HighSeverity > 0 {
		actions = append(actions, "tratar reclamacoes de alta severidade antes de ampliar aquisicao")
	}
	if users.TotalWithConsent > 0 {
		actions = append(actions, "segmentar comunicacoes usando origem e segmento dos usuarios com consentimento")
	}
	return actions
}

func roundMoneyMap(values map[string]float64) map[string]float64 {
	rounded := make(map[string]float64, len(values))
	for key, value := range values {
		rounded[key] = roundMoney(value)
	}
	return rounded
}

func roundProductRevenue(products []ProductInsight) []ProductInsight {
	for i := range products {
		products[i].Revenue = roundMoney(products[i].Revenue)
	}
	return products
}

func roundMoney(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}
