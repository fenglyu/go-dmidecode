package smbios

import (
	//	"encoding/binary"
	"fmt"
	"log"

	"github.com/digitalocean/go-smbios/smbios"
)

const (
	headerLen = 4
)

type StringKW struct {
	Keyword string
	Type    uint8
	Offset  uint8
}

// Same as defined in dmidecode
// https://github.com/mirror/dmidecode/blob/master/dmiopt.c#L150
// The Offset is calculated from the beginning of `Structure`
// While Structure's Formatted attribute is from the end of `Strucure` Header(4 BYTE)
var string_keyword = []*StringKW{
	{"bios-vendor", 0, 0x04},
	{"bios-version", 0, 0x05},
	{"bios-release-date", 0, 0x08},
	{"bios-revision", 0, 0x15},
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
	Table map[string]*StringKW
	ep    smbios.EntryPoint
	ss    []*smbios.Structure
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

	table := make(map[string]*StringKW)
	for _, v := range string_keyword {
		table[v.Keyword] = v
	}
	dt.Table = table
	return dt
}

func (dmit *DMITable) Version() string {
	// Determine SMBIOS version and table location from entry point.
	major, minor, rev := dmit.ep.Version()
	addr, size := dmit.ep.Table()

	return fmt.Sprintf("SMBIOS %d.%d.%d - table: address: %#x, size: %d\n",
		major, minor, rev, addr, size)
}

func (dmit *DMITable) GetResultByKeyword(keyword string) string {

	if _, ok := dmit.Table[keyword]; !ok {
		return ""
	}

	sk := dmit.Table[keyword]

	var s *smbios.Structure
	for _, st := range dmit.ss {
		if sk.Type == st.Header.Type {
			s = st
			break
		}
	}

	if sk.Offset-uint8(headerLen) >= s.Header.Length {
		return ""
	}

	offset := sk.Offset
	key := (s.Header.Type << 8) | offset
	switch keyword {
	case "bios-revision", "firmware-revision":
		k := key - 1 - headerLen
		if s.Formatted[k] != 0xFF && s.Formatted[k+1] != 0xFF {
			return fmt.Sprintf("%d.%d\n", s.Formatted[k], s.Formatted[k+1])
		}
		break
	case "system-uuid": /* dmi_system_uuid() */
		fmt.Println(keyword)
	case "chassis-type": /* dmi_chassis_type() */
		fmt.Println(keyword)
	case "processor-family": /* dmi_processor_family() */
		fmt.Println(keyword)
	case "processor-frequency": /* dmi_processor_frequency() */
		fmt.Println(keyword)
	default:
		return dmit.dmi_to_string(s, offset)
	}

	return ""
}

func (dmit *DMITable) dmi_to_string(s *smbios.Structure, offset uint8) string {

	offset -= uint8(headerLen)
	if int(offset) >= len(s.Strings) {
		return ""
	}

	return s.Strings[offset]
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
