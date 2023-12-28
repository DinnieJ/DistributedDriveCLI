package helpers

import (
	"bufio"
	"os"
	"regexp"
)

type Config struct {
	ConfigName string
	Data       map[string]string
}

func LoadConfig(path string) (map[string]map[string]string, error) {
	var file, err = os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var configuration = make(map[string]map[string]string)
	var scanner = bufio.NewScanner(file)
	var re = regexp.MustCompile(`^\[(?P<Name>[A-Za-z0-9_-]+)\]$`)
	var reEnv = regexp.MustCompile(`^(?P<Key>[A-Za-z0-9_-]+)=(?P<Value>.*)$`)
	var currentConfigName = ""
	for scanner.Scan() {
		var line = scanner.Text()
		var matches = re.FindStringSubmatch(line)
		if idx := re.SubexpIndex("Name"); idx != -1 && len(matches) > 0 {
			currentConfigName = matches[idx]
			configuration[currentConfigName] = make(map[string]string)
			continue
		}

		if currentConfigName != "" {
			var envMatches = reEnv.FindStringSubmatch(line)
			if idxKey, idxValue := reEnv.SubexpIndex("Key"), reEnv.SubexpIndex("Value"); idxKey != -1 && idxValue != -1 {
				configuration[currentConfigName][envMatches[idxKey]] = envMatches[idxValue]
			}
		}
	}

	return configuration, nil
}
