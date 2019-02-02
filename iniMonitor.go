//  ---------------------------------------------------------------------------
//
//  iniMonitor.go
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
	"sync"
	"sync/atomic"
	"time"
)

var (
	coreId          uint32
	forceUpdateChan = make(chan interface{}, 0)
	monitors        map[string]*monIni
	monitorLock     sync.Mutex
	pollFreqSec     uint32
)

type monIni struct {
	changeCount int
	iniFile     *IniCfg
	name        string
	subscribers map[uint32]*monSubscriber
}

type monSubscriber struct {
	callback func(*IniCfg, int)
	id       uint32
}

// ClearSubscribers clears the list of subscriber functions for the given
// ini file.
func ClearSubscribers(cfg *IniCfg) {
	monitorLock.Lock()
	defer monitorLock.Unlock()

	mon := getMonIni(cfg)
	mon.subscribers = make(map[uint32]*monSubscriber)
}

// ForceUpdate cause an ini poll to happen via manual request.
func ForceUpdate() {
	forceUpdateChan <- nil
}

// SetPollFreqSec sets the frequency at which the underlying ini file is
// checked for changes.
func SetPollFreqSec(freqSec uint32) {
	atomic.StoreUint32(&pollFreqSec, freqSec)
}

// Subscribe notifies the ini system that the given callback should be called
// upon any changes to the given ini file.
func Subscribe(cfg *IniCfg, callback func(*IniCfg, int)) uint32 {
	newId := atomic.AddUint32(&coreId, 1)

	monitorLock.Lock()
	defer monitorLock.Unlock()

	mon := getMonIni(cfg)

	sub := &monSubscriber{
		callback: callback,
		id:       newId,
	}

	mon.subscribers[sub.id] = sub

	callback(mon.iniFile, 0)

	return newId
}

// Unsubscribe removes the callback identified by id from the list of
// subscribers for the given ini file.
func Unsubscribe(cfg *IniCfg, id uint32) {
	monitorLock.Lock()
	defer monitorLock.Unlock()

	mon := getMonIni(cfg)
	delete(mon.subscribers, id)
}

// getMonIni creates, or gets, a monIni object for tracking callbacks
// associated with the given ini file.
func getMonIni(cfg *IniCfg) *monIni {
	if monitors == nil {
		monitors = make(map[string]*monIni)
	}

	mon, ok := monitors[cfg.Path]
	if !ok {
		mon = &monIni{
			name:        cfg.Path,
			iniFile:     cfg,
			subscribers: make(map[uint32]*monSubscriber),
		}

		monitors[cfg.Path] = mon
	}

	return mon
}

// monitorInis is executed within a separate goroutine to periodically
// stat monitored ini files. Changes are detected by comparing an ini file's
// ModTime with the last recorded ModTime (at the time the file was last
// parsed). If changes are detected, all subscribed callbacks are called.
func monitorInis() {
	defer iniShutdown.Complete()

	for {
		monitorLock.Lock()

		for k1 := range monitors {
			mon := monitors[k1]
			info, err := os.Stat(mon.iniFile.Path)
			if err != nil {
				continue
			}

			// file hasn't been written to
			if mon.iniFile.ModTime.After(info.ModTime()) ||
				mon.iniFile.ModTime.Equal(info.ModTime()) {
				continue
			}

			// something has changed - reparse
			mon.changeCount++
			mon.iniFile.Reparse()

			// notify
			for k2 := range mon.subscribers {
				sub := mon.subscribers[k2]
				sub.callback(mon.iniFile, mon.changeCount)
			}
		}

		monitorLock.Unlock()

		select {
		case <-forceUpdateChan:
		case <-time.After(time.Duration(atomic.LoadUint32(&pollFreqSec)) * time.Second):
		case <-iniShutdown.Signal:
			return
		}
	}
}
