package pkg

import (
	"os"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)


type DBconfig struct {
	Host string
	Port string
	Username string
	Password string
	DBname string
	SSLMode string
}

func LoadEnv() error {
	return godotenv.Load()
}

func GetDBconfig() DBconfig{
	return DBconfig{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBname: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
}

func GetPort() string {
	return viper.GetString("port")
}


func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

