package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Spending struct {
	Id       string `csv:"Id"`
	Date     string `csv:"Date"`
	Group    string `csv:"Group"`
	Amount   int32  `csv:"Amount"`
	Currency string `csv:"Currency"`
	Note     string `csv:"Note"`
	Wallet   string `csv:"Wallet"`
	Type     SpendingType
}

type SpendingType string

const (
	INCOME  SpendingType = "INCOME"
	EXPENSE SpendingType = "EXPENSE"
)

func (t SpendingType) IsValid() bool {
	switch t {
	case INCOME, EXPENSE:
		return true
	default:
		return false
	}
}

func (t SpendingType) String() string {
	return string(t)
}

func main() {
	f, err := os.Open("financial_report.csv")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var spending []*Spending

	if err := gocsv.UnmarshalFile(f, &spending); err != nil {
		panic(err)
	}

	for _, s := range spending {
		if s.Amount < 0 {
			s.Type = EXPENSE
			s.Amount = -s.Amount
		} else {
			s.Type = INCOME
		}
		if !s.Type.IsValid() {
			panic("Invalid spending type")
		}
	}

	writeSpending(spending)
}

func writeSpending(spending []*Spending) {
	url := "http://localhost:8086"
	token := "eTjVDmFXk38b-6312uMIctjZGUnCuyil_hRQaioiP7HDOyXixL4pu_TEWVd5a_hhlP4rzE72WpsLAAabxmr2hQ=="
	org := "15f0762da5e84762"
	bucket := "financial"

	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPIBlocking(org, bucket)

	for _, s := range spending {
		fmt.Printf("Writing %s %s %s %s %d %s\n", s.Wallet, s.Date, s.Type, s.Group, s.Amount, s.Currency)
		date, err := time.Parse("02/01/2006", s.Date)
		if err != nil {
			panic(err)
		}
		p := influxdb2.NewPoint(s.Wallet,
			map[string]string{"group": s.Group, "type": s.Type.String()},
			map[string]interface{}{"amount": s.Amount, "currency": s.Currency, "note": s.Note},
			date)
		err = writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			panic(err)
		}
	}

	err := writeAPI.Flush(context.Background())
	if err != nil {
		panic(err)
	}

	client.Close()
}
