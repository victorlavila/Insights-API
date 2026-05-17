package controller

import (
	"net/http"

	"api_insights/internal/model"
	"api_insights/internal/view"
)

type InsightController struct {
	model *model.InsightModel
}

func NewInsightController(model *model.InsightModel) *InsightController {
	return &InsightController{model: model}
}

func (c *InsightController) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", c.health)
	mux.HandleFunc("GET /insights/report", c.report)
	return mux
}

func (c *InsightController) health(w http.ResponseWriter, _ *http.Request) {
	view.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (c *InsightController) report(w http.ResponseWriter, r *http.Request) {
	report, err := c.model.GenerateReport(r.Context())
	if err != nil {
		view.Error(w, http.StatusInternalServerError, err)
		return
	}
	view.JSON(w, http.StatusOK, report)
}
