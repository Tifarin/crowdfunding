package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	DB        DBConfig
	JWTConfig JWTConfig
	Mailer    Mailer
}

type ServerConfig struct {
	Addr            string
	WebsocketAddr   string
	WriteTimeout    int
	ReadTimeout     int
	GraceFulTimeout int
	Registration    bool
	ForceTLS        bool
}

type DBConfig struct {
	Name            string
	Host            string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
}

type JWTConfig struct {
	Issuer            string
	Secret            string
	TokenLifeTimeHour int
}

type Mailer struct {
	Server     string
	Port       int
	Username   string
	Password   string
	UseTls     bool
	Sender     string
	MaxAttempt int
}


func InitConfig() Config {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")

	var configuration Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return configuration
}
