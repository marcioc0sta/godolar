package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// reads response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	q := Quote{}
	err = json.Unmarshal(body, &q)

	if err != nil {
		panic(err)
	}

	// creates a txt file with response
	f, err := os.Create("response.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("Dolar: %s", q.USDBRL.Bid))
	if err != nil {
		fmt.Println(err)
		return
	}
}
