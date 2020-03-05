package main

import (
	"fmt"

	"github.com/fenglyu/smbios"
)

func main() {

	dmit := smbios.NewDMITable()
	fmt.Println(dmit.Version())

	for k, _ := range dmit.Table {
		fmt.Printf("[%s] %s\n", k, dmit.Query(k))
	}
}
