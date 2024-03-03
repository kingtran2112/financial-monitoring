package importing

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type influxClient interface {
	WritePoint(p *write.Point) error
}

type importingService struct {
	influx influxClient
}

func (is *importingService) Import(path string) error {
	spending, err := is.getDataFromFile(path)
	if err != nil {
		return err
	}
	is.writeSpending(spending)
	return nil
}

func (is *importingService) getDataFromFile(fileName string) ([]*Spending, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var spending []*Spending

	if err := gocsv.UnmarshalFile(f, &spending); err != nil {
		return nil, err
	}

	for _, s := range spending {
		if s.Amount < 0 {
			s.Type = EXPENSE
			s.Amount = -s.Amount
		} else {
			s.Type = INCOME
		}
		if !s.Type.IsValid() {
			return nil, fmt.Errorf("invalid spending type: %v", s.Type)
		}
	}

	return spending, nil
}

// Todo: consider to return (n, error). With n is the total of written spendings.
func (is *importingService) writeSpending(spending []*Spending) {
	for _, s := range spending {
		p, err := is.spendingToPoint(s)
		if err != nil {
			log.Printf("convert spending to point: %v\n", err)
			continue
		}
		if err := is.influx.WritePoint(p); err != nil {
			log.Printf("write point: %v\n", err)
			continue
		}
	}
}

func (is *importingService) spendingToPoint(s *Spending) (*write.Point, error) {
	fmt.Printf("Writing %s %s %s %s %d %s\n", s.Wallet, s.Date, s.Type, s.Group, s.Amount, s.Currency)
	date, err := time.Parse("02/01/2006", s.Date)
	if err != nil {
		return nil, err
	}
	return influxdb2.NewPoint(
		s.Wallet,
		map[string]string{"group": s.Group, "type": s.Type.String()},
		map[string]interface{}{"amount": s.Amount, "currency": s.Currency, "note": s.Note},
		date,
	), nil
}

func NewService(influxClient influxClient) *importingService {
	return &importingService{
		influx: influxClient,
	}
}
