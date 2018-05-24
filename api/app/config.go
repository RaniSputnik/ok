package app

import "github.com/RaniSputnik/ok/api/store"

type Config struct {
	Addr  string
	Store store.Store
}

func (c Config) withSensibleDefaults() Config {
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.Store == nil {
		c.Store = store.NewInMemory()
	}
	return c
}