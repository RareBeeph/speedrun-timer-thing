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

func OpenConfigFile() *Config {
	cfgfilename, searcherr := xdg.SearchConfigFile(config_path)
	if searcherr != nil {
		// Assume the file does not exist; create default one
		s, _ := xdg.ConfigFile(config_path)
		newfile, _ := os.Create(s)                   // Unhandled potential error
		confbytes, _ := json.Marshal(default_config) // Unhandled potential error
		newfile.Write(confbytes)
	}

	// TODO: handle these errors
	cfgfile, openerr := os.Open(cfgfilename)
	if openerr != nil {
		log.Print("config open error")
		// etc
	}

	cfgfilebytes, cfgreaderr := io.ReadAll(cfgfile)
	if cfgreaderr != nil {
		log.Print("config read error")
		// etc
	}

	conf := &Config{}
	cfgunmarshalerr := json.Unmarshal(cfgfilebytes, conf)
	if cfgunmarshalerr != nil {
		log.Print("config unmarshal error")
		// etc
	}

	return conf
}
