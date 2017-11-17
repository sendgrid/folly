package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var cfg *Config

func main() {
	log.Println("Running your load")

	cfg = new(Config)
	err := envconfig.Process("RUNNER", cfg)
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
		return
	}

	for i := 0; i < cfg.Iterations; i++ {
		var wg sync.WaitGroup
		start := time.Now()
		wg.Add(cfg.Concurrency)
		for j := 0; j < cfg.Concurrency; j++ {
			go func() {
				defer wg.Done()
				_, err := http.Get(fmt.Sprintf("http://%s:%d/work", cfg.Address, cfg.Port))
				if err != nil {
					log.Printf("Error in GET: %v", err)
				}
			}()
		}
		wg.Wait()
		t := time.Since(start)
		seconds := float64(t.Nanoseconds()) / 1e9
		rps := float64(cfg.Concurrency) / seconds
		log.Printf("Iteration %d took %v with RPS of %f", i, t, rps)
	}
	log.Println("Done")
}
