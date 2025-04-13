package initializers

import (
	"fmt"
	"os"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

var (
	localScrapper *twitterscraper.Scraper
)

func InitializeScrapper() error {
	initOnce.Do(func() {
		// Get ENV variables
		authToken := os.Getenv("AUTH_TOKEN")
		csrfToken := os.Getenv("CSRF_TOKEN")

		// Set up Twitter scraper
		scraper := twitterscraper.New()
		scraper.SetAuthToken(twitterscraper.AuthToken{Token: authToken, CSRFToken: csrfToken})
		if !scraper.IsLoggedIn() {
			initError = fmt.Errorf("Invalid AuthToken or CSRFToken")
			return
		}

		localScrapper = scraper
	})
	return initError
}

// Get the local scrapper
func GetScrapper() (*twitterscraper.Scraper, error) {
	if err := InitializeScrapper(); err != nil {
		return nil, err
	}

	mutex.RLock()
	defer mutex.RUnlock()
	return localScrapper, nil
}
