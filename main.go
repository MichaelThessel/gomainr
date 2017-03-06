package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/MichaelThessel/domainr/app"
	"github.com/MichaelThessel/domainr/cache"
	"github.com/MichaelThessel/domainr/file"
	"github.com/MichaelThessel/domainr/search"
	"github.com/MichaelThessel/domainr/search/source"

	"github.com/BurntSushi/toml"
)

type configPaths struct {
	configFile string
	baseDir    string
	dataDir    string
}

type config struct {
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

	configCreated, err := generateConfig()
	if err != nil {
		log.Panic(err)
	}

	// Exit with message that new config file has been created
	if configCreated {
		fmt.Println("New configuration created please update:", cp.configFile)
		os.Exit(1)
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
	baseDir := ".domainr"
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
	if c.NameCheap.Enabled {
		searchSource = source.Get(c.NameCheap, source.NameCheapSource)
	} else if c.GoDaddy.Enabled {
		searchSource = source.Get(c.GoDaddy, source.GoDaddySource)
	} else {
		fmt.Println("No search source enabled please update:", cp.configFile)
		os.Exit(1)
	}
	cache := cache.New(cp.dataDir)

	return search.New(searchSource, cache)
}

// generateConfig generates config files and directories
func generateConfig() (bool, error) {
	created := false
	// Create base directory
	err := file.CreateDirectory(cp.baseDir, 0700)
	if err != nil {
		return created, err
	}

	// Create data directory
	err = file.CreateDirectory(cp.dataDir, 0700)
	if err != nil {
		return created, err
	}

	// Create config file
	created, err = createConfigFile(cp.configFile)
	if err != nil {
		return created, err
	}

	return created, nil
}

// loadConfig loads the configuration from the config file
func loadConfig() error {
	configData, err := ioutil.ReadFile(cp.configFile)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(configData), &c)
	if err != nil {
		return err
	}

	return nil
}

// createConfigFile creates the default config file
func createConfigFile(fileName string) (bool, error) {
	fd, created, err := file.CreateFile(fileName)

	if err != nil {
		return created, err
	}

	// Exit if no new file was created
	if created == false {
		return created, nil
	}

	err = initConfig(fd)
	if err != nil {
		return created, err
	}

	return created, nil
}

// initConfig writes the default config to config file
func initConfig(fd *os.File) error {
	defaultConfig := `[namecheap]
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
