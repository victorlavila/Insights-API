package model

import (
	"context"
	"sync"
	"time"
)

type InsightRepository interface {
	ListSales(ctx context.Context) ([]Sale, error)
	ListSuggestions(ctx context.Context) ([]Suggestion, error)
	ListComplaints(ctx context.Context) ([]Complaint, error)
	ListUserCaptures(ctx context.Context) ([]UserCapture, error)
}

type MockStore struct {
	mu           sync.RWMutex
	sales        []Sale
	suggestions  []Suggestion
	complaints   []Complaint
	userCaptures []UserCapture
}

func NewMockStore() *MockStore {
	now := time.Now().UTC()

	return &MockStore{
		sales: []Sale{
			{ID: "s1", Product: " Plano Pro ", Amount: 149.90, Channel: "Web", CreatedAt: now.Add(-72 * time.Hour)},
			{ID: "s2", Product: "plano pro", Amount: 149.90, Channel: " app ", CreatedAt: now.Add(-48 * time.Hour)},
			{ID: "s3", Product: "Consultoria", Amount: 799.00, Channel: "web", CreatedAt: now.Add(-24 * time.Hour)},
			{ID: "s4", Product: "Plano Basic", Amount: 49.90, Channel: "Parceiros", CreatedAt: now.Add(-12 * time.Hour)},
			{ID: "s5", Product: "consultoria", Amount: 799.00, Channel: "web", CreatedAt: now.Add(-6 * time.Hour)},
		},
		suggestions: []Suggestion{
			{ID: "sg1", UserID: "u1", Text: "Criar dashboard com filtros por periodo", Category: " Produto ", CreatedAt: now.Add(-70 * time.Hour)},
			{ID: "sg2", UserID: "u2", Text: "Integrar relatorios por email", Category: "produto", CreatedAt: now.Add(-30 * time.Hour)},
			{ID: "sg3", UserID: "u3", Text: "Melhorar onboarding inicial", Category: "Experiencia", CreatedAt: now.Add(-20 * time.Hour)},
		},
		complaints: []Complaint{
			{ID: "c1", UserID: "u4", Text: "Checkout lento no celular", Severity: 4, CreatedAt: now.Add(-50 * time.Hour)},
			{ID: "c2", UserID: "u5", Text: "Demora no suporte", Severity: 3, CreatedAt: now.Add(-25 * time.Hour)},
			{ID: "c3", UserID: "u6", Text: "Erro ao exportar relatório", Severity: 5, CreatedAt: now.Add(-10 * time.Hour)},
		},
		userCaptures: []UserCapture{
			{UserID: "u1", Source: "Landing Page", Segment: "SMB", Consent: true, CreatedAt: now.Add(-90 * time.Hour)},
			{UserID: "u2", Source: "Ads", Segment: "Enterprise", Consent: true, CreatedAt: now.Add(-60 * time.Hour)},
			{UserID: "u3", Source: "ads", Segment: "smb", Consent: true, CreatedAt: now.Add(-40 * time.Hour)},
			{UserID: "u4", Source: "Indicação", Segment: "SMB", Consent: false, CreatedAt: now.Add(-15 * time.Hour)},
		},
	}
}

func (s *MockStore) ListSales(ctx context.Context) ([]Sale, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return clone(s.sales), ctx.Err()
}

func (s *MockStore) ListSuggestions(ctx context.Context) ([]Suggestion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return clone(s.suggestions), ctx.Err()
}

func (s *MockStore) ListComplaints(ctx context.Context) ([]Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return clone(s.complaints), ctx.Err()
}

func (s *MockStore) ListUserCaptures(ctx context.Context) ([]UserCapture, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return clone(s.userCaptures), ctx.Err()
}

func clone[T any](items []T) []T {
	out := make([]T, len(items))
	copy(out, items)
	return out
}
