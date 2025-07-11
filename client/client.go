package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("[ERROR] Erro ao acessar servidor: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("[ERROR] Erro ao acessar servidor: %v", err)
	}
	defer res.Body.Close()

	var cotacao Cotacao
	err = json.NewDecoder(res.Body).Decode(&cotacao)
	if err != nil {
		log.Fatalf("[ERROR] Erro ao decodificar resposta: %v", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("[ERROR] Erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("Dólar: " + cotacao.Bid)
	if err != nil {
		log.Fatalf("[ERROR] Erro ao escrever no arquivo: %v", err)
	}

	log.Print("[Info] Cotação escrita em cotacao.txt")
}
