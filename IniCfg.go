//  ---------------------------------------------------------------------------
//
//  IniCfg.go
//
//  Copyright (c) 2015, Jared Chavez.
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

package ini

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// GetSection returns a pointer to the requested IniSection object, or nil
// if an IniSection object with that name is not present within the configuration
func (this *IniCfg) GetSection(sectionName string) *IniSection {
	val, ok := this.Sections[cleanIniToken(sectionName)]
	if ok {
		return val
	}

	return VoidSection
}

// Reparse forces a config file to be re-read and all IniSections and
// IniValues to be reparsed. After parsing is complete, all hashes are also
// recomputed.
func (this *IniCfg) Reparse() {
	this.Sections = make(map[string]*IniSection, 0)
	this.keys = make([]string, 0)
	this.parseConfig()
	this.computeHash()
}

// RawString prints the ini file exactly as read in from disk,
// along with last write time and file hash.
func (this *IniCfg) RawString() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(
		"Config: %s (%s)\n",
		this.Name,
		this.ConfigVer,
	))

	buf.WriteString(fmt.Sprintf(
		"LastWriteTime: %s\n", this.ModTime,
	))

	buf.WriteString("====================================================================\n\n")

	buf.WriteString(this.Raw)

	return buf.String()
}

// String prints a human-readable text representation of the IniCfg
// object heirarchy.
func (this *IniCfg) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(
		"Config: %s (%s)\n",
		this.Name,
		this.ConfigVer,
	))

	buf.WriteString(fmt.Sprintf(
		"LastWriteTime: %s\n", this.ModTime,
	))

	for i := range this.keys {
		buf.WriteString(this.Sections[this.keys[i]].String())
	}

	return buf.String()
}

// computeHash recomputes the sha1 hash for the config file based
// on a combination of the hashes of all IniSections.
func (this *IniCfg) computeHash() {
	hash := sha1.New()

	for i := range this.keys {
		io.WriteString(hash, fmt.Sprintf(
			"%s:%s",
			this.keys[i],
			this.Sections[this.keys[i]].ConfigVer,
		))
	}

	this.ConfigVer = fmt.Sprintf("%x", hash.Sum(nil))
}

// getSection either returns a pre-existing IniSection with the given
// sectionName, or creates a new one and returns that.
func (this *IniCfg) getSection(sectionName string) *IniSection {
	secName := cleanIniToken(sectionName)

	if val, ok := this.Sections[secName]; ok {
		return val
	}

	sec := newIniSection(secName)
	this.Sections[secName] = sec
	this.keys = append(this.keys, secName)
	sort.Strings(this.keys)

	return sec
}

// parseConfig opens the ini file, parses its sections and key/val
// pairs, and recomputes hashes for all IniSections.
func (this *IniCfg) parseConfig() {
	f, err := os.Open(this.Path)
	if err != nil {
		return
	}
	defer f.Close()

	info, err := os.Stat(this.Path)
	if err != nil {
		return
	}

	this.ModTime = info.ModTime()

	var curSection *IniSection
	scanner := bufio.NewScanner(f)

	var buf bytes.Buffer

	for scanner.Scan() {
		if scanner.Err() != nil {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		buf.WriteString(fmt.Sprintf("%s\n", line))

		if len(line) < 1 ||
			strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, ";") {
			continue
		}

		section := secRegexp.FindStringSubmatch(line)
		if len(section) > 0 && curSection == nil {
			// new section
			curSection = this.getSection(section[1])
		} else if len(section) > 0 && curSection != nil {
			// next section
			curSection = this.getSection(section[1])
		} else if curSection == nil {
			// orphaned lines
			continue
		}

		// poorly formatted lines
		keyval := keyvalRegexp.FindStringSubmatch(line)
		if len(keyval) < 1 {
			continue
		}

		// curSection keyval
		curSection.AddValue(keyval[1], keyval[2])
	}

	this.Raw = string(buf.Bytes())

	// compute hash of each section
	for i := range this.keys {
		this.Sections[this.keys[i]].ComputeHash()
	}
}
