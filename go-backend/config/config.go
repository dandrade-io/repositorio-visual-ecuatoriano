package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"time"

	"github.com/joho/godotenv"
)

var (
	DBUser             string
	DBPassword         string
	DBHost             string
	DBPort             string
	DBName             string
	EmbedServiceURL    string
	EmbedServiceApiKey string
	PublicBaseURL      string
	LLMBaseUrl         string
	LLMApiKey          string
	LLMModel           string
	KCIssuerURL        string
	KCBaseURL          string
	KCRealm            string
	CAPath             string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envPath := filepath.Join(filepath.Dir(wd), ".env")
	_ = godotenv.Load(envPath)

	DBUser = os.Getenv("POSTGRES_USER")
	DBPassword = os.Getenv("POSTGRES_PASSWORD")
	DBHost = os.Getenv("POSTGRES_HOST")
	DBPort = os.Getenv("POSTGRES_PORT")
	DBName = os.Getenv("POSTGRES_DB")

	EmbedServiceURL = os.Getenv("EMBED_SERVICE_BASE_URL")
	EmbedServiceApiKey = os.Getenv("EMBED_SERVICE_API_KEY")
	PublicBaseURL = os.Getenv("BACKEND_BASE_URL")

	LLMBaseUrl = os.Getenv("LLM_BASE_URL")
	LLMApiKey = os.Getenv("LLM_API_KEY")
	LLMModel = os.Getenv("LLM_NAME")

	KCIssuerURL = os.Getenv("KC_ISSUER_URL")
	KCBaseURL = os.Getenv("KC_BASE_URL")
	KCRealm = os.Getenv("KC_REALM")

	CAPath = os.Getenv("CA_PATH")
}

func ClientWithPEM(caPath string) (*http.Client, error) {
	if caPath == "" {
		return nil, fmt.Errorf("no certifica file found")
	}

	pem, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		pool = x509.NewCertPool()
	}

	if !pool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("no valid certificate found in CA file")
	}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}

	return &http.Client{
		Transport: tr,
		Timeout:   60 * time.Second,
	}, nil
}
