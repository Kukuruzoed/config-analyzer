package config

import "strings"

type Config map[string]any

func (config *Config) FindByKey(key string) (bool, []any) {
	var results = make([]any, 0)
	findKeyRecursive(config, key, &results)
	if len(results) == 0 {
		return false, nil
	}
	return true, results
}

func findKeyRecursive(config *Config, key string, results *[]any) {
	for k, v := range *config {
		if k == key || strings.Contains(k, key) {
			*results = append(*results, v)
		} else {
			nv, ok := v.(map[string]any)
			if ok {
				childConf := Config(nv)
				findKeyRecursive(&childConf, key, results)
			}
		}
	}
}
