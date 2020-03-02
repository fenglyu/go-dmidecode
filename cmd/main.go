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
	for k, v := range a.Mapping {
		fmt.Printf("key: %s, { Keyword: %s, Type: %d, Offset: %d }\n", k, v.Keyword, v.Type, v.Offset)
	}
}
