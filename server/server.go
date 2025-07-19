package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

type Server struct {
	db *sql.DB
	client *http.Client
}

const createTable = `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		var_bid TEXT,
		pct_change TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	);`

const insertCotacao = `
	INSERT INTO cotacoes(
		code,
		codein,
		name,
		high,
		low,
		var_bid,
		pct_change,
		bid,
		ask,
		timestamp,
		create_date
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

const urlCotacao = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {

	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatalf("[ERROR] Erro ao se conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("[ERROR] Erro ao criar tabela no banco de dados: %v", err)
	}

	client := &http.Client{}
	
	server := Server{db, client}

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", server.cotacaoHandler)
	log.Println("[INFO] Servidor HTTP rodando em: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (s *Server) cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] %s %s - requisição recebida", r.Method, r.URL.Path)

	cotacao, err := s.findCotacao()
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

	err = s.saveCotacao(cotacao)
	if err != nil {
		log.Printf("[ERROR] Erro ao salvar cotação: %v", err)
		return
	}
}

func (s *Server) findCotacao() (*USDExchange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlCotacao, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
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

func (s *Server) saveCotacao(cotacao *USDExchange) error {
	stmt, err := s.db.Prepare(insertCotacao)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	
	_, err = stmt.ExecContext(
		ctx,
		cotacao.UsdBrl.Code,
		cotacao.UsdBrl.Codein,
		cotacao.UsdBrl.Name,
		cotacao.UsdBrl.High,
		cotacao.UsdBrl.Low,
		cotacao.UsdBrl.VarBid,
		cotacao.UsdBrl.PctChange,
		cotacao.UsdBrl.Bid,
		cotacao.UsdBrl.Ask,
		cotacao.UsdBrl.Timestamp,
		cotacao.UsdBrl.CreateDate)
	if err != nil {
		return err
	}
	
	return nil
}
