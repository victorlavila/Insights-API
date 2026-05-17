package model

import (
	"strings"
	"unicode"
)

type normalizedData struct {
	sales        []Sale
	suggestions  []Suggestion
	complaints   []Complaint
	userCaptures []UserCapture
	issues       []string
}

func normalizeSale(sale Sale) (Sale, string) {
	sale.Product = normalizeLabel(sale.Product)
	sale.Channel = normalizeLabel(sale.Channel)
	if sale.Amount < 0 {
		sale.Amount = 0
		return sale, "venda " + sale.ID + " possuia valor negativo e foi zerada"
	}
	return sale, ""
}

func normalizeSuggestion(suggestion Suggestion) Suggestion {
	suggestion.Text = normalizeText(suggestion.Text)
	suggestion.Category = normalizeLabel(suggestion.Category)
	return suggestion
}

func normalizeComplaint(complaint Complaint) (Complaint, string) {
	complaint.Text = normalizeText(complaint.Text)
	if complaint.Severity < 1 {
		complaint.Severity = 1
		return complaint, "reclamacao " + complaint.ID + " possuia severidade abaixo do minimo"
	}
	if complaint.Severity > 5 {
		complaint.Severity = 5
		return complaint, "reclamacao " + complaint.ID + " possuia severidade acima do maximo"
	}
	return complaint, ""
}

func normalizeUserCapture(capture UserCapture) UserCapture {
	capture.Source = normalizeLabel(capture.Source)
	capture.Segment = strings.ToUpper(strings.TrimSpace(capture.Segment))
	return capture
}

func normalizeLabel(value string) string {
	return strings.ToLower(normalizeText(value))
}

func normalizeText(value string) string {
	return strings.Join(strings.FieldsFunc(strings.TrimSpace(value), unicode.IsSpace), " ")
}
