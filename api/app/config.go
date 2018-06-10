package app

import (
	"github.com/RaniSputnik/ok/api/auth"
	"github.com/RaniSputnik/ok/api/store"
	"github.com/RaniSputnik/ok/api/store/inmemory"
)

type Config struct {
	Addr    string
	Store   store.Store
	AuthSvc auth.Service
}

func (c Config) withSensibleDefaults() Config {
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.Store == nil {
		c.Store = inmemory.New()
	}
	if c.AuthSvc == nil {
		c.AuthSvc = auth.NewHMAC([]byte("todo-this-is-a-default-secret"))
	}
	return c
}
