package common

import "os"

type Config struct {
	Port string `json:"PORT"`

	DatabaseHost string `json:"DATABASE_HOST"`
	DatabasePort string `json:"DATABASE_PORT"`
	DatabaseUser string `json:"DATABASE_USER"`
	DatabasePass string `json:"DATABASE_PASS"`
	DatabaseName string `json:"DATABASE_NAME"`
}

func GetConfig() *Config {
	return &Config{
		Port:         ThreeTermString(len(os.Getenv("PORT")) > 0, os.Getenv("PORT"), "50051"),
		DatabaseHost: ThreeTermString(len(os.Getenv("DATABASE_HOST")) > 0, os.Getenv("DATABASE_HOST"), "localhost"),
		DatabasePort: ThreeTermString(len(os.Getenv("DATABASE_PORT")) > 0, os.Getenv("DATABASE_PORT"), "5432"),
		DatabaseUser: ThreeTermString(len(os.Getenv("DATABASE_USER")) > 0, os.Getenv("DATABASE_USER"), "root"),
		DatabasePass: ThreeTermString(len(os.Getenv("DATABASE_PASS")) > 0, os.Getenv("DATABASE_PASS"), "1234"),
		DatabaseName: ThreeTermString(len(os.Getenv("DATABASE_NAME")) > 0, os.Getenv("DATABASE_NAME"), "remoter"),
	}
}
