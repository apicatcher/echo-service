package config

import (
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/apicatcher/echo-service/pkg/util"
)

var C = new(Config)

type Config struct {
	Server ServerConfig `toml:"server"`
}

type ServerConfig struct {
	Port int `toml:"port"`
}

func init() {
	split := "/"
	if runtime.GOOS == "windows" {
		split = "\\"
	}
	var err error
	dir, _ := os.Getwd()
	if strings.Contains(dir, split+"internal"+split) {
		index := strings.Index(dir, split+"internal"+split)
		dir = dir[:index]
	}
	c, err := os.Open(dir + split + "config.toml")
	if err != nil {
		util.AbnormalExit(err)
	}
	data, err := io.ReadAll(c)
	if err != nil {
		util.AbnormalExit(err)
	}
	if _, err = toml.Decode(string(data), C); err != nil {
		util.AbnormalExit(err)
	}
}
