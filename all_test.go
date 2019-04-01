//  ---------------------------------------------------------------------------
//
//  all_test.go
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
    "os"
    "testing"
    "time"
)

var stopTest = make(chan bool, 1)


func TestIni(t *testing.T) {
    cfg := New("./test.ini")

    // test ini properties
    if cfg.Name != "test" {
        t.Error("expected test, got " + cfg.Name)
    }

    if len(cfg.Paths) != 1 {
        t.Errorf("expected cfg.Paths to have 1 entry, got %d", len(cfg.Paths))
    }

    if cfg.Paths[0] != "./test.ini" {
        t.Error("expected ./test.ini, got " + cfg.Paths[0])
    }

    t.Log(cfg)
}

func TestMultiFileIni(t *testing.T) {
    cfg := newIniCfgFromFiles([]string{"./test.ini", "./test2.ini"})
    
    if cfg.Name != "test" {
        t.Error("expected test, got " + cfg.Name)
    }

	sec := cfg.GetSection("section_1")
	if sec == VoidSection {
		t.Errorf("Expected a valid section (from the base file).")
	}

	// This value comes from the main ini.
	if sec.GetFirstVal("key1").GetValStr(0, "") != "value1" {
		t.Errorf("Expected key1 to have value1 (from the base file).")
	}

	// This value comes from the secondary ini.
	if sec.GetFirstVal("key3").GetValStr(0, "") != "value3" {
		t.Errorf("Expected key3 to have value3 (from the merged file).")
	}

	// This section comes from the secondary ini.
	if cfg.GetSection("section3") == VoidSection {
		t.Errorf("Expected a valid section (from the merged file).")
	}
}

func TestIniMonitor(t *testing.T) {
    cfg := New("./test.ini")
    Subscribe(cfg, onCfgChange)
    err := touch(cfg.Paths[0])
    if err != nil {
        t.Fatal(err)
    }
    
    select {
    case <-stopTest:
        t.Log("Config change notified successfully")
    case <-time.After(2 * time.Second):
        t.Fatal("Config change notification not received after 2 seconds")
    }
}

func TestIniMonitorOnSecondFile(t *testing.T) {
    cfg := newIniCfgFromFiles([]string{"./test.ini", "./test2.ini"})
    Subscribe(cfg, onCfgChange)
    err := touch(cfg.Paths[1])
    if err != nil {
        t.Fatal(err)
    }
    
    select {
    case <-stopTest:
        t.Log("Config change notified successfully")
    case <-time.After(2 * time.Second):
        t.Fatal("Config change notification not received after 2 seconds")
    }
}

func onCfgChange(cfg *IniCfg, changeCount int) {
    stopTest<- true
}

func touch(path string) error {
    now := time.Now()
    return os.Chtimes(path, now, now)
}
