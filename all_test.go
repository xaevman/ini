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

var stopTest = make(chan bool, 0)


func TestIni(t *testing.T) {
    cfg := New("./test.ini")

    // test ini properties
    if cfg.Name != "test" {
        t.Error("expected test, got " + cfg.Name)
    }

    if cfg.Path != "./test.ini" {
        t.Error("expected ./test.ini, got " + cfg.Path)
    }

    t.Log(cfg)
}

func TestIniMonitor(t *testing.T) {
    cfg := New("./test.ini")
    Subscribe(cfg, onCfgChange)
    err := touch(cfg.Path)
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
