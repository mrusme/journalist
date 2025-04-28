package lib

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Debug string

	Admin struct {
		Username string
		Password string
	}

	Database struct {
		Type       string
		Connection string
	}

	Server struct {
		BindIP   string
		Port     string
		Endpoint struct {
			Api string
			Web string
		}
	}

	Feeds struct {
		AutoRefresh string
	}
}

func Cfg() (Config, error) {
	viper.SetDefault("Debug", "false")

	viper.SetDefault("Admin.Username", "admin")
	viper.SetDefault("Admin.Password", "admin")

	viper.SetDefault("Database.Type", "sqlite3")
	viper.SetDefault("Database.Connection", "file:ent?mode=memory&cache=shared&_fk=1")

	viper.SetDefault("Server.BindIP", "127.0.0.1")
	viper.SetDefault("Server.Port", "8000")
	viper.SetDefault("Server.Endpoint.Api", "http://127.0.0.1:8000/api")
	viper.SetDefault("Server.Endpoint.Web", "http://127.0.0.1:8000/web")

	viper.SetDefault("Feeds.AutoRefresh", "900")

	viper.SetConfigName("journalist.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$XDG_CONFIG_HOME/")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("journalist")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	config = *ParseDatabaseURL(&config)

	return config, nil
}

func ParseDatabaseURL(config *Config) *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return config
	}

	dbURL, err := url.Parse(databaseURL)
	if err != nil {
		return config
	}

	host, port, _ := net.SplitHostPort(dbURL.Host)
	dbname := strings.TrimLeft(dbURL.Path, "/")
	user := dbURL.User.Username()
	password, _ := dbURL.User.Password()

	switch dbURL.Scheme {
	case "postgresql", "postgres":
		if port == "" {
			port = "5432"
		}
		config.Database.Type = "postgres"
		config.Database.Connection = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s",
			host, port, dbname, user, password,
		)
	case "mysql":
		if port == "" {
			port = "3306"
		}
		config.Database.Type = "mysql"
		config.Database.Connection = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=True",
			user, password, host, port, dbname,
		)
	}

	return config
}
