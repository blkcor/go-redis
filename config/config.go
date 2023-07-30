package config

import (
	"bufio"
	"errors"
	"github.com/blkcor/go-redis/lib/logger"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ServerProperties struct {
	Bind           string `cfg:"bind"`
	Port           int    `cfg:"port"`
	AppendOnly     bool   `cfg:"appendOnly"`
	AppendFilename string `cfg:"appendFilename"`
	MaxClients     int    `cfg:"maxClients"`
	RequirePass    string `cfg:"requirePass"`
	Databases      int    `cfg:"databases"`

	Peers []string `cfg:"peers"`
	Self  string   `cfg:"self"`
}

// Properties holds global config properties
var (
	Properties *ServerProperties
	// ErrInvalidConfig is an error for invalid configuration format.
	ErrInvalidConfig = errors.New("invalid configuration format")
)

func init() {
	//default config
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       6379,
		AppendOnly: false,
	}
}

func SetupConfig(configFile string) {
	Properties = &ServerProperties{}
	err := LoadConfig(configFile, Properties)
	if err != nil {
		logger.Error("Error loading config file: %s", err.Error())
		panic(err)
	}
}

// LoadConfig loads the configuration from the given configFile into the config interface.
func LoadConfig(configFile string, config interface{}) error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidConfig
	}

	if rv.Elem().Kind() != reflect.Struct {
		return ErrInvalidConfig
	}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue // Skip empty lines and comments
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue // Invalid line format, skip
		}

		fieldName := strings.TrimSpace(parts[0])
		fieldValue := strings.TrimSpace(parts[1])

		field := rv.Elem().FieldByName(fieldName)
		if !field.IsValid() {
			continue // Unknown field, skip
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(fieldValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(fieldValue)
			if err != nil {
				continue // Invalid int value, skip
			}
			field.SetInt(int64(intValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(fieldValue)
			if err != nil {
				continue // Invalid bool value, skip
			}
			field.SetBool(boolValue)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
