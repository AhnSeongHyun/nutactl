package config

import (
	"os"
)

type Config struct {
	NutanixUrl string
	UserName   string
	Password   string
}

func MakeConfig() Config {
	NutanixUrl := os.Getenv("NUTANIX_URL")
	UserName := os.Getenv("NUTANIX_USERNAME")
	Password := os.Getenv("NUTANIX_PASSWORD")
	if len(NutanixUrl) == 0 {
		panic("no value NUTANIX_URL")
	}
	if len(UserName) == 0 {
		panic("no value NUTANIX_USERNAME")
	}
	if len(Password) == 0 {
		panic("no value NUTANIX_PASSWORD")
	}
	return Config{NutanixUrl: NutanixUrl, UserName: UserName, Password: Password}
}
