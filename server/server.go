package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type USDExchange struct {
	UsdBrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoHandler)
	log.Println("[INFO] Servidor HTTP rodando em: http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] %s %s - requisição recebida", r.Method, r.URL.Path)
	cotacao, err := findUSDExchange()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[ERROR] Erro ao consumir API: %v", err)
		return
	}
	res := CotacaoResponse{cotacao.UsdBrl.Bid}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	log.Printf("[INFO] %s %s - resposta enviada com status %d", r.Method, r.URL.Path, http.StatusOK)
}

func findUSDExchange() (*USDExchange, error) {
	res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var exchange USDExchange
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		return nil, err
	}

	return &exchange, nil
}
