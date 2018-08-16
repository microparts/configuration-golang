package config

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestReadConfigs(t *testing.T) {
	t.Parallel()
	t.Run("Success parsing", func(t *testing.T) {
		os.Setenv("STAGE", "test")
		configBytes, err := ReadConfigs("./test/good")
		os.Unsetenv("STAGE")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type cfg struct {
			Debug bool `yaml:"debug"`
			Log   struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			} `yaml:"log"`
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		}

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			Debug: true,
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "warn", Format: "json"},
			Host: "localhost",
			Port: "8080",
		}

		assert.EqualValues(t, refConfig, config)
	})

	t.Run("Fail dir not found", func(t *testing.T) {
		_, err := ReadConfigs("")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}
