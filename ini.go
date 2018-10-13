//  ---------------------------------------------------------------------------
//
//  ini.go
//
//  Copyright (c) 2015, Jared Chavez.
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

// Package ini provides basic facilities for parsing and monitoring
// changes to ini format configuration files.
package ini

import (
	"github.com/xaevman/crash"
	"github.com/xaevman/shutdown"

	"path"
	"regexp"
	"strings"
	"time"
)

// DefaultPollFreqSec is the default amount of time, in seconds,
// that the ini system will wait between checks to see if a given
// ini file has changed.
const DefaultPollFreqSec = 10

// Regex to parse lines containing section definitions.
const sectionRegexFmt = "^\\s*\\[(.*)\\]\\s*$"

// Regex to parse key/value lines.
const keyvalRegexFmt = "^\\s*(.*?)\\s*=\\s*(.*?)\\s*$"

// Shared regexp objects.
var (
	secRegexp    = regexp.MustCompile(sectionRegexFmt)
	keyvalRegexp = regexp.MustCompile(keyvalRegexFmt)
)

var iniShutdown = shutdown.New()

// IniCfg represents a single ini configuration file, containing pointers
// to the IniSections contained within it. ConfigVer is a consistent hash
// of all IniSections within the file which is not influenced by whitespace
// or comments.
type IniCfg struct {
	ConfigVer string
	Name      string
	Path      string
	ModTime   time.Time
	Sections  map[string]*IniSection
	Raw       string

	keys []string
}

// IniSection represents a section within an ini file. Section names are
// enclosed with square brackets (ex: [my.section.name]). Section names
// are case insensitive and cleaned up by forced removal of framing white
// space. White space within the section name is converted to an underline.
// There may be multiple instances of a given section name present within
// the ini file, but only one IniSection object will be created. In the case
// where multiple instances of an ini section exist within the file, the
// child key/value objects will be coalesced within the single IniSection
// object. ConfigVer is a consistent hash of all the Key/Value entries
// within the section which is not influenced by whitespace or comments.
type IniSection struct {
	ConfigVer string
	Name      string
	Values    map[string][]*IniValue

	keys []string
}

// IniValue represents a single Key/Value within a config section. Multiple
// instances of of a given key may be present within a section, and separate
// IniValue objects will exist for each instance of the key. Keys and values
// are separated by an equal sign. The value side of the key/value pair is
// split on a comma delimeter and trimmed of any enclosing whitepsace.
type IniValue struct {
	Name   string
	Values []string
}

// New returns a pointer to a new IniCfg object for the given file path.
func New(iniPath string) *IniCfg {
	return newIniCfg(iniPath)
}

func Shutdown() {
	iniShutdown.Start()
	if iniShutdown.WaitForTimeout() {
		panic("Shutdown Timeout")
	}
}

// newIniCfg returns a pointer to a new IniCfg object for the given file path.
func newIniCfg(iniPath string) *IniCfg {
	iniName := strings.TrimSuffix(
		path.Base(iniPath),
		path.Ext(iniPath),
	)

	cfg := IniCfg{
		Name: iniName,
		Path: iniPath,
	}

	cfg.Reparse()

	return &cfg
}

// newIniSection returns a pointer to a new IniSection object for the named
// section.
func newIniSection(sectionName string) *IniSection {
	sec := IniSection{
		Name:   cleanIniToken(sectionName),
		Values: make(map[string][]*IniValue, 0),
		keys:   make([]string, 0),
	}

	return &sec
}

// newIniValue returns a pointer to a new IniValue object for the given
// key.
func newIniValue(key, valstring string) *IniValue {
	val := IniValue{
		Name: cleanIniToken(key),
	}

	val.parseValues(valstring)

	return &val
}

// cleanIniToken takes a given token (usually a section name or config key)
// and forces it to lower case, trims enclosing white space, and converts any
// interior spaces to underscores.
func cleanIniToken(token string) string {
	clean := strings.ToLower(strings.TrimSpace(token))
	clean = strings.Replace(clean, " ", "_", -1)

	return clean
}

// init initializes the ini package. Primarily, it spawns a goroutine
// which is responsible for handling ini change monitoring.
func init() {
	SetPollFreqSec(DefaultPollFreqSec)
	go func() {
		defer crash.HandleAll()
		monitorInis()
	}()
}
