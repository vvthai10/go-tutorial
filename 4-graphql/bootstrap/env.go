package bootstrap

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 	string 			`mapstructure:"APP_ENV"`
	Port          			string 			`mapstructure:"PORT"`
	DBHost                 	string 			`mapstructure:"DB_HOST"`
	DBPort                 	string 			`mapstructure:"DB_PORT"`
	DBUser                 	string 			`mapstructure:"DB_USER"`
	DBPass                 	string 			`mapstructure:"DB_PASS"`
	DBName                 	string 			`mapstructure:"DB_NAME"`
	RabbitUrl               string 			`mapstructure:"RABBITMQ_URL"`
	RabbitWaitTime          string 			`mapstructure:"RABBITMQ_WAITTIME"`
	RabbitAttempts          string 			`mapstructure:"RABBITMQ_ATTEMPTS"`
	DBSSLMode              	string 			`mapstructure:"DB_SSL_MODE"`
	GinMode					string 			`mapstructure:"GIN_MODE"`
	LogLevel  				int    			`mapstructure:"LOG_LEVEL"`
	TokenSymmetricKey  		string    		`mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  	time.Duration   `mapstructure:"ACCESS_TOKEN_DURATION"`
	AuthorizationPayloadKey string 			`mapstructure:"AUTHORIZATION_PAYLOAD_KEY"`
	AuthorizationType 		string 			`mapstructure:"AUTHORIZATION_TYPE"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
