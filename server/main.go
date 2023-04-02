package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseFile = "./quotes.db"
)

type Quote struct {
	USDBRL struct {
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
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS quotes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code REAL,
			code_in REAL,
			name REAL,
			high REAL,
			low REAL,
			var_bid REAL,
			pct_change REAL,
			bid REAL,
			ask REAL,
			timestamp REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		priceHandler(w, r, db)
	})

	http.ListenAndServe(":8080", nil)
}

func priceHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// creates context
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// creates request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// executes request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	// reads response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	q := Quote{}
	err = json.Unmarshal(body, &q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = insertQuote(db, q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// returns response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(q)
}

func insertQuote(db *sql.DB, q Quote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO quotes (code, code_in, name, high, low, var_bid, pct_change, bid, ask, timestamp, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		q.USDBRL.Code,
		q.USDBRL.Codein,
		q.USDBRL.Name,
		q.USDBRL.High,
		q.USDBRL.Low,
		q.USDBRL.VarBid,
		q.USDBRL.PctChange,
		q.USDBRL.Bid,
		q.USDBRL.Ask,
		q.USDBRL.Timestamp,
		q.USDBRL.CreateDate,
	)
	if err != nil {
		return err
	}

	return nil
}
