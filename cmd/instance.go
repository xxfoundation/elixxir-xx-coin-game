////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package cmd initializes the CLI and config parsers as well as the logger

package cmd

import (
	"encoding/csv"
	jww "github.com/spf13/jwalterweatherman"
	"io"
	"os"
	"strconv"
	"sync"
)

type Addresses struct {
	addressMap map[string]uint64
	sync.Mutex
}

func ReadAddresses(path string) *Addresses {
	// Open the file
	csvfile, err := os.Open(path)
	if err != nil {
		jww.FATAL.Panicf("Couldn't open the csv file: %+v", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	addresses := &Addresses{
		addressMap: make(map[string]uint64),
	}

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			jww.FATAL.Panicf("Unable to read line: %+v", err)
		}
		if len(record) != 2 {
			jww.FATAL.Panicf("Line is formatted incorrectly: %v", record)
		}
		jww.DEBUG.Printf("Address: %s, Amount: %s", record[0], record[1])
		addresses.addressMap[record[0]], err = strconv.ParseUint(record[1], 10, 32)
		if err != nil {
			jww.FATAL.Panicf("Unable to process value: %+v", err)
		}
	}
	jww.INFO.Printf("Processed %d records!", len(addresses.addressMap))
	return addresses
}
