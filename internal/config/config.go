package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	User     string
	Password string
	OldUrl   string
	NewUrl   string
	Token    string
}

var (
	cfg = &Config{}
)

func New() (c *Config, err error) {
	err = loadEnv(cfg)
	if err == nil {
		c = cfg
		return
	}

	slog.Info("Env not contains Trying to get from file", "msg", err)

	err = godotenv.Load(".env")
	if err != nil {
		slog.Error("Can't load environments", "err", err)
		return
	}

	err = loadEnv(cfg)
	if err != nil {
		slog.Error("Env not contains needed vars: %s", "err", err)
		return
	}

	c = cfg
	return
}

func loadEnv(cfg *Config) (err error) {
	vars := []string{"User", "Password", "OldUrl", "NewUrl"}
	missed := make([]string, 0, len(vars))

	for _, vr := range vars {
		v, ok := os.LookupEnv(vr)
		if !ok {
			missed = append(missed, vr)
		} else {
			switch vr {
			case "User":
				cfg.User = v
			case "Password":
				cfg.Password = v
			case "OldUrl":
				cfg.OldUrl = v
			case "NewUrl":
				cfg.NewUrl = v
			}
		}
	}

	if len(missed) > 0 {
		return fmt.Errorf("missing environment vars: %s", strings.Join(missed, ", "))
	}

	return
}
