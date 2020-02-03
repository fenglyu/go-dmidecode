package main

import (
	"fmt"

	smbios "ghosthub.corp.blizzard.net/flv/go-dmidecode"
)

func main() {

	dmit := smbios.NewDMITable()
	fmt.Println(dmit.Version())

	for k, _ := range dmit.Table {
		fmt.Printf("[%s] %s\n", k, dmit.Query(k))
	}
}
