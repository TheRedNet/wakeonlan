package config

import (
	"encoding/json"
	"io/ioutil"
)

type Device struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
}

type Config struct {
	Devices       []Device `json:"devices"`
	IsAuthEnabled bool     `json:"isAuthEnabled"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
}

var configPath string

func GetConfigPath() string {
	if configPath == "" {
		configPath = "config.json"
	}
	return configPath
}

func SetConfigPath(path string) {
	configPath = path
}

func LoadDevices() (map[string]string, error) {
	config, err := LoadConfig()
	var devices []Device
	devices = config.Devices
	deviceMap := make(map[string]string)
	for _, d := range devices {
		deviceMap[d.Name] = d.Mac
	}
	return deviceMap, err
}

func SaveDevices(deviceMap map[string]string) error {
	// convert deviceMap to devices slice
	var devices []Device
	for name, mac := range deviceMap {
		devices = append(devices, Device{Name: name, Mac: mac})
	}
	var config Config
	config, err := LoadConfig()
	config.Devices = devices
	err = SaveConfig(config)
	return err
}

func LoadConfig() (Config, error) {
	data, err := ioutil.ReadFile(GetConfigPath())
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
func SaveConfig(config Config) error {
	// marshal config to JSON data
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	// write data to file
	return ioutil.WriteFile(GetConfigPath(), data, 0644)
}
