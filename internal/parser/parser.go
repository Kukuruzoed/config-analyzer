package parser

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/Kukuruzoed/config-analyzer/internal/config"
	"gopkg.in/yaml.v3"
)

type Config = config.Config

func ParseFile(path string) (Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parse(data)
}

func ParseReader(r io.Reader) (Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return parse(data)
}

func parse(data []byte) (Config, error) {
	var result Config

	err := json.Unmarshal(data, &result)

	if err == nil {
		return result, err
	}

	err = yaml.Unmarshal(data, &result)
	if err == nil {
		return result, err
	}

	return result, errors.New("Некорректный формат конфигурационного файла")
}
