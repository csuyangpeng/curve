package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"
)

type stocksResponse struct {
	UpdatedAt string  `json:"updated_at"`
	Count     int     `json:"count"`
	Stocks    []Stock `json:"stocks"`
}

type codesResponse struct {
	Codes []string `json:"codes"`
}

type errorResponse struct {
	Detail string `json:"detail"`
}

func frontendDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "../frontend"
	}
	return filepath.Join(filepath.Dir(file), "..", "frontend")
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func handleStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := getStocks()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(errorResponse{Detail: "Failed to fetch stock data"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stocksResponse{
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Count:     len(stocks),
		Stocks:    stocks,
	})
}

func handleCodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(codesResponse{Codes: stockCodes})
}

func main() {
	dir := frontendDir()
	static := http.FileServer(http.Dir(dir))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/stocks", withCORS(handleStocks))
	mux.HandleFunc("/api/codes", withCORS(handleCodes))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(dir, "index.html"))
	})
	mux.Handle("/static/", http.StripPrefix("/static/", static))

	addr := ":8000"
	log.Printf("server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
