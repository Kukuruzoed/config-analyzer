package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kukuruzoed/config-analyzer/internal/analyzer"
	"github.com/Kukuruzoed/config-analyzer/internal/parser"
	"github.com/Kukuruzoed/config-analyzer/internal/rules"
)

type AnalyzeRequest struct {
	Content string `json:"content"`
}

type AnalyzeResponse struct {
	Issues []IssueDTO `json:"issues"`
	Total  int        `json:"total"`
}

type IssueDTO struct {
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
}

type Server struct {
	analyzer *analyzer.Analyzer
	port     int
}

func New(a *analyzer.Analyzer, port int) *Server {
	return &Server{analyzer: a, port: port}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /analyze", s.handleAnalyze)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("HTTP сервер запущен на %s", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "невалидный JSON: "+err.Error())
		return
	}

	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "поле content обязательно")
		return
	}

	config, err := parser.ParseString(req.Content)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "ошибка парсинга: "+err.Error())
		return
	}

	issues := s.analyzer.Analyze(config)

	resp := AnalyzeResponse{
		Issues: toDTO(issues),
		Total:  len(issues),
	}

	w.Header().Set("Content-Type", "application/json")
	if len(issues) > 0 {
		w.WriteHeader(http.StatusMultiStatus)
	}
	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func toDTO(issues []rules.Issue) []IssueDTO {
	dto := make([]IssueDTO, len(issues))
	for i, issue := range issues {
		dto[i] = IssueDTO{
			Severity:       string(issue.Severity),
			Description:    issue.Description,
			Recommendation: issue.Recommendation,
		}
	}
	return dto
}
