package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	DBUser          string
	DBPassword      string
	DBHost          string
	DBPort          string
	DBName          string
	EmbedServiceURL string
	PublicBaseURL   string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envPath := filepath.Join(filepath.Dir(wd), ".env")

	// Optional: in containers the variables come from the environment.
	_ = godotenv.Load(envPath)

	DBUser = os.Getenv("POSTGRES_USER")
	DBPassword = os.Getenv("POSTGRES_PASSWORD")
	DBHost = os.Getenv("POSTGRES_HOST")
	DBPort = os.Getenv("POSTGRES_PORT")
	DBName = os.Getenv("POSTGRES_DB")
	EmbedServiceURL = os.Getenv("EMBED_SERVICE_BASE_URL")
	PublicBaseURL = os.Getenv("BACKEND_BASE_URL")
}
