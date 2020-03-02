package smbios

import (
	"testing"
)

func TestGetResultByKeyword(t *testing.T) {
	tests := []struct{
		name: string
		keyword: string
		result: string
		ok: bool
	}{
		{
			name: "test bios-vendor",
			keyword: "bios-vendor",
			result: "Hewlett-Packard",
			ok: true,
		},
		{
			name: "test system-product-name",
			keyword: "system-product-name",
			result: "HP Z420 Workstation",
			ok: true,
		}
    }	
	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			dmi := &DMITable{}
			dmi.GetResultByKeyword()
		}
	}
}
