package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// ReadConfigs Reads yaml files from configuration directory with sub folders
// as application stage and merges config files in one configuration per stage
func ReadConfigs(cfgPath string) ([]byte, error) {
	if cfgPath == "" {
		cfgPath = "./configuration"
	}
	cfgPath = strings.TrimRight(cfgPath, "/")
	iSay("Config path: `%v`", cfgPath)

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return nil, err
	}

	stage := getStage()

	var (
		fileList        = map[string][]string{}
		stageDir        string
		configPathDepth int
	)

	err := filepath.Walk(cfgPath, func(path string, f os.FileInfo, err error) error {
		pathLen := len(strings.Split(path, "/"))

		if cfgPath == path {
			configPathDepth = pathLen + 1
			if strings.Contains(cfgPath, "./") {
				configPathDepth = pathLen
			}
			return nil
		}

		if pathLen > configPathDepth+1 {
			return filepath.SkipDir
		}

		if f.IsDir() && pathLen == configPathDepth {
			stageDir = f.Name()
			return nil
		}

		if filepath.Ext(f.Name()) == ".yaml" && (stageDir == "defaults" || stageDir == stage) {
			fileList[stageDir] = append(fileList[stageDir], f.Name())
		}

		return nil
	})
	if err != nil {
		iSay("Some error while walking through config dir!: `%+v`", err)
		return nil, err
	}

	// check defaults config existence. Fall down if not
	if _, ok := fileList["defaults"]; !ok || len(fileList["defaults"]) == 0 {
		log.Fatal("[config] defaults config is not found! Fall down.")
	}

	fileListResult := make(map[string][]string)
	configs := make(map[string]map[string]interface{})
	for folder, files := range fileList {
		for _, file := range files {
			configBytes, _ := ioutil.ReadFile(cfgPath + "/" + folder + "/" + file)

			var configFromFile map[string]map[string]interface{}

			if err = yaml.Unmarshal(configBytes, &configFromFile); err != nil {
				log.Fatalf("[config] %s %s config read fal! Fall down.", folder, file)
			}

			if _, ok := configFromFile[folder]; !ok {
				iSay("File %s excluded from %s (it is not for this stage)!", file, folder)
				continue
			}

			if _, ok := configs[folder]; !ok {
				configs[folder] = configFromFile[folder]
				continue
			}

			cc := configs[folder]
			if err := mergo.Merge(&cc, configFromFile[folder], mergo.WithOverride); err != nil {
				log.Fatalf("[config] merging files in folder error: %s", err)
			}
			configs[folder] = cc

			fileListResult[folder] = append(fileListResult[folder], file)
		}
	}

	iSay("Config files: `%+v`", fileListResult)

	config := configs["defaults"]
	c, ok := configs[stage]
	if ok {
		if err := mergo.Merge(&config, c, mergo.WithOverride); err != nil {
			log.Fatalf("[config] merging error: %s", err)
		}

		iSay("Stage `%s` config is loaded and merged with `defaults`", stage)
	}

	return yaml.Marshal(config)
}

// iSay Logs in stdout when quiet mode is off
func iSay(pattern string, args ...interface{}) {
	// if quietMode == false {
	log.Printf("[config] "+pattern, args...)
	// }
}

// getStage Load configuration for stage with fallback to 'development'
func getStage() (stage string) {
	stage = GetEnv("STAGE", "development")
	iSay("Current stage: `%s`", stage)
	return
}

// GetEnv Getting var from ENV with fallback param on empty
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
