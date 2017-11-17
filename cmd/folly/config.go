package main

// Config is the configuration struct
type Config struct {
	Port        int  `envconfig:"PORT" default:"8080"`
	DurationMax int  `envconfig:"DURATION_MAX" default:"800"`
	DurationMin int  `envconfig:"DURATION_MIN" default:"500"`
	Memory      int  `envconfig:"MEMORY" default:"1024"`
	CPUProfile  bool `envconfig:"CPU_PROFILE" default:"false"`
}
