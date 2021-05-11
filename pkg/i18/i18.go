package i18

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type i18Map struct {
	sync.Mutex
	i18Map map[string]map[string]string
}

var i18 i18Map

// ParseJSON parse json i18 file
func ParseJSON(lang string) error {
	i18.Lock()
	defer i18.Unlock()
	jsonFile, err := os.Open("./i18/" + lang + ".json")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &i18.i18Map)
	if err != nil {
		return err
	}

	return nil
}

// Get translate value
func Get(group, key string) string {
	i18.Lock()
	defer i18.Unlock()
	if val, ok := i18.i18Map[group][key]; ok {
		return val
	}
	return key
}
