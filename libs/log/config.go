package logger

import (
	"encoding/json"
	"fmt"
	"os"
)

type LogConfig struct {
	configFile     string
	LogSplitPolicy string
	LogRootDir     string
	LogLevel       string
	NamePrefix     string
	SplitError     bool
	DetailInfo     bool
}

func NewLogConfig(path string) *LogConfig {
	return &LogConfig{
		configFile:     path,
		LogSplitPolicy: "perhour",
		LogRootDir:     "./log",
		LogLevel:       "info",
		NamePrefix:     "smartgo",
		SplitError:     true,
		DetailInfo:     false,
	}
}

func (self *LogConfig) LoadConfig() error {
	file, err := os.Open(self.configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&self)
	if err != nil {
		return err
	}

	return nil
}

func (self LogConfig) String() string {
	return fmt.Sprintf("{ LogSplitPolicy:[%s], LogRootDir:[%s], LogLevel:[%s], NamePrefix:[%s], DetailInfo:[%v], SplitError:[%v] }",
		self.LogSplitPolicy,
		self.LogRootDir,
		self.LogLevel,
		self.NamePrefix,
		self.DetailInfo,
		self.SplitError)
}
