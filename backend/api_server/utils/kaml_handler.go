package utils

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"

	"api_server/logger"
)

type Config struct {
	FilePath string
	YamlData map[interface{}]interface{}
}

type ConfigInterface interface {
	FromYaml(p string) map[interface{}]interface{}
	ToYaml(p string)
	FromKaml(p string) map[interface{}]interface{}
	ToKaml(p string)
	SetValue(h string, k string, v interface{})
	GetValue(h string, k string) interface{}
}

func New() *Config {
	return &Config{}
}

func (c *Config) FromYaml(path string) map[interface{}]interface{} {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Error("File is Not Exist - ", path)
		return nil
	}

	r, err := os.ReadFile(path)
	if err != nil {
		logger.Error("File Open Failed - ", path, " : ", err)
		return nil
	}

	err = yaml.Unmarshal(r, &c.YamlData)
	if err != nil {
		logger.Error("Data Unmarshal Failed - ", err)
	}

	return c.YamlData
}

func (c *Config) FromKaml(path string) map[interface{}]interface{} {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Error("File is Not Exist - ", path)
		return nil
	}

	s := CreateSecurity()
	dec := s.DecryptFile(path)

	err := yaml.Unmarshal([]byte(dec), &c.YamlData)
	if err != nil {
		logger.Error("Data Unmarshal Failed - ", err)
		return nil
	}

	return c.YamlData
}

func (c *Config) ToYaml(path string) {
	if c.YamlData == nil {
		logger.Error("No Yaml Data")
		return
	}

	yamlBytes, err := yaml.Marshal(c.YamlData)
	if err != nil {
		logger.Error("Data Marshaling Failed - ", err)
		return
	}

	err = os.WriteFile(path, yamlBytes, os.FileMode(0666))
	if err != nil {
		logger.Error("Data Writing Failed - ", path, " : ", err)
		return
	}
}

func (c *Config) ToKaml(path string) {
	if c.YamlData == nil {
		logger.Error("No Kaml Data")
		return
	}

	body := yamlToJson(c.YamlData)
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		logger.Error("Data Marshaling Failed - ", err)
		return
	}

	s := CreateSecurity()
	if err := s.EncryptFile(string(jsonBytes), path); err != nil {
		logger.Debug(err)
	}
}

func (c *Config) SetValue(h string, k string, v interface{}) {
	if c.YamlData == nil {
		logger.Error("No Yaml Data, Yaml Data is nil")
		return
	}

	hd := c.YamlData[h].(map[string]interface{})
	hd[k] = v
}

func (c *Config) GetValue(h string, k string) interface{} {

	if c.YamlData == nil {
		logger.Error("No Yaml Data")
		return nil
	}

	hd := c.YamlData[h].(map[string]interface{})

	return hd[k]
}

func yamlToJson(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = yamlToJson(v)
		}
		return m2
	case map[string]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k] = yamlToJson(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = yamlToJson(v)
		}
	}
	return i
}

func CreateTrainConfigurationFile() string { // string is File Path
	// DataPath
	// BackBone
	// Width
	// Height
	// Train Batch Size
	// Valid Batch Size
	// Print Freq
	// Val Freq
	// Optimizer
	// Base LR
	// Scheduler
	// Epochs
	// Weight
	// Save Top K
	// Class List []
	// Default Config File
	return ""
}
