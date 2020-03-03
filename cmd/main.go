package main

import (
	//	"encoding/binary"
	"fmt"
	//"log"
	//"github.com/digitalocean/go-smbios/smbios"
	"github.com/fenglyu/smbios"
)

func main() {
	a := smbios.NewDMITable()
	fmt.Println(a.Version())
	//
	fmt.Println(a.GetEntriesByType(17))

	//	for i, v := range smbios.StringKeyword {
	//		fmt.Printf("[%d] Keyword: %s, Type: %d, Offset: %d\n", i, v.Keyword, v.Type, v.Offset)
	//	}
	//
	for k, v := range a.Table {
		fmt.Printf("key: %s, { Keyword: %s, Type: %d, Offset: %d }\n", k, v.Keyword, v.Type, v.Offset)
	}

	fmt.Println(a.GetResultByKeyword("firmware-revision"))
	fmt.Println(a.GetResultByKeyword("bios-version"))
	fmt.Println(a.GetResultByKeyword("bios-vendor"))
	fmt.Println(a.GetResultByKeyword("bios-version"))
	fmt.Println(a.GetResultByKeyword("bios-release-date"))
	fmt.Println(a.GetResultByKeyword("bios-revision"))
	fmt.Println(a.GetResultByKeyword("firmware-revision")) /* 0x16 and 0x17 */
	fmt.Println(a.GetResultByKeyword("system-manufacturer"))
	fmt.Println(a.GetResultByKeyword("system-product-name"))
	fmt.Println(a.GetResultByKeyword("system-version"))
	fmt.Println(a.GetResultByKeyword("system-serial-number"))
	fmt.Println(a.GetResultByKeyword("system-uuid")) /* dmi_system_uuid() */
	fmt.Println(a.GetResultByKeyword("system-family"))
	fmt.Println(a.GetResultByKeyword("baseboard-manufacturer"))
	fmt.Println(a.GetResultByKeyword("baseboard-product-name"))
	fmt.Println(a.GetResultByKeyword("baseboard-version"))
	fmt.Println(a.GetResultByKeyword("baseboard-serial-number"))
	fmt.Println(a.GetResultByKeyword("baseboard-asset-tag"))
	fmt.Println(a.GetResultByKeyword("chassis-manufacturer"))
	fmt.Println(a.GetResultByKeyword("chassis-type")) /* dmi_chassis_type() */
	fmt.Println(a.GetResultByKeyword("chassis-version"))
	fmt.Println(a.GetResultByKeyword("chassis-serial-number"))
	fmt.Println(a.GetResultByKeyword("chassis-asset-tag"))
	fmt.Println(a.GetResultByKeyword("processor-family")) /* dmi_processor_family() */
	fmt.Println(a.GetResultByKeyword("processor-manufacturer"))
	fmt.Println(a.GetResultByKeyword("processor-version"))
	fmt.Println(a.GetResultByKeyword("processor-frequency")) /* dmi_processor_frequency() */
}
