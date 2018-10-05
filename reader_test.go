package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestReadConfigs(t *testing.T) {
	t.Run("Success parsing common dirs and files", func(t *testing.T) {
		t.Parallel()
		err := os.Setenv("STAGE", "test")
		configBytes, err := ReadConfigs("./test/configuration")
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		err = os.Unsetenv("STAGE")

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

	t.Run("Success parsing symlinked files and dirs", func(t *testing.T) {
		t.Parallel()
		err := os.Setenv("STAGE", "test")
		configBytes, err := ReadConfigs("./test/symnlinkedConfigs")
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		err = os.Unsetenv("STAGE")

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

	t.Run("Success parsing symlinked files and dirs in root", func(t *testing.T) {
		t.Parallel()
		err := os.Setenv("STAGE", "test")
		configBytes, err := ReadConfigs("/cfgs")
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		err = os.Unsetenv("STAGE")

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
		t.Parallel()
		_, err := ReadConfigs("")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}
