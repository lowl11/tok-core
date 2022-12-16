package definition

import (
	"github.com/lowl11/lazyconfig/confapi"
	"github.com/lowl11/lazylog/logapi"
	"os"
)

type Configuration struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`

	Database struct {
		Connection     string `json:"connection"`
		MaxConnections int    `json:"max_connections"`
		Lifetime       int    `json:"lifetime"`
	} `json:"database"`

	Redis struct {
		Address  string `json:"address"`
		Password string `json:"password"`
	} `json:"redis"`

	User struct {
		CryptKey string `json:"crypt_key"`
	} `json:"user"`

	Image struct {
		BasePath string `json:"base_path"`
	} `json:"image"`
}

var Config Configuration
var Logger logapi.ILogger

func Init() {
	Config = Configuration{}
	isProduction := os.Getenv("env") == "production"

	logger := logapi.New().File("info", "logs")

	if err := confapi.Read(&Config, isProduction); err != nil {
		logger.Fatal(err, "Reading config error")
	}

	Logger = logger
	initServer()
}
