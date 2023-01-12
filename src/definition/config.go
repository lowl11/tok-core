package definition

import (
	"github.com/lowl11/lazyconfig/confapi"
	"github.com/lowl11/lazylog/logapi"
	"os"
)

type Configuration struct {
	Environment string
	Primary     bool

	Server struct {
		Port struct {
			Primary   string `json:"primary"`
			Secondary string `json:"secondary"`
		} `json:"port"`
	} `json:"server"`

	Database struct {
		Connection     string `json:"connection"`
		MaxConnections int    `json:"max_connections"`
		Lifetime       int    `json:"lifetime"`
	} `json:"database"`

	Mongo struct {
		Connection string `json:"connection"`
	} `json:"mongo"`

	Redis struct {
		Address  string `json:"address"`
		Password string `json:"password"`
	} `json:"redis"`

	User struct {
		CryptKey string `json:"crypt_key"`
	} `json:"user"`

	Image struct {
		BasePath string `json:"base"`
	} `json:"image"`
}

var Config Configuration
var Logger logapi.ILogger

func Init() {
	Config = Configuration{}

	// определение окружения (прод или нет)
	Config.Environment = os.Getenv("env")
	isProduction := Config.Environment == "production"

	// определение порта (primary, secondary)
	if len(os.Args) > 1 && os.Args[1] == "secondary" {
		Config.Primary = false
	} else {
		Config.Primary = true
	}

	// создание логгера
	logger := logapi.New().File("info", "logs")

	// чтение конфигов
	if err := confapi.Read(&Config, isProduction); err != nil {
		logger.Fatal(err, "Reading config error")
	}

	Logger = logger

	// создание объекта сервера
	initServer()
}
