package main

import (
	"context"
	"financial-monitoring/db"
	"financial-monitoring/gold"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-co-op/gocron/v2"
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
	s.Start()
	defer s.Shutdown()

	if err != nil {
		fmt.Printf("Error in main: %v", err)
		panic(err)
	}
	goldService := gold.NewService(influxClient)
	goldJob := gold.NewJob(s, goldService)
	goldJob.AddGoldPrice()

	// Keep the application keep running
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
}
