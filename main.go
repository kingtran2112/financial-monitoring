package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-co-op/gocron/v2"

	"financial-monitoring/db"
	"financial-monitoring/gold"
)

func main() {
	// TODO: These information should be imported from env
	// TODO: Create a flag for importing function
	url := "http://localhost:8086"
	token := "eTjVDmFXk38b-6312uMIctjZGUnCuyil_hRQaioiP7HDOyXixL4pu_TEWVd5a_hhlP4rzE72WpsLAAabxmr2hQ=="
	org := "my-org"
	bucket := "financial"

	influxClient := db.NewInfluxClient(url, token, org, bucket)
	defer influxClient.Close()
	// importingService := importing.NewService(influxClient)

	// importingService.Import("financial_report.csv")

	s, err := gocron.NewScheduler(
		gocron.WithLogger(
			gocron.NewLogger(gocron.LogLevelDebug),
		),
	)
	if err != nil {
		log.Fatalf("new scheduler: %s\n", err)
	}

	s.Start()
	defer func() {
		if err := s.Shutdown(); err != nil {
			log.Printf("shutdown scheduler: %s\n", err)
		}
	}()

	goldService := gold.NewService(influxClient)
	goldJob := gold.NewJob(s, goldService)
	goldJob.AddGoldPrice()

	// Keep the application running
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
}
