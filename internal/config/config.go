package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds application configuration values loaded from the environment.
type Config struct {
	AppPort        string
	DBURL          string
	JWTSecret      string
	JWTExpiryHours int
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	jwtExpiryHours := 24
	if value := os.Getenv("JWT_EXPIRY_HOURS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
		}
		jwtExpiryHours = parsed
	}

	// Return a literal struct. This stays on the stack and
	// avoids unnecessary heap allocation.
	return Config{
		AppPort:        appPort,
		DBURL:          dbURL,
		JWTSecret:      jwtSecret,
		JWTExpiryHours: jwtExpiryHours,
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
