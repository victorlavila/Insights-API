package model

import "time"

type Sale struct {
	ID        string
	Product   string
	Amount    float64
	Channel   string
	CreatedAt time.Time
}

type Suggestion struct {
	ID        string
	UserID    string
	Text      string
	Category  string
	CreatedAt time.Time
}

type Complaint struct {
	ID        string
	UserID    string
	Text      string
	Severity  int
	CreatedAt time.Time
}

type UserCapture struct {
	UserID    string
	Source    string
	Segment   string
	Consent   bool
	CreatedAt time.Time
}

type InsightReport struct {
	GeneratedAt         time.Time          `json:"generated_at"`
	TotalRevenue        float64            `json:"total_revenue"`
	TotalSales          int                `json:"total_sales"`
	AverageTicket       float64            `json:"average_ticket"`
	RevenueByChannel    map[string]float64 `json:"revenue_by_channel"`
	TopProducts         []ProductInsight   `json:"top_products"`
	SuggestionThemes    map[string]int     `json:"suggestion_themes"`
	ComplaintSummary    ComplaintInsight   `json:"complaint_summary"`
	UserCaptureSummary  UserInsight        `json:"user_capture_summary"`
	RecommendedActions  []string           `json:"recommended_actions"`
	NormalizationIssues []string           `json:"normalization_issues,omitempty"`
}

type ProductInsight struct {
	Product string  `json:"product"`
	Revenue float64 `json:"revenue"`
	Sales   int     `json:"sales"`
}

type ComplaintInsight struct {
	Total           int     `json:"total"`
	AverageSeverity float64 `json:"average_severity"`
	HighSeverity    int     `json:"high_severity"`
}

type UserInsight struct {
	TotalWithConsent int            `json:"total_with_consent"`
	BySource         map[string]int `json:"by_source"`
	BySegment        map[string]int `json:"by_segment"`
}
