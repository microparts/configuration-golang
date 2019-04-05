package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestReadConfigs(t *testing.T) {
	t.Run("Success parsing common dirs and files", func(t *testing.T) {
		configBytes, err := ReadConfigs("./test/configuration")
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
	t.Run("Success parsing complex dirs and files", func(t *testing.T) {
		configBytes, err := ReadConfigs("./test/configuration2")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type hbParams struct {
			AreaMapping map[string]string `yaml:"area_mapping"`
			Url         string            `yaml:"url"`
			Username    string            `yaml:"username"`
			Password    string            `yaml:"password"`
		}

		type cfg struct {
			HotelbookParams hbParams `yaml:"hotelbook_params"`
			Logging         string   `yaml:"logging"`
			DefaultList     []string `yaml:"default_list"`
			Databases       struct {
				Redis struct {
					Master struct {
						Username string `yaml:"username"`
						Password string `yaml:"password"`
					} `yaml:"master"`
				} `yaml:"redis"`
			} `yaml:"databases"`
		}

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			HotelbookParams: hbParams{
				AreaMapping: map[string]string{"KRK": "Krakow", "MSK": "Moscow", "CHB": "Челябинск"},
				Url:         "https://hotelbook.com/xml_endpoint",
				Username:    "TESt_USERNAME",
				Password:    "PASSWORD",
			},
			DefaultList: []string{"bar", "baz"},
			Logging:     "info",
			Databases: struct {
				Redis struct {
					Master struct {
						Username string `yaml:"username"`
						Password string `yaml:"password"`
					} `yaml:"master"`
				} `yaml:"redis"`
			}{Redis: struct {
				Master struct {
					Username string `yaml:"username"`
					Password string `yaml:"password"`
				} `yaml:"master"`
			}{Master: struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{Username: "R_USER", Password: "R_PASS"}}},
		}

		assert.EqualValues(t, refConfig, config)
	})

	t.Run("Success parsing symlinked files and dirs", func(t *testing.T) {
		configBytes, err := ReadConfigs("./test/symnlinkedConfigs")
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

	if GetEnv("IN_CONTAINER", "") == "true" {
		t.Run("Success parsing symlinked files and dirs in root", func(t *testing.T) {
			configBytes, err := ReadConfigs("/cfgs")
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
	}

	t.Run("Fail dir not found", func(t *testing.T) {

		_, err := ReadConfigs("")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}
