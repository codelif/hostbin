package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadDotEnv(paths ...string) error {
	if len(paths) == 0 {
		paths = []string{".env"}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		values, err := godotenv.Read(path)
		if err != nil {
			return err
		}

		for key, value := range values {
			if _, exists := os.LookupEnv(key); !exists {
				if err := os.Setenv(key, value); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
