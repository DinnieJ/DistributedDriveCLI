package helpers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
)

type Configuration struct {
	ConfigFilePath string
	Data           map[string]map[string]string
}

type Config struct {
	ConfigName string
	Data       map[string]string
}

func (c *Configuration) SetFilePath(path string) {
	c.ConfigFilePath = path
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
	configuration.Data = make(map[string]map[string]string)
	var scanner = bufio.NewScanner(file)
	var re = regexp.MustCompile(`^\[(?P<Name>[A-Za-z0-9_-]+)\]$`)
	var reEnv = regexp.MustCompile(`^(?P<Key>[A-Za-z0-9_-]+)=(?P<Value>.*)$`)
	var currentConfigName = ""
	for scanner.Scan() {
		var line = scanner.Text()
		if len(line) == 0 {
			continue
		}
		var matches = re.FindStringSubmatch(line)
		if idx := re.SubexpIndex("Name"); idx != -1 && len(matches) > 0 {
			currentConfigName = matches[idx]
			configuration.Data[currentConfigName] = make(map[string]string)
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

func (configuration *Configuration) WriteToConfigFile() error {
	var file, err = os.OpenFile(configuration.ConfigFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	file.Truncate(0)
	file.Seek(0, 0)
	defer file.Close()
	for configName, config := range configuration.Data {
		file.WriteString(fmt.Sprintf("[%s]\n", configName))
		for k, v := range config {
			file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
		file.WriteString("\n")
	}

	return nil
}

func (c *Configuration) GetConfig(name string) *Config {
	if c.Data[name] != nil {
		return nil
	}
	return &Config{
		ConfigName: name,
		Data:       c.Data[name],
	}
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
