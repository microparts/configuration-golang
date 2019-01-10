# golang-pkg

Yaml config reader (based on [spec rules](https://confluence.teamc.io/pages/viewpage.action?pageId=4227704))

## Usage example

```go
package main

import (
	"github.com/microparts/configuration-golang"
	"gopkg.in/yaml.v2"
	"log"
)


// Your app config structure. This must be related to yaml config file structure. Everything that is not
// in this struct will be passed and will not be processed.
// Keep in mind that inheritance must be implemented with `struct{}`
type ConfigStruct struct {
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

func main() {
	// Reader accept config path as param. Commonly it stored like STAGE in ENV.
	configPath := config.GetEnv("CONFIG_PATH", "./configuration")
	// Reading ALL config files in defaults configuration folder and recursively merge them with STAGE configs
	configBytes, err := config.ReadConfigs(configPath)
	if err != nil {
		log.Fatalf("config reading error: %+v", err)
	}

    var Config ConfigStruct 
    // unmarshal config into Config structure 
	err = yaml.Unmarshal(configBytes, &Config)
	if err != nil {
        log.Fatalf("config unmarshal error: %+v", err)
    }
} 
```