package tools

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Task struct {
	timer  time.Duration
	ticker time.Ticker
	quit   chan int
}

func NewTask(initialDelay time.Duration, tick time.Duration) Task {
	return Task{
		timer:  tick,
		ticker: *time.NewTicker(initialDelay),
		quit:   make(chan int),
	}
}

func (t Task) Start(task func(), runAsync bool) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			t.close()
			fmt.Println("Task stopped!")
		}()
		firstExec := true
		for {
			select {
			case <-t.ticker.C:

				if firstExec {
					t.ticker.Stop()
					t.ticker = *time.NewTicker(t.timer)
					firstExec = false
				}

				if runAsync {
					go task()
				} else {
					task()
				}

				break
			case <-t.quit:
				return
			case <-sigs:
				_ = t.Close()
				break
			}
		}

	}()

}

func (t *Task) Close() error {
	go func() {
		t.quit <- 1
	}()
	return nil
}

func (t *Task) close() {
	t.ticker.Stop()
}
