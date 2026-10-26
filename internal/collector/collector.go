package collector

import (
	"context"
	"log"
	"sync"
	"time"
)

type Collector struct {
	mu       sync.RWMutex
	interval time.Duration
	latest   *Snapshot
}

func New(interval time.Duration) *Collector {
	return &Collector{
		interval: interval,
	}
}

func (c *Collector) Run(ctx context.Context) {
	c.collectOnce()
	ticker := time.NewTicker(c.interval)
	for {
		select {
		case <-ticker.C:
			c.collectOnce()
		case <-ctx.Done():
			ticker.Stop()
			log.Println("stopped collect metrics")
			return
		}
	}
}

func (c *Collector) Latest() *Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.latest
}

func (c *Collector) collectOnce() {
	snap, err := SnapshotNow()
	if err != nil {
		log.Printf("failed to analyze server metrics: %v", err)
		return
	}
	c.mu.Lock()
	c.latest = snap
	c.mu.Unlock()
}
