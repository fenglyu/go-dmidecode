## Go native support *dmidecode -s keyword*

```
String keyword expected
  bios-vendor
  bios-version
  bios-release-date
  bios-revision
  firmware-revision
  system-manufacturer
  system-product-name
  system-version
  system-serial-number
  system-uuid
  system-family
  baseboard-manufacturer
  baseboard-product-name
  baseboard-version
  baseboard-serial-number
  baseboard-asset-tag
  chassis-manufacturer
  chassis-type
  chassis-version
  chassis-serial-number
  chassis-asset-tag
  processor-family
  processor-manufacturer
  processor-version
  processor-frequency

```

## Implemention
The underlying DMI decode/parse is based on [go-smbios](https://github.com/digitalocean/go-smbios), Some functions are simply a re-implemention of the C version [dmidecode](https://github.com/mirror/dmidecode)
SMBIOS Documention reference [DSP0134_3.1.1.pdf](https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.1.1.pdf)
