package initializers

import (
	"github.com/joho/godotenv"
)

func InitializeEVN() error {
	initOnce.Do(func() {
		err := godotenv.Load()
		if err != nil {
			initError = err
		}
	})
	return initError
}
