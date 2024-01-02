package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var ErrConfigNotFound = errors.New("Config not found")
var ErrConfigValueNotFound = errors.New("Config key value not found")

type Configuration struct {
	ConfigFilePath string
	Generated      bool
	Data           map[string]map[string]string
}

type Config struct {
	ConfigName string
	Data       map[string]string
}

func (c *Configuration) SetFilePath(path string) {
	c.ConfigFilePath = path
}

func (c *Configuration) Init() error {
	c.Data = make(map[string]map[string]string)
	if _, err := os.Stat(c.ConfigFilePath); errors.Is(err, os.ErrNotExist) {
		LogInfo.Println("Config File not found, creating file")
		var dir = path.Dir(c.ConfigFilePath)
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
		}
		if f, err := os.Create(c.ConfigFilePath); err != nil {
			return err
		} else {
			f.Close()
		}
		LogResult.Printf("Config file created at %s\n", c.ConfigFilePath)
	}

	return nil
}

func (configuration *Configuration) SetConfig(name string, key string, value string) error {
	if configuration.Data[name] == nil {
		configuration.Data[name] = make(map[string]string)
	}

	configuration.Data[name][key] = value
	return nil
}

func (configuration *Configuration) LoadConfig() error {
	var file, err = os.OpenFile(configuration.ConfigFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if configuration.Data == nil {
		configuration.Data = make(map[string]map[string]string)
	}
	var scanner = bufio.NewScanner(file)
	var re = regexp.MustCompile(`^\[(?P<Name>[A-Za-z0-9_-]+)\]$`)
	var reEnv = regexp.MustCompile(`^(?P<Key>[A-Za-z0-9_-]+)=(?P<Value>.*)$`)
	var currentConfigName = ""
	for scanner.Scan() {
		var line = scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		var matches = re.FindStringSubmatch(line)
		if idx := re.SubexpIndex("Name"); idx != -1 && len(matches) > 0 {
			currentConfigName = matches[idx]
			if configuration.Data[currentConfigName] == nil {
				configuration.Data[currentConfigName] = make(map[string]string)
			}
			continue
		}
		if currentConfigName != "" {
			var envMatches = reEnv.FindStringSubmatch(line)
			if idxKey, idxValue := reEnv.SubexpIndex("Key"), reEnv.SubexpIndex("Value"); idxKey != -1 && idxValue != -1 {
				configuration.Data[currentConfigName][envMatches[idxKey]] = envMatches[idxValue]
			}
		}
	}
	return nil
}

func (configuration *Configuration) SaveConfig() error {
	var file, err = os.OpenFile(configuration.ConfigFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Truncate(0)
	file.Seek(0, 0)
	if configuration.Generated {
		file.WriteString("# THIS IS AUTOGENERATED CONFIGURATION FILE\n")
		file.WriteString("# DO NOT MAKE AN ATTEMPT TO EDIT THIS\n\n")
	}
	for configName, config := range configuration.Data {
		file.WriteString(fmt.Sprintf("[%s]\n", configName))
		for k, v := range config {
			file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
		file.WriteString("\n")
	}

	return nil
}

func (c *Configuration) GetConfig(name string) (*Config, error) {
	if c.Data[name] == nil {
		return nil, fmt.Errorf("config [%s] not found", name)
	}
	return &Config{
		ConfigName: name,
		Data:       c.Data[name],
	}, nil
}

func (c *Configuration) GetOrError(configName, key string) (string, error) {
	config, err := c.GetConfig(configName)
	if err != nil {
		return "", err
	}

	value, err := config.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *Configuration) GetPrtString() string {
	var titleSprint = color.New(color.FgYellow, color.Bold)
	var keySprint = color.New(color.FgCyan)
	var valueSprint = color.New(color.FgMagenta)

	var str = ""
	for name, config := range c.Data {
		str += titleSprint.Sprintf("[%s]\n", name)
		for k, v := range config {
			str += fmt.Sprintf("%s=%s\n", keySprint.Sprint(k), valueSprint.Sprint(v))
		}
		str += "\n"
	}
	return str
}

func (cnf *Config) Get(key string) (string, error) {
	if cnf.Data[key] == "" {
		return "", fmt.Errorf("value for \"%s\" not found in config [%s]", key, cnf.ConfigName)
	}
	return cnf.Data[key], nil
}

func LoadAllConfig(configs ...*Configuration) error {
	for _, c := range configs {
		if err := c.Init(); err != nil {
			return err
		}
		if err := c.LoadConfig(); err != nil {
			return err
		}
	}
	return nil
}
