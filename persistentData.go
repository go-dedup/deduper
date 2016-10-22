////////////////////////////////////////////////////////////////////////////
// Porgram: PersistentData.go
// Purpose: Go Persistent Data to/from disk
// Authors: Tong Sun (c) 2016, All rights reserved
// Credits:
//          https://play.golang.org/p/wT8_H44crC by Michael Jones
////////////////////////////////////////////////////////////////////////////

package main

import (
	"encoding/gob"
	"os"
)

func SaveState(persistName string, state interface{}) error {
	// create persistence file
	f, err := os.Create(persistName)
	if err != nil {
		return err
	}
	defer f.Close()

	// write persistemce file
	e := gob.NewEncoder(f).Encode(state)
  err = f.Close()
	if err != nil {
		return err
	}
	return e
}

func RestoreState(persistName string, state interface{}) error {
	// open persistence file
	f, err := os.Open(persistName)
	if err != nil {
		return err
	}
	defer f.Close()

	// read persistemce file
	err = gob.NewDecoder(f).Decode(state)
	return err
}
