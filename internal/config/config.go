package config

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	User         string
	Password     string
	OldUrl       string
	NewUrl       string
	Token        string
	MacAddresses []string
	Timeout      time.Duration
	MaxConn      int
}

var (
	max_conn = 5
	cfg      = &Config{}
	re       = regexp.MustCompile(`^([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})[:-]?([0-9A-Fa-f]{2})$`)
)

func New() (c *Config, err error) {
	err = cfg.loadEnv()
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

	err = cfg.loadEnv()
	if err != nil {
		slog.Error("Env not contains needed vars: %s", "err", err)
		return
	}

	cfg.parseFlags()

	c = cfg
	return
}

func (cfg *Config) loadEnv() error {
	vars := []string{"USERNAME", "PASSWORD", "OLD_URL", "NEW_URL", "MAX_CONN"}
	missed := make([]string, 0, len(vars))

	for _, vr := range vars {
		v, ok := os.LookupEnv(vr)
		if !ok {
			missed = append(missed, vr)
		} else {
			switch vr {
			case "USERNAME":
				cfg.User = v
			case "PASSWORD":
				cfg.Password = v
			case "OLD_URL":
				cfg.OldUrl = v
			case "NEW_URL":
				cfg.NewUrl = v
			case "MAX_CONN":
				num, err := strconv.Atoi(v)
				if err != nil {
					slog.Error("Can't parse MAX_CONN variable. Set to default value", "msg", max_conn)
					cfg.MaxConn = max_conn
				} else {
					cfg.MaxConn = num
				}
			}
		}
	}

	if len(missed) > 0 {
		return fmt.Errorf("missing environment vars: %s", strings.Join(missed, ", "))
	}

	return nil
}

func (cfg *Config) parseFlags() {
	var rawMacs string
	var path string
	var timeout time.Duration

	flag.StringVar(&rawMacs, "a", "", "Comma separated list of MAC addresses. "+
		"Do no use with -a arg together because -a will be overwrite with .yaml data even if file not exists or file is empty")
	flag.StringVar(&path, "p", "", "Path to .txt or .csv file with MAC addresses. "+
		"Do no use with -a arg together because -a will be overwrite with .yaml data even if file not exists or file is empty")
	flag.DurationVar(&timeout, "t", 60, "Timeout for responses")
	flag.Parse()

	cfg.Timeout = timeout * time.Second

	if path != "" {
		cfg.fromFile(path)
	}

	if len(rawMacs) == 0 {
		return
	}

	macs := strings.Split(rawMacs, ",")

	for _, mac := range macs {
		cfg.parseMacs(mac, re)
	}

	if len(cfg.MacAddresses) == 0 {
		slog.Info("No MAC addresses found")
		return
	}

	slog.Info("Found number of MAC", "msg", len(cfg.MacAddresses))
}

func parseMac(s string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(s)

	if len(matches) != 7 {
		slog.Error("Invalid MAC address format", "msg", s)
		return ""
	}

	return strings.ToLower(strings.Join(matches[1:], "-"))
}

func (cfg *Config) fromFile(p string) {
	file, err := os.Open(p)
	if err != nil {
		slog.Error("Error opening file", "err", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cfg.parseMacs(scanner.Text(), re)
	}
}

func (cfg *Config) parseMacs(m string, re *regexp.Regexp) {
	m = parseMac(m, re)
	if m != "" {
		cfg.MacAddresses = append(cfg.MacAddresses, m)
	}
}
