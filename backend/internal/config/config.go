package config

import (
	"log"
	"sync"
	"time"

	"github.com/jinzhu/now"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env             string `envconfig:"ENV"`
	LogLevel        string `envconfig:"LOG_LEVEL"`
	MongoURL        string `envconfig:"MONGODB_URL"`
	MongoDB         string `envconfig:"MONGODB_DATABASE"`
	TemporalURL     string `envconfig:"TEMPORAL_URL"`
	HTTPAddr        string `envconfig:"HTTP_ADDR"`
	JWTSecret       string `envconfig:"JWT_SECRET"`
	ApiRole         string `envconfig:"API_ROLE"`
	WorkerRole      string `envconfig:"WORKER_ROLE"`
	NotifierRole    string `envconfig:"NOTIFIER_ROLE"`
	TonListenerRole string `envconfig:"TON_LISTENER_ROLE"`
	TonWorkerRole   string `envconfig:"TON_WORKER_ROLE"`
	TgBotApiKey     string `envconfig:"TG_BOT_API_KEY"`
	CgBotApiKey     string `envconfig:"CG_BOT_API_KEY"`
	ProductRole     string `envconfig:"PRODUCT_ROLE"`
	AutofarmRole    string `envconfig:"AUTOFARM_ROLE"`
	TonNetwork      string `envconfig:"TON_NETWORK"`
	TonWallet       string `envconfig:"TON_WALLET"`

	PaymentGatewayURL             string `envconfig:"PAYMENT_GATEWAY_URL"`
	PaymentGatewayLogin           string `envconfig:"PAYMENT_GATEWAY_LOGIN"`
	PaymentGatewayKey             string `envconfig:"PAYMENT_GATEWAY_KEY"`
	PaymentGatewayPaymentSystemID string `envconfig:"PAYMENT_GATEWAY_PAYMENT_SYSTEM_ID"`
}

var (
	config    Config
	nowConfig now.Config
	once      sync.Once
)

// Get reads config from environment
func Get() *Config {
	once.Do(func() {

		//var cfgPath string
		//flag.StringVar(&cfgPath, "e", ".env", "environment variables (path to .env file)")
		//flag.Parse()

		//_, err := os.Stat(cfgPath)
		//if errors.Is(err, os.ErrNotExist) {
		//	log.Println("config file not found")
		//}
		//err = godotenv.Load(cfgPath) // load .env file
		//if err != nil {
		//	log.Println(err)
		//}

		err := envconfig.Process("", &config)
		if err != nil {
			log.Println(err)
		}

		log.Println(config)
	})
	return &config
}

func GetNowConfig() *now.Config {
	once.Do(func() {
		nowConfig = now.Config{
			WeekStartDay: time.Monday,
			TimeLocation: time.UTC,
		}
	})
	return &nowConfig
}
