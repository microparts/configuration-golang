package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type configMap map[interface{}]interface{}

var (
	stage     = "development"
	cfgPath   = flag.String("conf", "./configuration", "Path to directory with configuration")
	quietMode = flag.Bool("quiet", false, "Quiet mode")
)

// parseFlags Parse CLI flags.
func parseFlags() {
	flag.Parse()
	iSay("Config path: `%s`", cfgPath)
}

// ReadConfigs Reads yaml files from configuration directory with sub folders
// as application stage and merges config files in one configuration per stage
func ReadConfigs() ([]byte, error) {
	parseFlags()
	getStage()

	var (
		fileList = map[string][]string{}
		stageDir string
	)

	err := filepath.Walk(*cfgPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			stageDir = f.Name()
			return nil
		}

		if filepath.Ext(f.Name()) == ".yaml" {
			fileList[stageDir] = append(fileList[stageDir], f.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	iSay("Config files: `%+v`", fileList)

	configs := make(map[string]interface{})
	for folder, files := range fileList {
		for _, file := range files {
			configBytes, err := ioutil.ReadFile(*cfgPath + "/" + folder + "/" + file)

			var configFromFile map[string]interface{}

			err = yaml.Unmarshal(configBytes, &configFromFile)
			if err != nil {
				return nil, err
			}

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
func iSay(pattern string, args ... interface{}) {
	if *quietMode == false {
		log.Printf("[config] "+pattern, args)
	}
}

// mergeMaps Recursively merges interfaces
func mergeMaps(defaultMap interface{}, stageMap interface{}) (interface{}) {
	result := make(configMap)
	switch defMap := defaultMap.(type) {
	case configMap:
		switch staMap := stageMap.(type) {
		case configMap:

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
		switch staMap := stageMap.(type) {
		case configMap:

			if defMap == nil {
				result = staMap
			}
		case interface{}:
			return staMap
		}
	default:
		log.Fatalf("[config] Config merging error: structures type are not match each other")
	}

	return result
}

// getMapsKeys Get config map keys slice
func getMapsKeys(defMap configMap, staMap configMap) (configMap) {
	keys := make(configMap)
	for key := range defMap {
		keys[key] = true
	}
	for key := range staMap {
		keys[key] = true
	}
	return keys
}

// getStage Load configuration for stage with fallback to 'development'
func getStage() string {
	stage := getEnv("STAGE", "development")
	iSay("Current stage: `%s`", stage)
	return stage
}

// getEnv Getting var from ENV with fallback param on empty
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
