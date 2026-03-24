package config

import (
	"flag"
	"os"
	"strings"
)

var (
	FlagRunAddr string
	FlagBaseShortenedURLAddr string
	FileStoragePath string
	FileDatabaseDSN string
)

func Load() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagBaseShortenedURLAddr, "b", "", "base address of the resulting shortened URL")	
	flag.StringVar(&FileStoragePath, "f", "", "json file for storing url data")
	flag.StringVar(&FileDatabaseDSN, "d", "", "database DSN (postgres://...)")

	flagRunAddr := os.Getenv("SERVER_ADDRESS")
	if flagRunAddr != "" {
		FlagRunAddr = flagRunAddr
	}

	flagBaseShortenedURLAddr := os.Getenv("BASE_URL")
	if flagBaseShortenedURLAddr != "" {
		FlagBaseShortenedURLAddr = flagBaseShortenedURLAddr
	}

	flag.Parse()

	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv != "" {
		FileStoragePath = fileStoragePathEnv
	}
	
	fileDatabaseDSNEnv := os.Getenv("DATABASE_DSN")
	if fileDatabaseDSNEnv != "" {
		FileDatabaseDSN = fileDatabaseDSNEnv
	}

	// FlagBaseShortenedURLAddr = strings.TrimPrefix(FlagBaseShortenedURLAddr, "http://localhost" + FlagRunAddr)
	FlagBaseShortenedURLAddr = strings.TrimSuffix(FlagBaseShortenedURLAddr, "/")
}
