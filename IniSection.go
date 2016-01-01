//  ---------------------------------------------------------------------------
//
//  IniSection.go
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
    "bytes"
    "crypto/sha1"
    "fmt"
    "io"
    "sort"
)

// VoidSection is returned by GetSection so that subsequent calls to GetVal
// can be made without nil checking.
var VoidSection = newIniSection("void")

// AddValue adds a new IniValue object for the given key and value strings
// to the IniSection instance.
func (this *IniSection) AddValue(key, value string) {
    ckey   := cleanIniToken(key)
    newVal := newIniValue(ckey, value)

    if _, ok := this.Values[ckey]; ok {
        this.Values[ckey] = append(this.Values[ckey], newVal)
    } else {
        newValArray             := make([]*IniValue, 1)
        newValArray[0]           = newVal
        this.Values[newVal.Name] = newValArray
        this.keys                = append(this.keys, newVal.Name)
        sort.Strings(this.keys)
    }
}

// ComputeHash recomputes the sha1 hash of all the key/val pairs within
// the IniSection.
func (this *IniSection) ComputeHash() {
    hash := sha1.New()

    for x := range this.keys {
        for y := range this.Values[this.keys[x]] {
            io.WriteString(hash, fmt.Sprintf(
                "%s:%s",
                this.keys[x],
                this.Values[this.keys[x]][y],
            ))
        }
    }

    this.ConfigVer = fmt.Sprintf("%x", hash.Sum(nil))
}

// GetFirstVal returns a pointer to the first IniValue object with a matching
// key name. Returns nil if no relevant IniValue objects are present in the 
// section.
func (this *IniSection) GetFirstVal(valName string) *IniValue {
    vals := this.GetVals(valName)

    if len(vals) < 1 {
        return VoidValue
    }

    return vals[0]
}

// GetVals returns an array of pointers to IniValue objects with matching key
// names. Returns nil if no relevant IniValue objects are present in the section.
func (this *IniSection) GetVals(valName string) []*IniValue {
    vals, ok := this.Values[cleanIniToken(valName)]
    if ok {
        return vals
    }

    return make([]*IniValue, 0)
}

// String prints a human-readable representation of the IniSection and
// its children IniValue objects.
func (this *IniSection) String() string {
    var buf bytes.Buffer
    buf.WriteString(fmt.Sprintf(
        "Section :: %s (%s)\n",
        this.Name,
        this.ConfigVer,
    ))

    for x := range this.keys {
        for y := range this.Values[this.keys[x]] {
            buf.WriteString(this.Values[this.keys[x]][y].String())
        }
    }

    return buf.String()
}

