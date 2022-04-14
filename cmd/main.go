package main

import (
	"fmt"
	"log"

	smbios "github.com/fenglyu/go-dmidecode"
)

func main() {

	dmit, err := smbios.NewDMITable()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dmit.Version())

	for k, _ := range smbios.Table {
		fmt.Printf("[%s] %s\n", k, dmit.Query(k))
	}
}
