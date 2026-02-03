package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

var loadDotEnvOnce sync.Once
var loadDotEnvErr error

func loadDotEnvIfPresent(path string) error {
	loadDotEnvOnce.Do(func() {
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				loadDotEnvErr = nil
				return
			}
			loadDotEnvErr = err
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "#") {
				continue
			}
			if strings.HasPrefix(line, "export ") {
				line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
			}

			k, v, ok := strings.Cut(line, "=")
			if !ok {
				loadDotEnvErr = fmt.Errorf("invalid .env line: %q", line)
				return
			}
			key := strings.TrimSpace(k)
			val := strings.TrimSpace(v)
			if key == "" {
				loadDotEnvErr = fmt.Errorf("invalid .env line: %q", line)
				return
			}

			if len(val) >= 2 {
				if (val[0] == '\'' && val[len(val)-1] == '\'') || (val[0] == '"' && val[len(val)-1] == '"') {
					val = val[1 : len(val)-1]
				}
			}

			if _, exists := os.LookupEnv(key); exists {
				continue
			}
			if err := os.Setenv(key, val); err != nil {
				loadDotEnvErr = err
				return
			}
		}
		if err := scanner.Err(); err != nil {
			loadDotEnvErr = err
			return
		}
	})

	return loadDotEnvErr
}
