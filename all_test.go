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
    "fmt"
    "testing"
)


func TestIni(t *testing.T) {
    cfg := New("./test.ini")

    // test ini properties
    if cfg.Name != "test" {
        t.Error("expected test, got " + cfg.Name)
    }

    if cfg.Path != "./test.ini" {
        t.Error("expected ./test.ini, got " + cfg.Path)
    }

    fmt.Println(cfg)
}
