package dmidecode

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/go-smbios/smbios"
)

const (
	headerLen = 4
	outOfSpec = "<OUT OF SPEC>"
)

type StringKW struct {
	Keyword Keyword
	Type    uint8
	Offset  uint8
}

type Keyword string

const (
	KeywordBIOSVendor            = Keyword("bios-vendor")
	KeywordBIOSVersion           = Keyword("bios-version")
	KeywordBIOSReleaseDate       = Keyword("bios-release-date")
	KeywordBIOSRevision          = Keyword("bios-revision")
	KeywordFirmwareRevision      = Keyword("firmware-revision")
	KeywordSystemManufacturer    = Keyword("system-manufacturer")
	KeywordSystemProductName     = Keyword("system-product-name")
	KeywordSystemVersion         = Keyword("system-version")
	KeywordSystemSerialNumber    = Keyword("system-serial-number")
	KeywordSystemUUID            = Keyword("system-uuid")
	KeywordSystemFamily          = Keyword("system-family")
	KeywordBaseboardManufacturer = Keyword("baseboard-manufacturer")
	KeywordBaseboardProductName  = Keyword("baseboard-product-name")
	KeywordBaseboardVersion      = Keyword("baseboard-version")
	KeywordBaseboardSerialNumber = Keyword("baseboard-serial-number")
	KeywordBaseboardAssetTag     = Keyword("baseboard-asset-tag")
	KeywordChassisManufacturer   = Keyword("chassis-manufacturer")
	KeywordChassisType           = Keyword("chassis-type")
	KeywordChassisVersion        = Keyword("chassis-version")
	KeywordChassisSerialNumber   = Keyword("chassis-serial-number")
	KeywordChassisAssetTag       = Keyword("chassis-asset-tag")
	KeywordProcessorFamily       = Keyword("processor-family")
	KeywordProcessorManufacturer = Keyword("processor-manufacturer")
	KeywordProcessorVersion      = Keyword("processor-version")
	KeywordProcessorFrequency    = Keyword("processor-frequency")
)

// Same as defined in dmidecode
// https://github.com/mirror/dmidecode/blob/master/dmiopt.c#L150
// The Offset is calculated from the beginning of `Structure`
// While Structure's Formatted attribute is from the end of `Strucure` Header(4 BYTE)
var stringKeyword = []*StringKW{
	{KeywordBIOSVendor, 0, 0x04},
	{KeywordBIOSVersion, 0, 0x05},
	{KeywordBIOSReleaseDate, 0, 0x08},
	{KeywordBIOSRevision, 0, 0x15},
	{KeywordFirmwareRevision, 0, 0x17}, /* 0x16 and 0x17 */
	{KeywordSystemManufacturer, 1, 0x04},
	{KeywordSystemProductName, 1, 0x05},
	{KeywordSystemVersion, 1, 0x06},
	{KeywordSystemSerialNumber, 1, 0x07},
	{KeywordSystemUUID, 1, 0x08}, /* dmiSystemUUID() */
	{KeywordSystemFamily, 1, 0x1a},
	{KeywordBaseboardManufacturer, 2, 0x04},
	{KeywordBaseboardProductName, 2, 0x05},
	{KeywordBaseboardVersion, 2, 0x06},
	{KeywordBaseboardSerialNumber, 2, 0x07},
	{KeywordBaseboardAssetTag, 2, 0x08},
	{KeywordChassisManufacturer, 3, 0x04},
	{KeywordChassisType, 3, 0x05}, /* dmiChassisType() */
	{KeywordChassisVersion, 3, 0x06},
	{KeywordChassisSerialNumber, 3, 0x07},
	{KeywordChassisAssetTag, 3, 0x08},
	{KeywordProcessorFamily, 4, 0x06}, /* dmiProcessorFamily() */
	{KeywordProcessorManufacturer, 4, 0x07},
	{KeywordProcessorVersion, 4, 0x10},
	{KeywordProcessorFrequency, 4, 0x16}, /* dmiProcessorFrequency() */
}

