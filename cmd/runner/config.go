package main

// Config is the configuration struct
type Config struct {
	Address     string `envconfig:"ADDRESS" default:"localhost"`
	Port        int    `envconfig:"PORT" default:"8080"`
	Concurrency int    `envconfig:"CONCURRENCY" default:"100"`
	Iterations  int    `envconfig:"ITERATIONS" default:"10"`
}
