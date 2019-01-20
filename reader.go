package config

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var stage string

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

	getStage()

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

	iSay("Config files: `%+v`", fileList)

	// check defaults config existance. Fall down if not
	if _, ok := fileList["defaults"]; !ok || len(fileList["defaults"]) == 0 {
		log.Fatal("[config] defaults config is not found! Fall down.")
	}

	configs := make(map[string]map[string]interface{})
	for folder, files := range fileList {
		for _, file := range files {
			configBytes, _ := ioutil.ReadFile(cfgPath + "/" + folder + "/" + file)

			var configFromFile map[string]map[string]interface{}

			_ = yaml.Unmarshal(configBytes, &configFromFile)

			configs[folder] = configFromFile[folder]
		}
	}

	config := configs["defaults"]
	c, ok := configs[stage]
	if ok {
		if err := mergo.Merge(&config, c, mergo.WithOverride, mergo.WithAppendSlice); err != nil {
			log.Fatalf("config merging error: %s", err)
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
func getStage() {
	stage = GetEnv("STAGE", "development")
	iSay("Current stage: `%s`", stage)
}

// GetEnv Getting var from ENV with fallback param on empty
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
