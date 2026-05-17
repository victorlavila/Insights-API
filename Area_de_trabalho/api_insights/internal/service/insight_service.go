package service

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"api_insights/internal/domain"
	"api_insights/internal/repository"
)

type InsightService struct {
	repo repository.InsightRepository
}

func NewInsightService(repo repository.InsightRepository) *InsightService {
	return &InsightService{repo: repo}
}

func (s *InsightService) GenerateReport(ctx context.Context) (domain.InsightReport, error) {
	data, err := s.loadAndNormalize(ctx)
	if err != nil {
		return domain.InsightReport{}, err
	}

	var wg sync.WaitGroup
	salesCh := make(chan salesInsight, 1)
	suggestionsCh := make(chan map[string]int, 1)
	complaintsCh := make(chan domain.ComplaintInsight, 1)
	usersCh := make(chan domain.UserInsight, 1)

	wg.Add(4)
	go run(&wg, salesCh, func() salesInsight { return summarizeSales(data.sales) })
	go run(&wg, suggestionsCh, func() map[string]int { return summarizeSuggestions(data.suggestions) })
	go run(&wg, complaintsCh, func() domain.ComplaintInsight { return summarizeComplaints(data.complaints) })
	go run(&wg, usersCh, func() domain.UserInsight { return summarizeUsers(data.userCaptures) })

	wg.Wait()
	close(salesCh)
	close(suggestionsCh)
	close(complaintsCh)
	close(usersCh)

	sales := <-salesCh
	suggestions := <-suggestionsCh
	complaints := <-complaintsCh
	users := <-usersCh

	return domain.InsightReport{
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

func (s *InsightService) loadAndNormalize(ctx context.Context) (normalizedData, error) {
	type result[T any] struct {
		items []T
		err   error
	}

	salesCh := make(chan result[domain.Sale], 1)
	suggestionsCh := make(chan result[domain.Suggestion], 1)
	complaintsCh := make(chan result[domain.Complaint], 1)
	usersCh := make(chan result[domain.UserCapture], 1)

	go func() {
		items, err := s.repo.ListSales(ctx)
		salesCh <- result[domain.Sale]{items: items, err: err}
	}()
	go func() {
		items, err := s.repo.ListSuggestions(ctx)
		suggestionsCh <- result[domain.Suggestion]{items: items, err: err}
	}()
	go func() {
		items, err := s.repo.ListComplaints(ctx)
		complaintsCh <- result[domain.Complaint]{items: items, err: err}
	}()
	go func() {
		items, err := s.repo.ListUserCaptures(ctx)
		usersCh <- result[domain.UserCapture]{items: items, err: err}
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
	sales []domain.Sale,
	suggestions []domain.Suggestion,
	complaints []domain.Complaint,
	userCaptures []domain.UserCapture,
) normalizedData {
	data := normalizedData{
		sales:        make([]domain.Sale, 0, len(sales)),
		suggestions:  make([]domain.Suggestion, 0, len(suggestions)),
		complaints:   make([]domain.Complaint, 0, len(complaints)),
		userCaptures: make([]domain.UserCapture, 0, len(userCaptures)),
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
	topProducts      []domain.ProductInsight
}

func summarizeSales(sales []domain.Sale) salesInsight {
	byProduct := make(map[string]domain.ProductInsight)
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

	topProducts := make([]domain.ProductInsight, 0, len(byProduct))
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

func summarizeSuggestions(suggestions []domain.Suggestion) map[string]int {
	themes := make(map[string]int)
	for _, suggestion := range suggestions {
		themes[suggestion.Category]++
	}
	return themes
}

func summarizeComplaints(complaints []domain.Complaint) domain.ComplaintInsight {
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

	return domain.ComplaintInsight{
		Total:           len(complaints),
		AverageSeverity: roundMoney(average),
		HighSeverity:    highSeverity,
	}
}

func summarizeUsers(captures []domain.UserCapture) domain.UserInsight {
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

	return domain.UserInsight{
		TotalWithConsent: totalWithConsent,
		BySource:         bySource,
		BySegment:        bySegment,
	}
}

func recommendActions(
	sales salesInsight,
	suggestions map[string]int,
	complaints domain.ComplaintInsight,
	users domain.UserInsight,
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

func roundProductRevenue(products []domain.ProductInsight) []domain.ProductInsight {
	for i := range products {
		products[i].Revenue = roundMoney(products[i].Revenue)
	}
	return products
}

func roundMoney(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}
