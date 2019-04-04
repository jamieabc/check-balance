package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type tx struct {
	Balance   int    `json:"ref_balance"`
	PrevValue int    `json:"value"`
	Time      string `json:"confirmed"`
	Input     int    `json:"tx_input_n"`
}

type transactions struct {
	TxRefs  []tx `json:txrefs`
	Balance int  `json:"balance"`
}

func main() {
	jsonFile, err := os.Open("history.json")
	if nil != err {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var txs transactions

	json.Unmarshal(byteValue, &txs)

	balance := make(map[string]float64)

	for i := 0; i < len(txs.TxRefs); i++ {
		tx := txs.TxRefs[i]
		if 0 != tx.Input {
			continue
		}
		t, _ := time.Parse("2006-01-02T15:04:05Z", tx.Time)
		balance[t.Format("2006-01-02")] +=
			toCoin(tx.PrevValue - tx.Balance)
	}
	showBalance(balance)
	fmt.Printf("current balance: %f\n", toCoin(txs.Balance))
}

func toCoin(value int) float64 {
	return float64(value) / 100000000
}

func showBalance(m map[string]float64) {
	sum := float64(0)
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := m[key]
		fmt.Printf("date: %s, daily spend: %f\n", key, value)
		sum += value
	}
	fmt.Printf("total spend: %f\n", sum)
}
