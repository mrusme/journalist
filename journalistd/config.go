package journalistd

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
  Admin                  struct {
    Username             string
    Password             string
  }

  Database               struct {
    Type                 string
    Connection           string
  }

  Server                 struct {
    BindIP               string
    Port                 string
    Endpoint             struct {
      Web                string
    }
  }
}

func Cfg() (Config, error) {
  viper.SetDefault("Admin.Username", "admin")
  viper.SetDefault("Admin.Password", "admin")

  viper.SetDefault("Database.Type", "sqlite3")
  viper.SetDefault("Database.Connection", "file:ent?mode=memory&cache=shared&_fk=1")

  viper.SetDefault("Server.BindIP", "127.0.0.1")
  viper.SetDefault("Server.Port", "8000")
  viper.SetDefault("Server.Endpoint.Api", "http://127.0.0.1:8000/api")
  viper.SetDefault("Server.Endpoint.Web", "http://127.0.0.1:8000/web")

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

  return config, nil
}

