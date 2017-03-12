package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/MichaelThessel/gomainr/app"
	"github.com/MichaelThessel/gomainr/cache"
	"github.com/MichaelThessel/gomainr/file"
	"github.com/MichaelThessel/gomainr/search"
	"github.com/MichaelThessel/gomainr/search/source"

	"github.com/BurntSushi/toml"
)

type configPaths struct {
	configFile string
	baseDir    string
	dataDir    string
}

type config struct {
	DNS       *source.DNSConfig
	NameCheap *source.NameCheapConfig
	GoDaddy   *source.GoDaddyConfig
}

var c *config
var a *app.App
var cp *configPaths

func init() {
	if err := initPaths(); err != nil {
		log.Panic(err)
	}

	if err := generateConfig(); err != nil {
		log.Panic(err)
	}

	if err := loadConfig(); err != nil {
		log.Panic(err)
	}
}

func main() {
	s := initSearch()
	a = app.New(s)
	defer a.Close()

	// Main loop
	a.Loop()
}

// initPaths sets config and data storage paths
func initPaths() error {
	configFile := "config"
	baseDir := ".gomainr"
	dataDir := "data"

	usr, err := user.Current()
	if err != nil {
		return err
	}

	cp = new(configPaths)

	cp.baseDir = usr.HomeDir + string(os.PathSeparator) + baseDir
	cp.dataDir = cp.baseDir + string(os.PathSeparator) + dataDir
	cp.configFile = cp.baseDir + string(os.PathSeparator) + configFile

	return nil
}

// initSearch initializes the searcher
func initSearch() *search.Search {
	var searchSource source.Source
	if c.DNS != nil && c.DNS.Enabled {
		searchSource = source.Get(c.DNS, source.DNSSource)
	} else if c.NameCheap != nil && c.NameCheap.Enabled {
		searchSource = source.Get(c.NameCheap, source.NameCheapSource)
	} else if c.GoDaddy != nil && c.GoDaddy.Enabled {
		searchSource = source.Get(c.GoDaddy, source.GoDaddySource)
	} else {
		fmt.Println("No search source enabled please update:", cp.configFile)
		os.Exit(1)
	}
	cache := cache.New(cp.dataDir)

	return search.New(searchSource, cache)
}

// generateConfig generates config files and directories
func generateConfig() error {
	// Create base directory
	if err := file.CreateDirectory(cp.baseDir, 0700); err != nil {
		return err
	}

	// Create data directory
	if err := file.CreateDirectory(cp.dataDir, 0700); err != nil {
		return err
	}

	// Create config file
	if err := createConfigFile(cp.configFile); err != nil {
		return err
	}

	return nil
}

// loadConfig loads the configuration from the config file
func loadConfig() error {
	configData, err := ioutil.ReadFile(cp.configFile)
	if err != nil {
		return err
	}

	if _, err := toml.Decode(string(configData), &c); err != nil {
		return err
	}

	return nil
}

// createConfigFile creates the default config file
func createConfigFile(fileName string) error {
	fd, created, err := file.CreateFile(fileName)

	if err != nil {
		return err
	}

	// Exit if no new file was created
	if created == false {
		return nil
	}

	if err := initConfig(fd); err != nil {
		return err
	}

	return nil
}

// initConfig writes the default config to config file
func initConfig(fd *os.File) error {
	defaultConfig := `[dns]
Enabled = true
[namecheap]
APIUser = ""
APIToken = ""
UserName = ""
Enabled = false
[godaddy]
Key = ""
Secret = ""
Enabled = false
`
	_, err := fd.WriteString(defaultConfig)
	return err
}
