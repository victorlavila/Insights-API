package service

import (
	"context"
	"testing"

	"api_insights/internal/domain"
	"api_insights/internal/repository"
)

func TestNormalizeData(t *testing.T) {
	data := normalizeData(
		[]domain.Sale{{ID: "s1", Product: " Plano Pro ", Amount: -10, Channel: " WEB "}},
		[]domain.Suggestion{{ID: "sg1", Category: " Produto ", Text: "  muitos   filtros  "}},
		[]domain.Complaint{{ID: "c1", Severity: 8, Text: " erro   crítico "}},
		[]domain.UserCapture{{UserID: "u1", Source: " Ads ", Segment: " smb ", Consent: true}},
	)

	if got := data.sales[0].Product; got != "plano pro" {
		t.Fatalf("produto normalizado = %q", got)
	}
	if got := data.sales[0].Amount; got != 0 {
		t.Fatalf("valor negativo deveria virar 0, recebeu %.2f", got)
	}
	if got := data.suggestions[0].Text; got != "muitos filtros" {
		t.Fatalf("texto normalizado = %q", got)
	}
	if got := data.complaints[0].Severity; got != 5 {
		t.Fatalf("severidade deveria ser limitada em 5, recebeu %d", got)
	}
	if got := data.userCaptures[0].Segment; got != "SMB" {
		t.Fatalf("segmento normalizado = %q", got)
	}
	if len(data.issues) != 2 {
		t.Fatalf("issues = %d, queria 2", len(data.issues))
	}
}

func TestGenerateReport(t *testing.T) {
	service := NewInsightService(repository.NewMockStore())

	report, err := service.GenerateReport(context.Background())
	if err != nil {
		t.Fatalf("GenerateReport retornou erro: %v", err)
	}

	if report.TotalSales != 5 {
		t.Fatalf("TotalSales = %d, queria 5", report.TotalSales)
	}
	if report.TotalRevenue != 1947.7 {
		t.Fatalf("TotalRevenue = %.2f, queria 1947.70", report.TotalRevenue)
	}
	if report.RevenueByChannel["web"] != 1747.9 {
		t.Fatalf("receita web = %.2f, queria 1747.90", report.RevenueByChannel["web"])
	}
	if report.TopProducts[0].Product != "consultoria" {
		t.Fatalf("top produto = %q, queria consultoria", report.TopProducts[0].Product)
	}
	if report.ComplaintSummary.HighSeverity != 2 {
		t.Fatalf("reclamacoes graves = %d, queria 2", report.ComplaintSummary.HighSeverity)
	}
	if report.UserCaptureSummary.TotalWithConsent != 3 {
		t.Fatalf("usuarios com consentimento = %d, queria 3", report.UserCaptureSummary.TotalWithConsent)
	}
	if len(report.RecommendedActions) == 0 {
		t.Fatal("esperava recomendacoes no relatorio")
	}
}
