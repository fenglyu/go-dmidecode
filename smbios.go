package smbios

import (
	"fmt"
	"log"

	"github.com/digitalocean/go-smbios/smbios"
)

type StringKW struct {
	Keyword string
	Type    uint8
	Offset  uint8
}

var string_keyword = []*StringKW{
	{"bios-vendor", 0, 0x04},
	{"bios-version", 0, 0x05},
	{"bios-release-date", 0, 0x08},
	{"bios-revision", 0, 0x15},     /* 0x14 and 0x15 */
	{"firmware-revision", 0, 0x17}, /* 0x16 and 0x17 */
	{"system-manufacturer", 1, 0x04},
	{"system-product-name", 1, 0x05},
	{"system-version", 1, 0x06},
	{"system-serial-number", 1, 0x07},
	{"system-uuid", 1, 0x08}, /* dmi_system_uuid() */
	{"system-family", 1, 0x1a},
	{"baseboard-manufacturer", 2, 0x04},
	{"baseboard-product-name", 2, 0x05},
	{"baseboard-version", 2, 0x06},
	{"baseboard-serial-number", 2, 0x07},
	{"baseboard-asset-tag", 2, 0x08},
	{"chassis-manufacturer", 3, 0x04},
	{"chassis-type", 3, 0x05}, /* dmi_chassis_type() */
	{"chassis-version", 3, 0x06},
	{"chassis-serial-number", 3, 0x07},
	{"chassis-asset-tag", 3, 0x08},
	{"processor-family", 4, 0x06}, /* dmi_processor_family() */
	{"processor-manufacturer", 4, 0x07},
	{"processor-version", 4, 0x10},
	{"processor-frequency", 4, 0x16}, /* dmi_processor_frequency() */
}

type DMIType uint8

type DMITable struct {
	Mapping map[string]*StringKW
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

	mapping := make(map[string]*StringKW)
	for _, v := range string_keyword {
		mapping[v.Keyword] = v
	}
	dt.Mapping = mapping
	return dt
}

func (dmit *DMITable) Version() string {
	// Determine SMBIOS version and table location from entry point.
	major, minor, rev := dmit.ep.Version()
	addr, size := dmit.ep.Table()

	return fmt.Sprintf("SMBIOS %d.%d.%d - table: address: %#x, size: %d\n",
		major, minor, rev, addr, size)
}

//func (dmit *DMITable)

/*
bios-vendor, bios-ver‐sion, bios-release-date, system-manufacturer, system-product-name, system-version, system-serial-number, system-uuid, baseboard-manufacturer,baseboard-product-name, baseboard-version, baseboard-serial-number, baseboard-asset-tag, chassis-manufacturer, chassis-type, chassis-version, chassis-serial-number, chassis-asset-tag, processor-family, processor-manufacturer,  processor-version,  processor-frequency.
*/
func (dmit *DMITable) GetResultByKeyword(keyword string) string {

	if _, ok := dmit.Mapping[keyword]; ok == nil {
		return ""
	}

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

	return ""
}

func (dmit *DMITable) GetEntriesByType(et DMIType) string {
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
