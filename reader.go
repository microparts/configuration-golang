package config

import (
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

		if filepath.Ext(f.Name()) == ".yaml" {
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

	configs := make(map[string]interface{})
	for folder, files := range fileList {
		for _, file := range files {
			configBytes, _ := ioutil.ReadFile(cfgPath + "/" + folder + "/" + file)

			var configFromFile map[string]interface{}

			yaml.Unmarshal(configBytes, &configFromFile)

			var config interface{}
			if configs[folder] == nil {
				config = configFromFile[folder]
			} else {
				config = mergeMaps(configs[folder], configFromFile[folder])
			}
			configs[folder] = config
		}
	}

	var config = configs["defaults"]
	if c, ok := configs[stage]; ok {
		config = mergeMaps(config, c)
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

// mergeMaps Recursively merges interfaces
func mergeMaps(defaultMap interface{}, stageMap interface{}) interface{} {
	result := make(map[interface{}]interface{})
	switch defMap := defaultMap.(type) {
	case map[interface{}]interface{}:
		switch staMap := stageMap.(type) {
		case map[interface{}]interface{}:

			keys := getMapsKeys(defMap, staMap)
			for key := range keys {
				defVal, defOk := defMap[key]
				staVal, staOk := staMap[key]

				if defOk && !staOk {
					result[key] = defVal
				} else if !defOk && staOk {
					result[key] = staVal
				} else {
					result[key] = mergeMaps(defVal, staVal)
				}
			}
		}
	case interface{}:
		return stageMap
	}

	return result
}

// getMapsKeys Get config map keys slice
func getMapsKeys(defMap map[interface{}]interface{}, staMap map[interface{}]interface{}) map[interface{}]interface{} {
	keys := make(map[interface{}]interface{})
	for key := range defMap {
		keys[key] = true
	}
	for key := range staMap {
		keys[key] = true
	}
	return keys
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
