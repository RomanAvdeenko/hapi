package main

import (
	"flag"
	"hapi/internal/app/apiserver"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	CONFIG_PATH = "configs/apiserver.yaml"
)

func main() {
	configPath := flag.String("config-path", CONFIG_PATH, "path to config file")
	flag.Parse()

	config := apiserver.NewConfig()
	err := cleanenv.ReadConfig(*configPath, config)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(apiserver.Start(config))
}
