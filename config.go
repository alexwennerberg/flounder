package main

import (
	"github.com/BurntSushi/toml"
	"path/filepath"
)

const HiddenFolder = ".hidden"
const GemlogFolder = "gemlog"

type Config struct {
	FilesDirectory     string
	TemplatesDirectory string
	Host               string
	HttpPort           int
	SiteTitle          string
	Debug              bool
	SecretKey          string
	DBFile             string
	AnalyticsDBFile    string
	LogFile            string
	GeminiCertStore    string
	CookieStoreKey     string
	OkExtensions       []string
	MaxFileBytes       int
	MaxFilesPerUser    int
	MaxUserBytes       int64
	SMTPServer         string
	SMTPUsername       string
	SMTPPassword       string
	EnableSFTP         bool
	HostKeyPath        string
}

func getConfig(filename string) (Config, error) {
	var config Config
	// Attempt to overwrite defaults from file
	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		return config, err
	}
	// Workaround for how some of my path fns are written
	config.FilesDirectory, _ = filepath.Abs(config.FilesDirectory)
	return config, nil
}
