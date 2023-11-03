package main

import (
	"fmt"
	"os"
	"github.com/jetsadawwts/go-restapi/config"
)


func envPath() string {
	if len(os.Args) == 1 {
		return	".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())
	fmt.Println(cfg.App())
	fmt.Println(cfg.Db())
	fmt.Println(cfg.Jwt())
}