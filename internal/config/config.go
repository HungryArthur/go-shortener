package config

import (
	"flag"
	"os"
	"strings"
)

var FlagRunAddr string

var FlagBaseShortenedURLAddr string

func Load() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagBaseShortenedURLAddr, "b", "", "base address of the resulting shortened URL")

	flagRunAddr := os.Getenv("SERVER_ADDRESS")
	if flagRunAddr != "" {
		FlagRunAddr = flagRunAddr
	}

	flagBaseShortenedURLAddr := os.Getenv("BASE_URL")
	if flagBaseShortenedURLAddr != "" {
		FlagBaseShortenedURLAddr = flagBaseShortenedURLAddr
	}

	flag.Parse()

	// FlagBaseShortenedURLAddr = strings.TrimPrefix(FlagBaseShortenedURLAddr, "http://localhost" + FlagRunAddr)
	FlagBaseShortenedURLAddr = strings.TrimSuffix(FlagBaseShortenedURLAddr, "/")
}
