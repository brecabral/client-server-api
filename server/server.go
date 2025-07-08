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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoHandler)
	log.Println("Iniciando Servidor...")
	http.ListenAndServe(":8080", mux)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Requisição Recebida")
	cotacao, err := findUSDExchange()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Erro ao consumir API: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)
	log.Println("Requisição Respondida")
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
