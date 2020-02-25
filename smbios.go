package main

import (
	"fmt"
	"log"

	"github.com/digitalocean/go-smbios/smbios"
)

type DMIType uint8

type DMITable struct {
	Type    uint8
	Keyword string
}

func (dmi *DMITable) GetEntriesByType(entry_type DMIType) string {
	// Find SMBIOS data in operating system-specific location.
	rc, ep, err := smbios.Stream()
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}
	// Be sure to close the stream!
	defer rc.Close()

	// Decode SMBIOS structures from the stream.
	d := smbios.NewDecoder(rc)
	ss, err := d.Decode()
	if err != nil {
		log.Fatalf("failed to decode structures: %v", err)
	}

	// Determine SMBIOS version and table location from entry point.
	major, minor, rev := ep.Version()
	addr, size := ep.Table()

	fmt.Printf("SMBIOS %d.%d.%d - table: address: %#x, size: %d\n",
		major, minor, rev, addr, size)

	for _, s := range ss {
		//	fmt.Printf("---> %d\n", s.Header.Type)
		if s.Header.Type == dtype {
			return s.Strings
		}
	}
}
