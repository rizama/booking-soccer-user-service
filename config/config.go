package config

import (
	"os"
	"user-service/common/utils"

	"github.com/sirupsen/logrus"
	_ "github.com/spf13/viper/remote"
)

var Config AppConfig

type AppConfig struct {
	Port                  int            `json:"port"`
	AppName               string         `json:"appName"`
	AppEnv                string         `json:"appEnv"`
	SignatureKey          string         `json:"signatureKey"`
	Database              DatabaseConfig `json:"database"`
	RateLimiterRequest    int            `json:"rateLimiterRequest"`
	RateLimiterTimeSecond int            `json:"rateLimiterTimeSecond"`
	JwtSecretKey          string         `json:"jwtSecretKey"`
	JwtExpirationTime     int            `json:"jwtExpirationTime"`
}

type DatabaseConfig struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	User                  string `json:"user"`
	Name                  string `json:"name"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	MaxOpenConnection     int    `json:"maxOpenConnection"`
	MaxLifetimeConnection int    `json:"maxLifetimeConnection"`
	MaxIdleConnection     int    `json:"maxIdleConnection"`
	MaxIdleTime           int    `json:"maxIdleTime"`
}

func Init() {
	err := utils.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("Failed load config json local %v", err)
		err = utils.BindFromConsul(&Config, os.Getenv("CONSUL_HTPP_URL"), os.Getenv("CONSUL_CONFIG_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
