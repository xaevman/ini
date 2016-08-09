//  ---------------------------------------------------------------------------
//
//  IniValue.go
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
    "fmt"
    "strconv"
    "strings"
)

// VoidValue is returned by GetFirstValue so that subsequent calls to GetValx
// can be made without nil checking.
var VoidValue = newIniValue("void", "")

// GetValBool retrieves the value at the given offset and attempts to parse
// and return it as a boolean value. If either the offset is invalid, or
// parsing fails, the supplied default value is returned.
func (this *IniValue) GetValBool(offset int, defVal bool) bool {
    if offset >= len(this.Values) {
        return defVal
    }

    bVal, err := strconv.ParseBool(this.Values[offset])
    if err != nil {
        return defVal
    }

    return bVal
}

// GetValFloat retrieves the value at the given offset and attempts to parse
// and return it as a 32bit float. If either the offset is invalid, or parsing
// fails, the supplied default value is returned.
func (this *IniValue) GetValFloat(offset int, defVal float32) float32 {
    if offset >= len(this.Values) {
        return defVal
    }

    fVal, err := strconv.ParseFloat(this.Values[offset], 32)
    if err != nil {
        return defVal
    }

    return float32(fVal)
}

// GetValFloat64 retrieves the value at the given offset and attempts to parse
// and return it as a 64bit float. If either the offset is invalid, or parsing
// fails, the supplied default value is returned.
func (this *IniValue) GetValFloat64(offset int, defVal float64) float64 {
    if offset >= len(this.Values) {
        return defVal
    }

    fVal, err := strconv.ParseFloat(this.Values[offset], 64)
    if err != nil {
        return defVal
    }

    return fVal
}

// GetValInt retrieves the value at the given offset and attempts to parse
// and return it as a 32bit signed integer. If either the offset is invalid,
// or parsing fails, the supplied default value is returned.
func (this *IniValue) GetValInt(offset int, defVal int) int {
    if offset >= len(this.Values) {
        return defVal
    }

    iVal, err := strconv.ParseInt(this.Values[offset], 10, 32)
    if err != nil {
        return defVal
    }

    return int(iVal)
}

// GetValInt64 retrieves the value at the given offset and attempts to parse
// and return it as a 64bit signed integer. If either the offset is invalid,
// or parsing fails, the supplied default value is returned.
func (this *IniValue) GetValInt64(offset int, defVal int64) int64 {
    if offset >= len(this.Values) {
        return defVal
    }

    iVal, err := strconv.ParseInt(this.Values[offset], 10, 64)
    if err != nil {
        return defVal
    }

    return iVal
}

// GetValStr retrieves the value at the given offset and returns it as a
// string value. If teh offset is invalid the supplied default value is
// returned.
func (this *IniValue) GetValStr(offset int, defVal string) string {
    if offset >= len(this.Values) {
        return defVal
    }

    if this.Values[offset] == "" {
        return defVal
    }

    return this.Values[offset]
}

// GetValUint retrieves the value at the given offset and attempts to parse
// and return it as a 32bit unsigned integer. If either the offset is invalid,
// or the parsing fails, the supplied default value is returned.
func (this *IniValue) GetValUint(offset int, defVal uint) uint {
    if offset >= len(this.Values) {
        return defVal
    }

    uVal, err := strconv.ParseUint(this.Values[offset], 10, 32)
    if err != nil {
        return defVal
    }

    return uint(uVal)
}

// GetValUint64 retrieves the value at the given offset and attempts to parse
// and return it as a 64bit unsigned integer. If either the offset is invalid,
// or the parsing fails, the supplied default value is returned.
func (this *IniValue) GetValUint64(offset int, defVal uint64) uint64 {
    if offset >= len(this.Values) {
        return defVal
    }

    uVal, err := strconv.ParseUint(this.Values[offset], 10, 64)
    if err != nil {
        return defVal
    }

    return uVal
}

// String prints a more human-readable representation of the IniValue
// object.
func (this *IniValue) String() string {
    var buf bytes.Buffer
    buf.WriteString(fmt.Sprintf(
        "[Key: %s",
        this.Name,
    ))

    for i := range this.Values {
        buf.WriteString(fmt.Sprintf(
            " | Val(%d): %s",
            i,
            this.Values[i],
        ))
    }

    buf.WriteString("]\n")

    return buf.String()
}

// parseValues strips any line comments from the value string,
// splits the raw string into its individual comma-separated parts,
// and trims any enclosing whitespace before adding the array of values
// to the IniValue object.
func (this *IniValue) parseValues(valstring string) {
    valstring = stripEOLComment(valstring)

    // parse out delimiters
    valparts := strings.Split(valstring, ",")
    this.Values = make([]string, len(valparts))

    for i := range valparts {
        this.Values[i] = strings.TrimSpace(valparts[i])
    }
}

// stripEOLComment strips the comment section of any inline comment
// present in a value string.
func stripEOLComment(value string) string {
    // look for trailing line comments
    idx := strings.Index(value, "#")
    if idx < 0 {
        return strings.TrimSpace(value)
    }

    return strings.TrimSpace(value[:idx])
}
