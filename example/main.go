package main

import (
	"fmt"
	"log"
	"time"

	"github.com/expo21xx/nxconfig"
)

type config struct {
	Host     string
	Port     uint16
	PGConfig pgconfig `name:"PG"`
}

type pgconfig struct {
	Host     string `default:"localhost"`
	Port     uint16 `default:"5432"`
	Username string `usage:"postgres username"`
	Password string
	Timeout  time.Duration `name:"connection-timeout"`
}

func main() {
	var cfg config

	err := nxconfig.Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", cfg)
}
