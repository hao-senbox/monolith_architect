package config

import "os"

type Config struct {
	Port        string
	MongoURI    string
	MongoDB     string
	Clouldinary string
	VNPayConfig VNPayConfig
}

type VNPayConfig struct {
	TmnCode    string
	HashSecret string
	PaymentUrl string
	Version    string
	Command    string
	CurrCode   string
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8005"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27015"),
		MongoDB:     getEnv("MONGO_DB", "modular-monolith"),
		Clouldinary: getEnv("CLOUDINARY_URL", ""),
		VNPayConfig: VNPayConfig{
			TmnCode:    getEnv("VN_PAY_TMN_CODE", ""),
			HashSecret: getEnv("VN_PAY_HASH_SECRET", ""),
			PaymentUrl: getEnv("VN_PAY_PAYMENT_URL", ""),
			Version:    getEnv("VN_PAY_VERSION", "2.1.0"),
			Command:    getEnv("VN_PAY_COMMAND", "pay"),
			CurrCode:   getEnv("VN_PAY_CURR_CODE", "VND"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
