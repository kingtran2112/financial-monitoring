package gold

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type GoldJob struct {
	scheduler   gocron.Scheduler
	goldService GoldService
}

type GoldService interface {
	FetchGoldPrice() (int, error)
}

func NewJob(scheduler gocron.Scheduler, goldService GoldService) *GoldJob {
	return &GoldJob{scheduler: scheduler, goldService: goldService}
}

func (gj *GoldJob) AddGoldPrice() error {
	j, err := gj.scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(func() {
			fmt.Println("Start task!")
			currentPrice, err := gj.goldService.FetchGoldPrice()
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			fmt.Printf("Current gold price: %d \n", currentPrice)
		}),
	)
	if err != nil {
		return err
	}
	if err != nil {
		fmt.Printf("Error: %v \n", err)
	}
	fmt.Printf("Job id: %d!", j.ID())
	return nil
}
