package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"sort"
	"strings"

	"github.com/MichaelThessel/gomainr/file"
	"github.com/MichaelThessel/gomainr/search"

	"github.com/jroimartin/gocui"
)

type App struct {
	gui         *gocui.Gui
	currentView int
	s           *search.Search
	state       *state
}

type state struct {
	Parts1   []string
	Parts2   []string
	Tlds     []string
	Domains  []string
	Settings map[string]bool
}

func New(s *search.Search) *App {
	a := new(App)

	a.s = s

	var err error
	a.gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	// Default settings
	a.state = new(state)
	a.state.Settings = map[string]bool{
		"TLDSubstitutions": false,
	}

	a.initGui()

	return a
}

// Loop starts the GUI loop
func (a *App) Loop() {
	if err := a.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// Close closes the app
func (a *App) Close() {
	a.gui.Close()
}

// initGui initializes the GUI
func (a *App) initGui() {
	le = lineEditor{gocui.DefaultEditor}

	// Defaults
	a.gui.Cursor = true
	a.gui.InputEsc = false
	a.gui.BgColor = gocui.ColorDefault
	a.gui.FgColor = gocui.ColorDefault

	// Set Layout function
	a.gui.SetManagerFunc(a.Layout)

	a.currentView = -1

	// Set Keybindings
	a.setKeyBindings()
}

// writeConsole writes string to console
func (a *App) writeConsole(text string, isError bool) {
	if isError {
		text = decorate(text, "red")
	}
	a.writeView(viewConsole, text)
}

// search builds the domain names from the parts and updates the result list with
// the available ones
func (a *App) search(g *gocui.Gui, v *gocui.View) error {
	a.clearView(viewDomain)

	if !a.validate() {
		return nil
	}

	a.writeConsole("Searching ...", false)

	// Generate domain list from parts
	domains := a.s.BuildQuery(
		a.state.Parts1,
		a.state.Parts2,
		a.state.Tlds,
		a.state.Settings["TLDSubstitutions"],
	)

	if len(domains) == 0 {
		a.writeConsole("No possible searches!", true)
	}

	jobs := make(chan string, len(domains))
	found := make(chan string)
	complete := make(chan bool)

	for _, domain := range domains {
		jobs <- domain
	}
	close(jobs)

	// Create the workers that fetch available domain names
	workerCount := len(domains)
	var apiErr error
	for i := 0; i < workerCount; i++ {
		go func(jobs <-chan string, found chan<- string, complete chan<- bool) {
			for domain := range jobs {
				available, err := a.s.IsAvailable(domain)
				if err != nil {
					apiErr = err
					break
				}
				if available {
					found <- domain
				}
			}
			complete <- true
		}(jobs, found, complete)
	}

	// Update the domain list as results come in
	foundDomains := []string{}
	go func(found <-chan string) {
		for domain := range found {
			foundDomains = append(foundDomains, domain)
			a.gui.Execute(func(g *gocui.Gui) error {
				sort.Strings(foundDomains)

				a.writeView(viewDomain, decorate(strings.Join(foundDomains, "\n"), "blue"))
				return nil
			})
		}

		a.gui.Execute(func(g *gocui.Gui) error {
			if apiErr != nil {
				a.writeConsole(
					fmt.Sprintf("API error: %s", apiErr),
					true,
				)
			} else {
				a.writeConsole(
					fmt.Sprintf(
						"Search complete: Scanned %d domain(s) - %d domain(s) available",
						len(domains),
						len(foundDomains),
					),
					false,
				)
			}
			return nil
		})
	}(found)

	// Signal completion to the domain list update goroutine
	go func(found chan string, complete chan bool) {
		i := 0
		for range complete {
			i++
			if i == workerCount {
				close(found)
				close(complete)
			}

		}
	}(found, complete)

	return nil
}

// saveModal opens the save modal
func (a *App) saveModal(g *gocui.Gui, v *gocui.View) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	a.showModal(
		viewSave,
		usr.HomeDir+string(os.PathSeparator),
		0.4,
		0.07,
	)

	a.writeConsole("Please enter the file name you want to save the results to.", false)

	return nil
}

// loadModal opens the load modal
func (a *App) loadModal(g *gocui.Gui, v *gocui.View) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	a.showModal(
		viewLoad,
		usr.HomeDir+string(os.PathSeparator),
		0.4,
		0.07,
	)

	a.writeConsole("Please enter the file name you want to load data from.", false)

	return nil
}

