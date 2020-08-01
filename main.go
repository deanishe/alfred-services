// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-08-01

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/pkg/errors"
	"howett.net/plist"
)

const (
	servicesList = "${HOME}/Library/Caches/com.apple.nsservicescache.plist"
	inputKey     = "input.json"
)

var (
	iconError   = &aw.Icon{Value: "icons/error.png"}
	iconWarning = &aw.Icon{Value: "icons/warning.png"}
)

var (
	wf *aw.Workflow

	fs       = flag.NewFlagSet("alfred-services", flag.ExitOnError)
	flagHelp = fs.Bool("h", false, "show this message and exit")
)

func init() {
	aw.IconError = iconError
	aw.IconWarning = iconWarning

	wf = aw.New()

	fs.SetOutput(os.Stderr)
}

func usage() {
	fmt.Fprint(fs.Output(), `alfred-services (-files|-services) [input...]

Alfred workflow to run macOS services

`)
	fs.PrintDefaults()
}

func pasteboardTypes() []string {
	if s := os.Getenv("PBOARD_TYPES"); s != "" {
		return strings.Split(s, "|")
	}

	var types []string
	data, err := util.Run("./PasteboardTypes.js")
	checkErr(err)
	checkErr(json.Unmarshal(data, &types))

	wf.Var("PBOARD_TYPES", strings.Join(types, "|"))
	return types
}

// Service is a macOS service.
type Service struct {
	Name    string   // name of service
	Types   []string // supported pasteboard types
	AppName string   // name of app service belongs to (optional)
	AppPath string   // path of application that defines service
}

// Title returns a more readable name.
func (s Service) Title() string {
	// Safari action has this weird name. Replace it with something better.
	if s.Name == "Search With %WebSearchProvider@" {
		return "Search Web"
	}
	return s.Name
}

// Icon returns a workflow icon for the service.
func (s Service) Icon() *aw.Icon { return &aw.Icon{Value: s.AppPath, Type: aw.IconTypeFileIcon} }

// Supports returns true if this service supports any of the given types.
func (s Service) Supports(types []string) bool {
	for _, t1 := range s.Types {
		for _, t2 := range types {
			if t1 == t2 {
				return true
			}
		}
	}
	return false
}

// read services from property list.
func loadServices() ([]Service, error) {
	var services []Service

	path := os.ExpandEnv(servicesList)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read services list")
	}
	v := struct {
		Apps map[string]struct {
			BundleID string `plist:"bundle_id"`
			Name     string `plist:"name"`
			Services []struct {
				Menu struct {
					Name string `plist:"default"`
				} `plist:"NSMenuItem"`
				Types []string `plist:"NSSendTypes"`
			} `plist:"service_dicts"`
		} `plist:"ServicesCache"`
	}{}

	if _, err := plist.Unmarshal(data, &v); err != nil {
		return nil, errors.Wrap(err, "unmarshal services list")
	}

	for path, app := range v.Apps {
		for _, v := range app.Services {
			services = append(services, Service{
				Name:    v.Menu.Name,
				Types:   v.Types,
				AppName: app.Name,
				AppPath: path,
			})
		}
	}

	return services, nil
}

func run() {
	checkErr(fs.Parse(wf.Args()))

	if *flagHelp {
		usage()
		return
	}

	var (
		query       = fs.Arg(0)
		types       []string
		allServices []Service
		services    []Service
		err         error
	)

	log.Printf("query=%q", query)

	types = pasteboardTypes()
	if len(types) == 0 {
		wf.Warn("No Data on Pasteboard", "")
		return
	}

	for _, s := range types {
		log.Printf("[pasteboard] type=%q", s)
	}

	allServices, err = loadServices()
	checkErr(err)
	log.Printf("%d total service(s)", len(allServices))

	for _, s := range allServices {
		if s.Supports(types) {
			services = append(services, s)
		}
	}
	log.Printf("%d service(s) support current pasteboard types", len(services))
	if len(services) == 0 {
		wf.Warn("No Matching Services", "No services support the current data")
		return
	}

	for _, s := range services {
		wf.NewItem(s.Title()).
			Subtitle(s.AppName).
			Arg(s.Name).
			Valid(true).
			Icon(s.Icon())
	}

	if query != "" {
		wf.Filter(query)
	}

	wf.WarnEmpty("No Matching Services", "Try a different query?")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	panic(err)
}
