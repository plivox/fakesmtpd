package config

import (
	"strings"

	"github.com/imdario/mergo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server *Server    `mapstructure:"server"`
	Log    *LogConfig `mapstructure:"log"`
}

type Server struct {
	Address  string `mapstructure:"address"`
	Domain   string `mapstructure:"domain"`
	TLS      bool   `mapstructure:"tls"`
	Insecure bool   `mapstructure:"insecure"`
	CACert   string `mapstructure:"cacert"`
}

type LogConfig struct {
	Level   string            `mapstructure:"level"`
	Console *LogConfigConsole `mapstructure:"console"`
	File    *LogConfigFile    `mapstructure:"file"`
	Syslog  *LogConfigSyslog  `mapstructure:"syslog"`
}

type LogConfigConsole struct {
	Color bool `mapstructure:"color"`
	Json  bool `mapstructure:"json"`
}

type LogConfigFile struct {
	Path string `mapstructure:"path"`
	Json bool   `mapstructure:"json"`
}

type LogConfigSyslog struct{}

func (c *Config) setDefault() error {
	return mergo.Merge(c, &Config{
		Server: &Server{
			Address:  "127.0.0.1:2525",
			Domain:   "localhost",
			TLS:      false,
			Insecure: false,
			CACert:   "",
		},
		Log: &LogConfig{
			Level: "info",
			Console: &LogConfigConsole{
				Color: true,
				Json:  false,
			},
		},
	})
}

func NewConfig(configPath string) *Config {
	var config = &Config{}

	if err := config.setDefault(); err != nil {
		log.Fatal().Err(err).Msg("Unable to set default configuration")
	}

	if configPath != "" {
		viper.SetConfigFile(configPath)

	} else {
		viper.SetConfigName("fakesmtpd")
		viper.SetConfigType("yaml")

		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/fakesmtpd")
		viper.AddConfigPath("/etc/fakesmtpd/")
	}

	viper.SetEnvPrefix("fakesmtpd")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Msgf("Unable to read the configuration file: %s\n", viper.ConfigFileUsed())
	}

	// Workaround because viper does not treat env vars the same as other config
	for _, key := range viper.AllKeys() {
		viper.Set(key, viper.Get(key))
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("Unable to parse configuration")
	}

	return config
}
