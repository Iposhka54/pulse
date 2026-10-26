package collector

import (
	"context"
	"log"
	"time"
)

type Collector struct {
	interval       time.Duration
	latestSnapshot *Snapshot
}

func (c *Collector) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			actualSnapshot, err := SnapshotNow()
			if err != nil {
				log.Printf("Failed to analyze server metrics: %v", err)
				continue
			}
			c.latestSnapshot = actualSnapshot
		case <-ctx.Done():
			ticker.Stop()
			log.Println("stopped collect metrics")
			return
		}
	}
}

func (c *Collector) Latest() *Snapshot {
	return c.latestSnapshot
}

func New(interval time.Duration) *Collector {
	return &Collector{
		interval: interval,
	}
}
