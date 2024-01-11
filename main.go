package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxApi "github.com/influxdata/influxdb-client-go/v2/api"
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
	url := "http://localhost:8086"
	token := "K3eRmCol-bR-gP6PspDxPhj9ZwmQZQetONjaA8nLWfPDyM8dHQfjbba1PjBQ-2oAFJdlGf1ai7AKML3d845zow=="
	client, closer := initDB(url, token)
	defer closer()

	org := "f25a23958d7394f3"
	bucket := "Financial"
	writeAPI := client.WriteAPIBlocking(org, bucket)

	spending := getDataFromFile("financial_report.csv")
	writeSpending(writeAPI, spending)
}

func initDB(url string, token string) (influxdb2.Client, func()) {
	client := influxdb2.NewClient(url, token)

	return client, func() {
		client.Close()
	}
}

func getDataFromFile(fileName string) []*Spending {
	f, err := os.Open(fileName)
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

	return spending
}

func writeSpending(writeAPI influxApi.WriteAPIBlocking, spending []*Spending) {
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
}
