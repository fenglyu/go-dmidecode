package smbios

import (
	"fmt"
	"log"

	"github.com/digitalocean/go-smbios/smbios"
)

type DMIType uint8

type DMITable struct {
	mapping map[uint8]*smbios.Structure
	ep      smbios.EntryPoint
	ss      []*smbios.Structure
}

func NewDMITable() *DMITable {

	dt := &DMITable{}
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
	dt.ep = ep
	dt.ss = ss
	return dt
}

func (dmit *DMITable) Version() string {
	// Determine SMBIOS version and table location from entry point.
	major, minor, rev := dmit.ep.Version()
	addr, size := dmit.ep.Table()

	return fmt.Sprintf("SMBIOS %d.%d.%d - table: address: %#x, size: %d\n",
		major, minor, rev, addr, size)
}

/*
bios-vendor, bios-ver‐sion, bios-release-date, system-manufacturer, system-product-name, system-version, system-serial-number, system-uuid, baseboard-manufacturer,baseboard-product-name, baseboard-version, baseboard-serial-number, baseboard-asset-tag, chassis-manufacturer, chassis-type, chassis-version, chassis-serial-number, chassis-asset-tag, processor-family, processor-manufacturer,  processor-version,  processor-frequency.
*/
func (dmit *DMITable) GetResultByKeyword(keyword string) interface{} {
	switch keyword {
	case "bios-vendor":
		fmt.Println(keyword)
	case "bios-version":
		fmt.Println(keyword)
	case "bios-release-date":
		fmt.Println(keyword)
	case "system-product-name":
		fmt.Println(keyword)
	}

	return nil
}

func (dmit *DMITable) GetEntriesByType(et DMIType) string {

	major, minor, rev := dmit.ep.Version()
	addr, size := dmit.ep.Table()

	fmt.Printf("SMBIOS %d.%d.%d - table: address: %#x, size: %d\n",
		major, minor, rev, addr, size)

	for _, s := range dmit.ss {
		//	fmt.Printf("---> %d\n", s.Header.Type)
		if s.Header.Type == uint8(et) {
			var con_str string
			for _, str := range s.Strings {
				con_str = fmt.Sprintf("%s\n%s", con_str, str)
			}
			return con_str
		}
	}
	return ""
}
