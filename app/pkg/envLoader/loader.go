package envLoader

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(paths ...string) {
	for _, path := range paths {
		if err := godotenv.Load(path); err != nil {
			log.Printf("No .env file found at %s (skipping): %v\n", path, err)
		}
	}
}
