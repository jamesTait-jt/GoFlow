package config

import (
	"flag"
	"fmt"
)

var defaultNumWorkers = 5

var defaultBrokerType = "redis"

var supportedBrokerTypes = []string{"redis"}

type Config struct {
	NumWorkers   int
	HandlersPath string
	BrokerType   string
	BrokerAddr   string
}

func LoadConfigFromFlags() *Config {
	c := &Config{}

	flag.IntVar(&c.NumWorkers, "num-workers", defaultNumWorkers, "Number of workers in the pool")
	flag.StringVar(&c.HandlersPath, "handlers-path", "", "Path to the location of the handler plugins")
	enumFlag(&c.BrokerType, "broker-type", supportedBrokerTypes, "Type of task broker (e.g. 'redis')")
	flag.StringVar(&c.BrokerAddr, "broker-addr", "", "Broker address (e.g., Redis address)")

	flag.Parse()

	return c
}

func enumFlag(target *string, name string, allowed []string, usage string) {
	flag.Func(name, usage, func(flagValue string) error {
		if flagValue == "" {
			*target = defaultBrokerType

			return nil
		}

		for _, allowedValue := range allowed {
			if flagValue == allowedValue {
				*target = flagValue
				return nil
			}
		}

		return fmt.Errorf("must be one of %v", allowed)
	})
}