// close Closes a view
func (a *App) closeModal(g *gocui.Gui, v *gocui.View) error {
	a.closeView(v.Name())
	return nil
}

// save saves the current state to a file
func (a *App) save(g *gocui.Gui, v *gocui.View) error {
	saveFile := strings.TrimSpace(v.Buffer())

	if saveFile == "" {
		return nil
	}

	fd, _, err := file.CreateFile(saveFile)
	if err != nil {
		a.writeConsole(fmt.Sprintf("Couldn't create file: %s", saveFile), false)
		return nil
	}

	json, err := json.MarshalIndent(a.state, "", "    ")
	if err != nil {
		a.writeConsole("Error generating output", false)
		return nil
	}

	if _, err := fd.Write(json); err != nil {
		a.writeConsole(fmt.Sprintf("Couldn't write to file: %s", saveFile), false)
		return nil
	}

	a.writeConsole(fmt.Sprintf("The results have been saved to: %s", saveFile), false)

	a.closeView(v.Name())

	return nil
}

// load loads state from a file
func (a *App) load(g *gocui.Gui, v *gocui.View) error {
	loadFile := strings.TrimSpace(v.Buffer())
	data, err := file.ReadFile(loadFile)
	if err != nil {
		a.writeConsole(fmt.Sprintf("Couldn't read: %s", loadFile), false)
		return nil
	}

	if err := json.Unmarshal(data, a.state); err != nil {
		a.writeConsole(fmt.Sprintf("Couldn't parse file: %s", loadFile), false)
		return nil
	}

	a.writeView(viewPart1, strings.Join(a.state.Parts1, " "))
	a.writeView(viewPart2, strings.Join(a.state.Parts2, " "))
	a.writeView(viewTLD, strings.Join(a.state.Tlds, " "))
	a.writeView(viewDomain, decorate(strings.Join(a.state.Domains, "\n"), "blue"))

	a.writeConsole(fmt.Sprintf("The results have been loaded from: %s", loadFile), false)

	a.closeView(v.Name())

	return nil
}

// toggleTLDSubstitutions toggles replacement of base domain ending with TLD
func (a *App) toggleTLDSubsitutions(g *gocui.Gui, v *gocui.View) error {
	a.setSetting("TLDSubstitutions", !a.state.Settings["TLDSubstitutions"])

	return nil
}

// setSetting updates a setting
func (a *App) setSetting(name string, value bool) {
	a.state.Settings[name] = value
}

// validate validates that the required fields are populated
func (a *App) validate() bool {
	if len(a.state.Parts1) == 0 {
		a.writeConsole("\"Parts 1\" cannot be empty! Please enter space seperated list of domain parts.", true)
		return false
	}

	// If TLD substitutions are disabled TLD needs to be set
	if len(a.state.Tlds) == 0 && !a.state.Settings["TLDSubstitutions"] {
		a.writeConsole("\"TLDs\" cannot be empty! Please enter space seperated list of TLDs to scan.", true)
		return false
	}

	// Validate TLDs
	tlds := a.parseLine(viewTLD)
	if err := search.ValidateTlds(tlds); err != nil {
		a.writeConsole(fmt.Sprintf("%s", err), true)
		return false
	}

	return true
}

// updateState saves the current state of views
func (a *App) updateState() {
	a.state.Parts1 = a.parseLine(viewPart1)
	a.state.Parts2 = a.parseLine(viewPart2)
	a.state.Tlds = a.parseLine(viewTLD)
	a.state.Domains = a.parseLine(viewDomain)
}

// updateViews updates the views based on the current state
func (a *App) updateViews() {
	if a.state.Settings["TLDSubstitutions"] {
		a.writeView(viewSettings, "[X] TLD substitutions")
	} else {
		a.writeView(viewSettings, "[ ] TLD substitutions")
	}
}

// getViewWords returns the list of words in a view (space separated)
func (a *App) parseLine(view string) []string {
	v, _ := a.gui.View(view)
	parts := strings.Fields(v.Buffer())

	m := make(map[string]bool)
	unique := make([]string, 0, len(parts))
	for _, part := range parts {
		if ok := m[part]; !ok {
			unique = append(unique, part)
			m[part] = true
		}
	}

	return unique
}

// decorate changes the color of a string
func decorate(s string, color string) string {
	switch color {
	case "blue":
		s = "\x1b[0;32m" + s
	case "red":
		s = "\x1b[0;31m" + s
	default:
		return s
	}

	return s + "\x1b[0m"
}
