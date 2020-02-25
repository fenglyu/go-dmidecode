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
	fmt.Println(a.GetEntriesByType(17))
}
