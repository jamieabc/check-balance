package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

const (
	fileName    = "history.json"
	host        = "https://api.blockcypher.com"
	apiEndpoint = "v1/ltc/main/addrs"
	query       = "limit=12000"
	address     = "Laz6tpw5xNhXWdiFoTfg4RnVmvVqRSFmkK"
)

type tx struct {
	Balance   int    `json:"ref_balance"`
	PrevValue int    `json:"value"`
	Time      string `json:"confirmed"`
	Input     int    `json:"tx_input_n"`
}

type transactions struct {
	TxRefs  []tx `json:"txrefs"`
	Balance int  `json:"balance"`
}

type balanceList map[string]float64

func main() {
	data, err := retrieveRemoteAddressHistory(address)
	if nil != err {
		fmt.Printf("retrieve remote address history with error: %s\n", err)
		return
	}

	err = writeDataToFile(fileName, data)
	if nil != err {
		fmt.Printf("write data to file with error: %s\n", err)
		return
	}

	var txs transactions
	err = json.Unmarshal(data, &txs)
	if nil != err {
		fmt.Printf("unmarshal json file with error: %s\n", err)
		return
	}

	balance := parseTransactions(txs)
	showBalance(balance)
	fmt.Printf("current balance: %f\n", toCoin(txs.Balance))
}

func retrieveRemoteAddressHistory(address string) ([]byte, error) {
	url := remoteURL(address)
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if nil != err {
		return nil, fmt.Errorf("new http request with error: %s", err)
	}

	req.Header.Set("User-Agent", "http")

	res, err := client.Do(req)
	if nil != err {
		return nil, fmt.Errorf("do http request with error: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return nil, fmt.Errorf("read body with error: %s", err)
	}

	return body, nil
}

func writeDataToFile(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, 0644)
}

func remoteURL(addr string) string {
	return fmt.Sprintf("%s/%s/%s?%s", host, apiEndpoint, addr, query)
}

func parseTransactions(txs transactions) balanceList {
	balance := make(balanceList)
	for i := 0; i < len(txs.TxRefs); i++ {
		tx := txs.TxRefs[i]
		if isReceiverTrx(tx) || isReceiveFund(tx) {
			continue
		}
		t, _ := time.Parse("2006-01-02T15:04:05Z", tx.Time)
		day := t.Format("2006-01-02")
		balance[day] += toCoin(tx.PrevValue - tx.Balance)
	}
	return balance
}

func isReceiverTrx(t tx) bool {
	return 0 != t.Input
}

func isReceiveFund(t tx) bool {
	return t.PrevValue < t.Balance
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
