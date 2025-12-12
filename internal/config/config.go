package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Email    EmailConfig
	Auth     AuthConfig
	Logging  LoggingConfig
	Cookie   CookieConfig
}

type ServerConfig struct {
	Port            string
	Environment     string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	PostgresURL string
	MongoDBURL  string
	RedisURL    string
	RedisDB     int
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

type OAuthConfig struct {
	Google   OAuthProviderConfig
	Facebook OAuthProviderConfig
}

type OAuthProviderConfig struct {
	Enabled      bool
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type LoggingConfig struct {
	Level         string
	SensitiveKeys []string
	IgnoredKeys   []string
}

type CookieConfig struct {
	Domain   string
	Secure   bool
	SameSite string
	HTTPOnly bool
}

type EmailConfig struct {
	ServiceEnabled bool
	FromAddress    string
	FromName       string
	FrontendURL    string
}

type AuthConfig struct {
	EmailPasswordEnabled    bool
	RequireEmailVerification bool
	VerificationTokenExpiry  time.Duration
	ResetTokenExpiry         time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (for local development)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Environment:     getEnv("ENVIRONMENT", "development"),
			ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			PostgresURL: getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/rest_api_db?sslmode=disable"),
			MongoDBURL:  getEnv("MONGODB_URL", "mongodb://admin:admin@localhost:27017"),
			RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
			RedisDB:     getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
			AccessExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			Issuer:        getEnv("JWT_ISSUER", "rest-api-service"),
		},
		OAuth: OAuthConfig{
			Google: OAuthProviderConfig{
				Enabled:      getEnvBool("OAUTH_GOOGLE_ENABLED", false),
				ClientID:     getEnv("OAUTH_GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("OAUTH_GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("OAUTH_GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			},
			Facebook: OAuthProviderConfig{
				Enabled:      getEnvBool("OAUTH_FACEBOOK_ENABLED", false),
				ClientID:     getEnv("OAUTH_FACEBOOK_CLIENT_ID", ""),
				ClientSecret: getEnv("OAUTH_FACEBOOK_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("OAUTH_FACEBOOK_REDIRECT_URL", "http://localhost:8080/auth/facebook/callback"),
			},
		},
		Logging: LoggingConfig{
			Level:         getEnv("LOG_LEVEL", "info"),
			SensitiveKeys: getEnvArray("LOG_SENSITIVE_KEYS", []string{"password", "hashed_password", "token", "secret", "access_token", "refresh_token"}),
			IgnoredKeys:   getEnvArray("LOG_IGNORED_KEYS", []string{}),
		},
		Cookie: CookieConfig{
			Domain:   getEnv("COOKIE_DOMAIN", "localhost"),
			Secure:   getEnvBool("COOKIE_SECURE", false),
			SameSite: getEnv("COOKIE_SAME_SITE", "Lax"),
			HTTPOnly: getEnvBool("COOKIE_HTTP_ONLY", true),
		},
		Email: EmailConfig{
			ServiceEnabled: getEnvBool("EMAIL_SERVICE_ENABLED", false),
			FromAddress:    getEnv("EMAIL_FROM_ADDRESS", "noreply@yourapp.com"),
			FromName:       getEnv("EMAIL_FROM_NAME", "Your App Name"),
			FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
		Auth: AuthConfig{
			EmailPasswordEnabled:     getEnvBool("EMAIL_PASSWORD_AUTH_ENABLED", false),
			RequireEmailVerification: getEnvBool("REQUIRE_EMAIL_VERIFICATION", true),
			VerificationTokenExpiry:  getEnvDuration("VERIFICATION_TOKEN_EXPIRY", 24*time.Hour),
			ResetTokenExpiry:         getEnvDuration("RESET_TOKEN_EXPIRY", 1*time.Hour),
		},
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}
	if c.OAuth.Google.Enabled && (c.OAuth.Google.ClientID == "" || c.OAuth.Google.ClientSecret == "") {
		return fmt.Errorf("Google OAuth is enabled but credentials are missing")
	}
	if c.OAuth.Facebook.Enabled && (c.OAuth.Facebook.ClientID == "" || c.OAuth.Facebook.ClientSecret == "") {
		return fmt.Errorf("Facebook OAuth is enabled but credentials are missing")
	}
	return nil
}

// Helper functions for reading environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvArray(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
