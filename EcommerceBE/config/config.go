package config

import "os"

type Config struct {
	Port string
	MongoURI string
	MongoDB string
	Clouldinary string
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8005"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27015"),
		MongoDB: getEnv("MONGO_DB", "modular-monolith"),
		Clouldinary: getEnv("CLOUDINARY_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

