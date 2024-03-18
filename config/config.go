package config

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/adrg/xdg"
)

type Config struct {
	LastSplitFile string
}

var default_config = Config{
	LastSplitFile: "",
}

const config_path = "speedruntimer/config"

func OpenConfigFile() (*Config, error) {
	cfgfilename, searcherr := xdg.SearchConfigFile(config_path)
	if searcherr != nil {
		// Assume the file does not exist; create default one
		var createerr error
		cfgfilename, createerr = createConfigFile()
		if createerr != nil {
			return nil, createerr
		}
	}

	cfgfile, openerr := os.Open(cfgfilename)
	if openerr != nil {
		log.Print("config open error")
		return nil, openerr
	}

	cfgfilebytes, cfgreaderr := io.ReadAll(cfgfile)
	if cfgreaderr != nil {
		log.Print("config read error")
		return nil, cfgreaderr
	}

	conf := &Config{}
	cfgunmarshalerr := json.Unmarshal(cfgfilebytes, conf)
	if cfgunmarshalerr != nil {
		log.Print("config unmarshal error")
		return nil, cfgunmarshalerr
	}

	return conf, nil
}

func createConfigFile() (string, error) {
	s, _ := xdg.ConfigFile(config_path) // Unhandled potential error
	newfile, filecreateerr := os.Create(s)
	if filecreateerr != nil {
		log.Print("config search error -> file create error")
		return "", filecreateerr
	}

	confbytes, marshalerr := json.Marshal(default_config)
	if marshalerr != nil {
		log.Print("config search error -> json marshal error")
		return "", marshalerr
	}

	newfile.Write(confbytes)

	return s, nil
}
