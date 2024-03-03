package gold

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
)

type GoldJob struct {
	scheduler   gocron.Scheduler
	goldService GoldService
}

type GoldService interface {
	FetchGoldPrice() (int, error)
	AddGoldPrice(price int) (int, error)
}

func NewJob(scheduler gocron.Scheduler, goldService GoldService) *GoldJob {
	return &GoldJob{scheduler: scheduler, goldService: goldService}
}

func (gj *GoldJob) AddGoldPrice() error {
	j, err := gj.scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(gocron.NewAtTime(7, 0, 0)),
		),
		gocron.NewTask(func() {
			fmt.Println("Start task!")
			currentPrice, err := gj.goldService.FetchGoldPrice()
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			gj.goldService.AddGoldPrice(currentPrice)
		}),
	)
	if err != nil {
		return err
	}
	return j.RunNow()
}
