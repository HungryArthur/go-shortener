package config

import (
	"flag"
	"strings"
)

var FlagRunAddr string

var FlagBaseShortenedURLAddr string


func Load() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	
	flag.StringVar(&FlagBaseShortenedURLAddr, "b", "", "base address of the resulting shortened URL")
	
	flag.Parse()
	
	// FlagBaseShortenedURLAddr = strings.TrimPrefix(FlagBaseShortenedURLAddr, "http://localhost" + FlagRunAddr)
	FlagBaseShortenedURLAddr = strings.TrimSuffix(FlagBaseShortenedURLAddr, "/")
}