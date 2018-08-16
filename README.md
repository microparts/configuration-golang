# golang-pkg

Читатель yaml конфигураций по спецификации https://confluence.teamc.io/pages/viewpage.action?pageId=4227704

## Пример использования

```go
package config

import (
	"gitlab.teamc.io/teamc.io/microservice/configuration/golang-pkg"
	"gopkg.in/yaml.v2"
)

var Config *confStruct

// confStruct file structure
type confStruct struct {
	Log struct {
		Level  string `yaml:"level"`
		Debug  bool   `yaml:"debug"`
		Format string `yaml:"yaml"`
	} `yaml:"log"`
	Database struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		User           string `yaml:"user"`
		Pass           string `yaml:"password"`
		Name           string `yaml:"name"`
		SslMode        string `yaml:"sslMode"`
		Logs           bool   `yaml:"logs"`
		MigrateOnStart bool   `yaml:"migrateOnStart"`
		MigrationPath  string `yaml:"migrationsPath"`
	} `yaml:"db"`
	Http struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"ws"`
}

func InitConfig() error {
	configBytes, err := config.ReadConfigs()
	if err != nil {
		return err
	}

	return yaml.Unmarshal(configBytes, &Config)
} 
```