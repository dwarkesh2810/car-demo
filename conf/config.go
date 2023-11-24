package conf

import "github.com/spf13/viper"

type Config struct {
	JwtSecret        string `mapstructure:"SECRET"`
	TwilioAccountSID string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken  string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioServiceSID string `mapstructure:"TWILIO_SERVICE_SID"`
	From             string `mapstructure:"FROM"`
	Password         string `mapstructure:"PASSWORD"`
	SmtpHost         string `mapstructure:"SMTP_HOST"`
	SmtpPort         string `mapstructure:"SMTP_PORT"`
	BaseStoragePath  string `mapstructure:"BASE_STORAGE_PATH"`
	MailSubject      string `mapstructure:"SUBJECT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
