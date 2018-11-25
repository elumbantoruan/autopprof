// Copyright 2018 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package autopprof provides a development-time
// library to collect pprof profiles from Go programs.
//
// This package is experimental and APIs may change.
package autopprof

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

// Profile represents a pprof profile.
type Profile interface {
	Capture() (profile string, err error)
}

// CPUProfile captures the CPU profile.
type CPUProfile struct {
	Duration time.Duration // 30 seconds by default
}

// Capture captures and writes the CPUProfile
func (p CPUProfile) Capture() (string, error) {
	dur := p.Duration
	if dur == 0 {
		dur = 30 * time.Second
	}

	f := newTemp()
	if err := pprof.StartCPUProfile(f); err != nil {
		return "", nil
	}
	time.Sleep(dur)
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		return "", nil
	}
	return f.Name(), nil
}

// HeapProfile captures the heap profile.
type HeapProfile struct{}

// Capture captures and writes the HeapProfile
func (p HeapProfile) Capture() (string, error) {
	f := newTemp()
	if err := pprof.WriteHeapProfile(f); err != nil {
		return "", nil
	}
	if err := f.Close(); err != nil {
		return "", nil
	}
	return f.Name(), nil
}

// GoRoutineProfile captures the goroutine profile.
type GoRoutineProfile struct{}

// Capture captures and writes the GoRoutine profile
func (p GoRoutineProfile) Capture() (string, error) {
	f := newTemp()

	if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", nil
	}

	return f.Name(), nil
}

// ThreadCreateProfile captures the thread create profile
type ThreadCreateProfile struct{}

// Capture captures and writes the ThreadCreateProfile
func (p ThreadCreateProfile) Capture() (string, error) {
	f := newTemp()

	if err := pprof.Lookup("threadcreate").WriteTo(f, 0); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", nil
	}

	return f.Name(), nil
}

// AllocsProfile captures the allocs profile
type AllocsProfile struct{}

// Capture captures and writes the AllocsProfile
func (p AllocsProfile) Capture() (string, error) {
	f := newTemp()

	if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", nil
	}

	return f.Name(), nil
}

// BlockProfile captures the block profile
type BlockProfile struct{}

// Capture captures and writes the BlockProfile
func (p BlockProfile) Capture() (string, error) {
	f := newTemp()

	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", nil
	}

	return f.Name(), nil
}

// MutexProfile captures the mutex profile
type MutexProfile struct{}

// Capture captures and writes the MutexProfile
func (p MutexProfile) Capture() (string, error) {
	f := newTemp()

	if err := pprof.Lookup("mutex").WriteTo(f, 0); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", nil
	}

	return f.Name(), nil
}

// Capture captures the given profiles at SIGINT
// and opens a browser with the collected profiles.
//
// Capture should be used in development-time
// and shouldn't be in production binaries.
func Capture(p Profile) {
	// TODO(jbd): As a library, we shouldn't be in the
	// business of signal handling. Provide a better way
	// trigger the capture.
	go capture(p)
}

func capture(p Profile) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT) // TODO(jbd): Add windows support.

	fmt.Println("Send SIGQUIT (CTRL+\\) to the process to capture...")

	for {
		<-c
		log.Println("Starting to capture.")

		profile, err := p.Capture()
		if err != nil {
			log.Printf("Cannot capture profile: %v", err)
		}

		// Open profile with pprof.
		log.Printf("Starting go tool pprof %v", profile)
		cmd := exec.Command("go", "tool", "pprof", "-http=:", profile)
		if err := cmd.Run(); err != nil {
			log.Printf("Cannot start pprof UI: %v", err)
		}
	}
}

func newTemp() (f *os.File) {
	f, err := ioutil.TempFile("", "profile-")
	if err != nil {
		log.Fatalf("Cannot create new temp profile file: %v", err)
	}
	return f
}
