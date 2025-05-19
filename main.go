package main

import (
	"fmt"

	"github.com/oleshko-g/gatorcli/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, _ := config.Read()
	cfg.SetUser("gena")
	cfg, _ = config.Read()
	fmt.Printf("%#v", cfg)
}
