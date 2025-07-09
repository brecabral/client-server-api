package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	res, err := http.Get("http://localhost:8080/cotacao")
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
