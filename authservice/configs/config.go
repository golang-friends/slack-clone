package configs

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Configuration struct - config for the Auth service
type Configuration struct {
	ListenInterface string   `json:"listen_interface"`
	DB              Database `json:"database"`
	JWTSecret       string   `json:"jwt_secret"`
	Perf            int      `json:"perf"`
}

// Database struct - config for the database
type Database struct {
	DBURI  string `json:"dburi"`
	DBName string `json:"dbname"`
	DBUser string `json:"dbuser"`
	DBPass string `json:"dbpassword"`
}

// ConfigFromFile parses the given file and returns the config
func ConfigFromFile(filename string) (Configuration, error) {
	var conf Configuration

	confjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return conf, errors.Wrapf(err, "Failed to open the config file at: %s", filename)
	}

	if err := json.Unmarshal(confjson, &conf); err != nil {
		return conf, errors.Wrapf(err, "Unable to parse config file at: %s", filename)
	}

	return conf, nil
}
