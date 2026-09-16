package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App     AppConfig
	DB      DatabaseConfig
	JWT     JWTConfig
	Payment PaymentConfig
	Upload  UploadConfig
	Rate    RateLimitConfig
}

type AppConfig struct {
	Name  string
	Port  string
	Env   string
	Debug bool
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type JWTConfig struct {
	Secret            string
	RefreshSecret     string
	Expiration        int
	RefreshExpiration int
}

type PaymentConfig struct {
	MidtransServerKey  string
	MidtransClientKey  string
	MidtransMerchantID string
	MidtransProduction bool
	XenditSecretKey    string
	XenditPublicKey    string
	StripeSecretKey    string
	StripeWebhook      string
}

type UploadConfig struct {
	Dir          string
	MaxUploadSize int64
}

type RateLimitConfig struct {
	Limit      int
	Expiration int
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:  getEnv("APP_NAME", "KPM Academy"),
			Port:  getEnv("APP_PORT", "3000"),
			Env:   getEnv("APP_ENV", "development"),
			Debug: getEnvBool("APP_DEBUG", true),
		},
		DB: DatabaseConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "kpm_academy"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 5),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "default-secret"),
			RefreshSecret:     getEnv("JWT_REFRESH_SECRET", "default-refresh-secret"),
			Expiration:        getEnvInt("JWT_EXPIRATION", 3600),
			RefreshExpiration: getEnvInt("JWT_REFRESH_EXPIRATION", 604800),
		},
		Payment: PaymentConfig{
			MidtransServerKey:  getEnv("MIDTRANS_SERVER_KEY", ""),
			MidtransClientKey:  getEnv("MIDTRANS_CLIENT_KEY", ""),
			MidtransMerchantID: getEnv("MIDTRANS_MERCHANT_ID", ""),
			MidtransProduction: getEnvBool("MIDTRANS_IS_PRODUCTION", false),
			XenditSecretKey:    getEnv("XENDIT_SECRET_KEY", ""),
			XenditPublicKey:    getEnv("XENDIT_PUBLIC_KEY", ""),
			StripeSecretKey:    getEnv("STRIPE_SECRET_KEY", ""),
			StripeWebhook:      getEnv("STRIPE_WEBHOOK_SECRET", ""),
		},
		Upload: UploadConfig{
			Dir:          getEnv("UPLOAD_DIR", "./uploads"),
			MaxUploadSize: getEnvInt64("MAX_UPLOAD_SIZE", 10485760),
		},
		Rate: RateLimitConfig{
			Limit:      getEnvInt("RATE_LIMIT", 100),
			Expiration: getEnvInt("RATE_LIMIT_EXPIRATION", 60),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func getEnvInt64(key string, defaultVal int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return i
}
