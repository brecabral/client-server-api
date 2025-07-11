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

	var cotacao Cotacao
	err := fetchCotacao(&cotacao)
	if err != nil {
		log.Fatalf("[ERROR] Erro ao acessar servidor: %v", err)
	}

	err = writeCotacao(&cotacao)
	if err != nil {
		log.Fatalf("[ERROR] Erro ao escrever no arquivo: %v", err)
	}

	log.Printf("[Info] Cotação (Dólar: %s) escrita em cotacao.txt", cotacao.Bid)
}

func fetchCotacao(cotacao *Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&cotacao)
	if err != nil {
		return err
	}

	return nil
}

func writeCotacao(cotacao *Cotacao) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString("Dólar: " + cotacao.Bid)
	if err != nil {
		return err
	}
	return nil
}
