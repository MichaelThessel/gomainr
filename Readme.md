# Gomainr


Terminal cli app that checks the availability of domains for different configurations of keywords.

[![asciicast](https://asciinema.org/a/106322.png)](https://asciinema.org/a/106322?autoplay=1)

## Installation

You need to have [Go](https://golang.org/) installed.

```
# go get github.com/MichaelThessel/gomainr
# gomainr
```

Upon first execution gomainr will create a config file, print the config path and exit. You will need to edit the configuration file, add your API credentials and start gomainr again. The config file will be located at:

```
# $HOME/.gomainr/config
```

## API Keys

Currently gomainr supports both the [NameCheap.com](https://www.namecheap.com/support/api/intro.aspx) and [GoDaddy.com](https://developer.godaddy.com/) APIs. To use the app you need to obtain an API key from either service.

To be allowed to use the NameCheap API you need to fullfill certain [conditions](https://www.namecheap.com/support/knowledgebase/article.aspx/9739/63/api--faq#c). It will also take up to 48 hours for NameCheap to activate your API access (if you ask nicely in the live chat they might do it right away though :). There are no restrictions for access to the GoDaddy API. Unless you already have a bunch of domains with NameCheap it's probably easiest to get a GoDaddy key.

## Usage

The main purpose of this tool is to find available domains for different keywords. I.e.:

* Keywords 1: foo bar
* Keywords 2: alice bob
* Tlds: com net

Will search for:

* fooalice.com
* fooalice.net
* baralice.com
* baralice.net
* foobob.com
* foobob.net
* barbob.com
* barbob.net

and return the available domain names.

Keywords 2 is optional, so you can just search for various domains among differnt TLDs.

You can save a session to a file and load it later again. This way you can view the results again without performing a new search. In addition this allows you to modify the keywords and repeat a search without typing the keywords all over again. 

## Keyboard Shortcuts

Shortcut | Action
---------|-------
<kbd>CTRL</kbd>+<kbd>q</kbd> | Quit
<kbd>CTRL</kbd>+<kbd>/</kbd> | Search
<kbd>UP</kbd>, <kbd>DOWN</kbd>, <kbd>TAB</kbd> | Navigate
<kbd>CTRL</kbd>+<kbd>j</kbd> | Scroll result list down
<kbd>CTRL</kbd>+<kbd>k</kbd> | Scroll result list up
<kbd>CTRL</kbd>+<kbd>s</kbd> | Save session
<kbd>CTRL</kbd>+<kbd>l</kbd> | Load session

## Notes

To speed up consecutive searches and to keep things light on the APIs gomainr caches API request results for 24hrs. If you want to flush the cache for some reason you can delete the contents of this directory:

```
# $HOME/.gomainr/data
```

## Thanks

This project utilizes the following 3rd party packages

* [GOCUI](https://github.com/jroimartin/gocui)
* [TOML](https://github.com/BurntSushi/toml)
* [diskv](https://github.com/peterbourgon/diskv)
* [go-namecheap](https://github.com/billputer/go-namecheap)
