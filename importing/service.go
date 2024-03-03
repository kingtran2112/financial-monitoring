package importing

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
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

type influxClient interface {
	WritePoint(p *write.Point) error
}

type importingService struct {
	influx influxClient
}

func (is *importingService) Import(path string) {
	spending := is.getDataFromFile(path)
	is.writeSpending(spending)
}

func (is *importingService) getDataFromFile(fileName string) []*Spending {
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

func (is *importingService) writeSpending(spending []*Spending) {
	for _, s := range spending {
		p := is.spendingToPoint(s)
		if err := is.influx.WritePoint(p); err != nil {
			panic(err)
		}
	}
}

func (is *importingService) spendingToPoint(s *Spending) *write.Point {
	fmt.Printf("Writing %s %s %s %s %d %s\n", s.Wallet, s.Date, s.Type, s.Group, s.Amount, s.Currency)
	date, err := time.Parse("02/01/2006", s.Date)
	if err != nil {
		panic(err)
	}
	return influxdb2.NewPoint(
		s.Wallet,
		map[string]string{"group": s.Group, "type": s.Type.String()},
		map[string]interface{}{"amount": s.Amount, "currency": s.Currency, "note": s.Note},
		date,
	)
}

func NewService(influxClient influxClient) *importingService {
	return &importingService{
		influx: influxClient,
	}
}
