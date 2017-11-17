package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	lambda "github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
)

// Handle is the entry point for a Lambda function instance.
func Handle(evt interface{}, ctx *lambda.Context) (string, error) {
	cfg = new(Config)
	err := envconfig.Process("FOLLY", cfg)
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
		return "Error loading configuration", err
	}
	cc = make(chan int, 1)
	work()
	return "Success", nil
}

var (
	cfg *Config
	cc  chan int
)

func main() {
	log.Println("This is folly!")

	// Load the configuration, or exit if there is an issue.
	cfg = new(Config)
	err := envconfig.Process("FOLLY", cfg)
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
		return
	}

	// Create a channel for counting the number of requests handled.
	cc = make(chan int, 200)

	// Go routine for profiling the request performance.
	go func() {
		count := 0
		start := time.Now()
		for {
			count += <-cc
			if count%100 == 0 {
				t := time.Since(start)
				seconds := float64(t.Nanoseconds()) / 1e9
				rps := float64(100) / seconds
				log.Printf("Received %d calls in %v with RPS of %f", count, t, rps)
				start = time.Now()

				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				log.Printf("\nAlloc = %v MB\nTotalAlloc = %v MB\nSys = %v MB\nNumGC = %v\n\n", m.Alloc/1024/1024, m.TotalAlloc/1024/1024, m.Sys/1024/1024, m.NumGC)
			}
		}
	}()

	// Exit cleanly on a sigterm
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Turn on CPU profiling based on configuration
	if cfg.CPUProfile {
		f, err := os.Create("folly.prof")
		if err != nil {
			log.Println(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	router := httprouter.New()
	router.Handle("GET", "/work", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		work()
	})

	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
		if err != nil {
			log.Printf("Error in http server: %v", err)
		}
	}()

	<-c
	log.Println("Your folly is complete")
}

// Perform memory and CPU intensive work.
func work() {
	mu := sync.RWMutex{}

	// Only chew on the CPU for a short duration
	stop := false
	cpuTime := rand.Intn(cfg.DurationMax - cfg.DurationMin)
	go func() {
		time.Sleep(time.Duration(cpuTime) * time.Millisecond)
		mu.Lock()
		stop = true
		mu.Unlock()
	}()

	// Allocate memory, and chew up cycles
	arr := make([]byte, cfg.Memory)
	for {
		for i := 0; i < cfg.Memory; i++ {
			arr[i] = byte(i + 1)
		}
		mu.RLock()
		if stop {
			mu.RUnlock()
			break
		}
		mu.RUnlock()
	}

	cc <- 1
}
