package smbios

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func Run(name string, arg ...string) string {

	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, arg...)
	var out, errbuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errbuf

	if err := cmd.Run(); err != nil {
		log.Fatal(errbuf.String())
		log.Fatal(err)
	}

	return out.String()
}

func TestGetResultByKeyword(t *testing.T) {
	tests := []struct {
		name    string
		keyword string
		//result  string
		ok bool
	}{
		{
			name:    "bios-vendor",
			keyword: "bios-vendor",
			ok:      true,
		},
		{
			name:    "bios-version",
			keyword: "bios-version",
			ok:      true,
		},
		{
			name:    "bios-release-date",
			keyword: "bios-release-date",
			ok:      true,
		},
		{
			name:    "bios-revision",
			keyword: "bios-revision",
			ok:      false,
		},
		{
			name:    "firmware-revision",
			keyword: "firmware-revision",
			ok:      false,
		},
		{
			name:    "system-manufacturer",
			keyword: "system-manufacturer",
			ok:      true,
		},
		{
			name:    "system-product-name",
			keyword: "system-product-name",
			ok:      true,
		},
		{
			name:    "system-version",
			keyword: "system-version",
			ok:      true,
		},
		{
			name:    "system-serial-number",
			keyword: "system-serial-number",
			ok:      true,
		},
		{
			name:    "system-uuid",
			keyword: "system-uuid",
			ok:      true,
		},
		{
			name:    "system-family",
			keyword: "system-family",
			ok:      false,
		},
		{
			name:    "baseboard-manufacturer",
			keyword: "baseboard-manufacturer",
			ok:      true,
		},
		{
			name:    "baseboard-product-name",
			keyword: "baseboard-product-name",
			ok:      true,
		},
		{
			name:    "baseboard-version",
			keyword: "baseboard-version",
			ok:      true,
		},
		{
			name:    "baseboard-serial-number",
			keyword: "baseboard-serial-number",
			ok:      true,
		},
		{
			name:    "baseboard-asset-tag",
			keyword: "baseboard-asset-tag",
			ok:      true,
		},
		{
			name:    "chassis-manufacturer",
			keyword: "chassis-manufacturer",
			ok:      true,
		},
		{
			name:    "chassis-type",
			keyword: "chassis-type",
			ok:      true,
		},
		{
			name:    "chassis-version",
			keyword: "chassis-version",
			ok:      true,
		},
		{
			name:    "chassis-serial-number",
			keyword: "chassis-serial-number",
			ok:      true,
		},
		{
			name:    "chassis-asset-tag",
			keyword: "chassis-asset-tag",
			ok:      true,
		},
		{
			name:    "processor-family",
			keyword: "processor-family",
			ok:      true,
		},
		{
			name:    "processor-manufacturer",
			keyword: "processor-manufacturer",
			ok:      true,
		},
		{
			name:    "processor-version",
			keyword: "processor-version",
			ok:      true,
		},
		{
			name:    "processor-frequency",
			keyword: "processor-frequency",
			ok:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// skip test bios-revision firmware-revision system-family
			// as dmidecode on most platform doesn't support them yet
			if !tt.ok {
				return
			}
			dmi := NewDMITable()
			result := dmi.Query(tt.keyword)
			// dmidecode's output always ends up with an newline
			expected := Run("dmidecode", "-s", tt.keyword)
			if !strings.EqualFold(result, strings.TrimSpace(expected)) {
				log.Fatal("Expected: ", strings.TrimSpace(expected), "\n", " Result: ", result)
			}
		})
	}
}
