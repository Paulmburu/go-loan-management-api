package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AppPort string
	DBURL   string
}

// Load reads environment variables and returns a Config VALUE.
// We return a value (not a pointer) because Config is small,
// read-only, and should be safely copied across the app.
func LoadConfig() (Config, error) {
	loadDotEnv()

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // default port
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	// Return a literal struct. This stays on the stack and
	// avoids unnecessary heap allocation.
	return Config{
		AppPort: appPort,
		DBURL:   dbURL,
	}, nil
}

func loadDotEnv() {
	for _, path := range candidateDotEnvPaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			key, value, found := strings.Cut(line, "=")
			if !found {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			value = strings.Trim(value, `"'`)

			if key == "" {
				continue
			}

			if os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}

		return
	}
}

func candidateDotEnvPaths() []string {
	paths := []string{".env"}

	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths,
			filepath.Join(cwd, ".env"),
			filepath.Join(cwd, "..", ".env"),
			filepath.Join(cwd, "..", "..", ".env"),
		)
	}

	return paths
}
