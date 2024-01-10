package main

import (
	"context"
	"financial-monitoring/gold"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	// goldService := gold.NewService()
	// price, err := goldService.FetchGoldPrice()
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)
	// }
	// fmt.Printf("Current gold price: %d", price)

	// panic("Implementing!!!!")
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
	goldService := gold.NewService()
	goldJob := gold.NewJob(s, goldService)
	goldJob.AddGoldPrice()

	// Keep the application keep running
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
}