type DMIType uint8

type DMITable struct {
	Table map[Keyword]*StringKW
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

	table := make(map[Keyword]*StringKW)
	for _, v := range stringKeyword {
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

func (dmit *DMITable) Query(keyword Keyword) string {
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
	// There is no such header type
	if s == nil {
		return ""
	}

	if sk.Offset >= s.Header.Length {
		return ""
	}

	offset := sk.Offset
	key := (uint16(s.Header.Type) << 8) | uint16(offset)
	switch keyword {
	case KeywordBIOSRevision, KeywordFirmwareRevision:
		key -= headerLen
		if s.Formatted[key-1] != 0xFF && s.Formatted[key] != 0xFF {
			return fmt.Sprintf("%d.%d", s.Formatted[key-1], s.Formatted[key])
		}
		break
	case KeywordSystemUUID:
		return dmit.dmiSystemUUID(s, int(offset)-headerLen)
	case KeywordChassisType:
		p := s.Formatted[int(offset)-headerLen]
		return dmit.dmiChassisType(p)
	case KeywordProcessorFamily:
		return dmit.dmiProcessorFamily(s)
	case KeywordProcessorFrequency:
		p := s.Formatted[int(offset)-headerLen:]
		return dmit.dmiProcessorFrequency(p)
	default:
		return dmit.dmiToString(s, int(offset))
	}

	return ""
}

// offset is counted from the beginning for the structure
// Structure's Formatted is count from the end of Header(4 BYTE long)
func (dmit *DMITable) dmiToString(s *smbios.Structure, offset int) string {
	offset -= headerLen
	idx := uint8(s.Formatted[offset])
	if int(idx) > len(s.Strings) || idx == 0 {
		return "Not Specified"
	}

	return s.Strings[int(idx)-1]
}

func (dmit *DMITable) dmiSystemUUID(s *smbios.Structure, offset int) string {
	only0xFF, only0x00 := true, true
	p := s.Formatted[offset:]
	for i := 0; i < 16 && (only0x00 || only0xFF); i++ {
		if p[i] != 0x00 {
			only0x00 = false
		}
		if p[i] != 0xFF {
			only0xFF = false
		}
	}
	if only0xFF {
		return fmt.Sprintln("Not Present")
	}
	if only0x00 {
		return fmt.Sprintln("Not Settable")
	}
	major, minor, _ := dmit.ep.Version()
	//fmt.Println(major, minor, rev)
	if major >= 3 || (major >= 2 && minor >= 6) {
		return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
			p[3], p[2], p[1], p[0], p[5], p[4], p[7], p[6], p[8], p[9], p[10], p[11], p[12], p[13], p[14], p[15])
	}
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x", p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7], p[8], p[9], p[10], p[11], p[12], p[13], p[14], p[15])
}

func (dmit *DMITable) dmiChassisType(code uint8) string {
	/* 7.4.1 */
	ctype := []string{
		"Other", /* 0x01 */
		"Unknown",
		"Desktop",
		"Low Profile Desktop",
		"Pizza Box",
		"Mini Tower",
		"Tower",
		"Portable",
		"Laptop",
		"Notebook",
		"Hand Held",
		"Docking Station",
		"All In One",
		"Sub Notebook",
		"Space-saving",
		"Lunch Box",
		"Main Server Chassis", /* CIM_Chassis.ChassisPackageType says "Main System Chassis" */
		"Expansion Chassis",
		"Sub Chassis",
		"Bus Expansion Chassis",
		"Peripheral Chassis",
		"RAID Chassis",
		"Rack Mount Chassis",
		"Sealed-case PC",
		"Multi-system",
		"CompactPCI",
		"AdvancedTCA",
		"Blade",
		"Blade Enclosing",
		"Tablet",
		"Convertible",
		"Detachable",
		"IoT Gateway",
		"Embedded PC",
		"Mini PC",
		"Stick PC", /* 0x24 */
	}

	code &= 0x7F /* bits 6:0 are chassis type, 7th bit is the lock bit */

	if code >= 0x01 && code <= 0x24 {
		return ctype[code-0x01]
	}
	return outOfSpec
}

func (dmit *DMITable) dmiProcessorFamily(s *smbios.Structure) string {
	data := s.Formatted
	var i, low, high int
	var code uint16

	/* 7.5.2 */
	family2 := []struct {
		value uint16
		name  string
	}{
		{0x01, "Other"},
		{0x02, "Unknown"},
		{0x03, "8086"},
		{0x04, "80286"},
		{0x05, "80386"},
		{0x06, "80486"},
		{0x07, "8087"},
		{0x08, "80287"},
		{0x09, "80387"},
		{0x0A, "80487"},
		{0x0B, "Pentium"},
		{0x0C, "Pentium Pro"},
		{0x0D, "Pentium II"},
		{0x0E, "Pentium MMX"},
		{0x0F, "Celeron"},
		{0x10, "Pentium II Xeon"},
		{0x11, "Pentium III"},
		{0x12, "M1"},
		{0x13, "M2"},
		{0x14, "Celeron M"},
		{0x15, "Pentium 4 HT"},

		{0x18, "Duron"},
		{0x19, "K5"},
		{0x1A, "K6"},
		{0x1B, "K6-2"},
		{0x1C, "K6-3"},
		{0x1D, "Athlon"},
		{0x1E, "AMD29000"},
		{0x1F, "K6-2+"},
		{0x20, "Power PC"},
		{0x21, "Power PC 601"},
		{0x22, "Power PC 603"},
		{0x23, "Power PC 603+"},
		{0x24, "Power PC 604"},
		{0x25, "Power PC 620"},
		{0x26, "Power PC x704"},
		{0x27, "Power PC 750"},
		{0x28, "Core Duo"},
		{0x29, "Core Duo Mobile"},
		{0x2A, "Core Solo Mobile"},
		{0x2B, "Atom"},
		{0x2C, "Core M"},
		{0x2D, "Core m3"},
		{0x2E, "Core m5"},
		{0x2F, "Core m7"},
		{0x30, "Alpha"},
		{0x31, "Alpha 21064"},
		{0x32, "Alpha 21066"},
		{0x33, "Alpha 21164"},
		{0x34, "Alpha 21164PC"},
		{0x35, "Alpha 21164a"},
		{0x36, "Alpha 21264"},
		{0x37, "Alpha 21364"},
		{0x38, "Turion II Ultra Dual-Core Mobile M"},
		{0x39, "Turion II Dual-Core Mobile M"},
		{0x3A, "Athlon II Dual-Core M"},
		{0x3B, "Opteron 6100"},
		{0x3C, "Opteron 4100"},
		{0x3D, "Opteron 6200"},
		{0x3E, "Opteron 4200"},
		{0x3F, "FX"},
		{0x40, "MIPS"},
		{0x41, "MIPS R4000"},
		{0x42, "MIPS R4200"},
		{0x43, "MIPS R4400"},
		{0x44, "MIPS R4600"},
		{0x45, "MIPS R10000"},
		{0x46, "C-Series"},
		{0x47, "E-Series"},
		{0x48, "A-Series"},
		{0x49, "G-Series"},
		{0x4A, "Z-Series"},
		{0x4B, "R-Series"},
		{0x4C, "Opteron 4300"},
		{0x4D, "Opteron 6300"},
		{0x4E, "Opteron 3300"},
		{0x4F, "FirePro"},
		{0x50, "SPARC"},
		{0x51, "SuperSPARC"},
		{0x52, "MicroSPARC II"},
		{0x53, "MicroSPARC IIep"},
		{0x54, "UltraSPARC"},
		{0x55, "UltraSPARC II"},
		{0x56, "UltraSPARC IIi"},
		{0x57, "UltraSPARC III"},
		{0x58, "UltraSPARC IIIi"},

		{0x60, "68040"},
		{0x61, "68xxx"},
		{0x62, "68000"},
		{0x63, "68010"},
		{0x64, "68020"},
		{0x65, "68030"},
		{0x66, "Athlon X4"},
		{0x67, "Opteron X1000"},
		{0x68, "Opteron X2000"},
		{0x69, "Opteron A-Series"},
		{0x6A, "Opteron X3000"},
		{0x6B, "Zen"},

		{0x70, "Hobbit"},

		{0x78, "Crusoe TM5000"},
		{0x79, "Crusoe TM3000"},
		{0x7A, "Efficeon TM8000"},

		{0x80, "Weitek"},

		{0x82, "Itanium"},
		{0x83, "Athlon 64"},
		{0x84, "Opteron"},
		{0x85, "Sempron"},
		{0x86, "Turion 64"},
		{0x87, "Dual-Core Opteron"},
		{0x88, "Athlon 64 X2"},
		{0x89, "Turion 64 X2"},
		{0x8A, "Quad-Core Opteron"},
		{0x8B, "Third-Generation Opteron"},
		{0x8C, "Phenom FX"},
		{0x8D, "Phenom X4"},
		{0x8E, "Phenom X2"},
		{0x8F, "Athlon X2"},
		{0x90, "PA-RISC"},
		{0x91, "PA-RISC 8500"},
		{0x92, "PA-RISC 8000"},
		{0x93, "PA-RISC 7300LC"},
		{0x94, "PA-RISC 7200"},
		{0x95, "PA-RISC 7100LC"},
		{0x96, "PA-RISC 7100"},

		{0xA0, "V30"},
		{0xA1, "Quad-Core Xeon 3200"},
		{0xA2, "Dual-Core Xeon 3000"},
		{0xA3, "Quad-Core Xeon 5300"},
		{0xA4, "Dual-Core Xeon 5100"},
		{0xA5, "Dual-Core Xeon 5000"},
		{0xA6, "Dual-Core Xeon LV"},
		{0xA7, "Dual-Core Xeon ULV"},
		{0xA8, "Dual-Core Xeon 7100"},
		{0xA9, "Quad-Core Xeon 5400"},
		{0xAA, "Quad-Core Xeon"},
		{0xAB, "Dual-Core Xeon 5200"},
		{0xAC, "Dual-Core Xeon 7200"},
		{0xAD, "Quad-Core Xeon 7300"},
		{0xAE, "Quad-Core Xeon 7400"},
		{0xAF, "Multi-Core Xeon 7400"},
		{0xB0, "Pentium III Xeon"},
		{0xB1, "Pentium III Speedstep"},
		{0xB2, "Pentium 4"},
		{0xB3, "Xeon"},
		{0xB4, "AS400"},
		{0xB5, "Xeon MP"},
		{0xB6, "Athlon XP"},
		{0xB7, "Athlon MP"},
		{0xB8, "Itanium 2"},
		{0xB9, "Pentium M"},
		{0xBA, "Celeron D"},
		{0xBB, "Pentium D"},
		{0xBC, "Pentium EE"},
		{0xBD, "Core Solo"},
		/* 0xBE handled as a special case */
		{0xBF, "Core 2 Duo"},
		{0xC0, "Core 2 Solo"},
		{0xC1, "Core 2 Extreme"},
		{0xC2, "Core 2 Quad"},
		{0xC3, "Core 2 Extreme Mobile"},
		{0xC4, "Core 2 Duo Mobile"},
		{0xC5, "Core 2 Solo Mobile"},
		{0xC6, "Core i7"},
		{0xC7, "Dual-Core Celeron"},
		{0xC8, "IBM390"},
		{0xC9, "G4"},
		{0xCA, "G5"},
		{0xCB, "ESA/390 G6"},
		{0xCC, "z/Architecture"},
		{0xCD, "Core i5"},
		{0xCE, "Core i3"},
		{0xCF, "Core i9"},

		{0xD2, "C7-M"},
		{0xD3, "C7-D"},
		{0xD4, "C7"},
		{0xD5, "Eden"},
		{0xD6, "Multi-Core Xeon"},
		{0xD7, "Dual-Core Xeon 3xxx"},
		{0xD8, "Quad-Core Xeon 3xxx"},
		{0xD9, "Nano"},
		{0xDA, "Dual-Core Xeon 5xxx"},
		{0xDB, "Quad-Core Xeon 5xxx"},

		{0xDD, "Dual-Core Xeon 7xxx"},
		{0xDE, "Quad-Core Xeon 7xxx"},
		{0xDF, "Multi-Core Xeon 7xxx"},
		{0xE0, "Multi-Core Xeon 3400"},

		{0xE4, "Opteron 3000"},
		{0xE5, "Sempron II"},
		{0xE6, "Embedded Opteron Quad-Core"},
		{0xE7, "Phenom Triple-Core"},
		{0xE8, "Turion Ultra Dual-Core Mobile"},
		{0xE9, "Turion Dual-Core Mobile"},
		{0xEA, "Athlon Dual-Core"},
		{0xEB, "Sempron SI"},
		{0xEC, "Phenom II"},
		{0xED, "Athlon II"},
		{0xEE, "Six-Core Opteron"},
		{0xEF, "Sempron M"},

		{0xFA, "i860"},
		{0xFB, "i960"},

		{0x100, "ARMv7"},
		{0x101, "ARMv8"},
		{0x104, "SH-3"},
		{0x105, "SH-4"},
		{0x118, "ARM"},
		{0x119, "StrongARM"},
		{0x12C, "6x86"},
		{0x12D, "MediaGX"},
		{0x12E, "MII"},
		{0x140, "WinChip"},
		{0x15E, "DSP"},
		{0x1F4, "Video Processor"},

		{0x200, "RV32"},
		{0x201, "RV64"},
		{0x202, "RV128"},
	}

	major, minor, _ := dmit.ep.Version()
	/* Special case for ambiguous value 0x30 (SMBIOS 2.0 only) */
	if major == 2 && minor == 0 && data[0x06-headerLen] == 0x30 && s.Header.Length >= 0x08 {
		manufacturer := dmit.dmiToString(s, 0x07)

		if strings.Contains(manufacturer, "Intel") || strings.EqualFold(manufacturer[:5], "Intel") {
		}
		return "Pentium Pro"
	}

	code = uint16(data[0x06-headerLen])
	if data[0x06-headerLen] == 0xFE && s.Header.Length >= 0x2A {
		code = binary.LittleEndian.Uint16(data[0x28-headerLen : 0x28-headerLen+2])
	}

	/* Special case for ambiguous value 0xBE */
	if code == 0xBE {
		if s.Header.Length >= 0x08 {
			manufacturer := dmit.dmiToString(s, 0x07)

			/* Best bet based on manufacturer string */
			if strings.Contains(manufacturer, "Intel") || strings.EqualFold(manufacturer[:5], "Intel") {
				return "Core 2"
			}
			if strings.Contains(manufacturer, "AMD") || strings.EqualFold(manufacturer[:3], "AMD") {
				return "K7"
			}
			return "Core 2 or K7"
		}
	}

	/* Perform a binary search */
	low = 0
	high = len(family2) - 1

	for {
		i = (low + high) / 2

		if family2[i].value == code {
			return family2[i].name
		}

		if low == high { /* Not found */
			return outOfSpec
		}
		if code < family2[i].value {
			high = i
		} else {
			low = i + 1
		}
	}
}

func (dmit *DMITable) dmiProcessorFrequency(p []byte) string {
	code := binary.LittleEndian.Uint16(p[:2])
	if code != 0 {
		return fmt.Sprintf("%d MHz", code)
	}
	return "Unknown"
}
