
/*
   The SMBIOS specification defines the following DMI types:

    Type   Information
    ────────────────────────────────────────────
       0   BIOS
       1   System
       2   Baseboard
       3   Chassis
       4   Processor
       5   Memory Controller
       6   Memory Module
       7   Cache
       8   Port Connector
       9   System Slots
      10   On Board Devices
      11   OEM Strings
      12   System Configuration Options
      13   BIOS Language
      14   Group Associations
      15   System Event Log

      16   Physical Memory Array
      17   Memory Device
      18   32-bit Memory Error
      19   Memory Array Mapped Address
      20   Memory Device Mapped Address
      21   Built-in Pointing Device
      22   Portable Battery
      23   System Reset
      24   Hardware Security
      25   System Power Controls
      26   Voltage Probe
      27   Cooling Device
      28   Temperature Probe
      29   Electrical Current Probe
      30   Out-of-band Remote Access
      31   Boot Integrity Services
      32   System Boot
      33   64-bit Memory Error
      34   Management Device
      35   Management Device Component
      36   Management Device Threshold Data
      37   Memory Channel
      38   IPMI Device
      39   Power Supply
      40   Additional Information
      41   Onboard Devices Extended Information
      42   Management Controller Host Interface

-s, --string KEYWORD
Only  display  the  value of the DMI string identified by KEYWORD.  KEYWORD must be a keyword from the following list: bios-vendor, bios-ver‐
sion, bios-release-date, system-manufacturer, system-product-name, system-version, system-serial-number, system-uuid, baseboard-manufacturer,
baseboard-product-name, baseboard-version, baseboard-serial-number, baseboard-asset-tag, chassis-manufacturer, chassis-type, chassis-version,
chassis-serial-number, chassis-asset-tag, processor-family, processor-manufacturer,  processor-version,  processor-frequency.   Each  keyword
corresponds to a given DMI type and a given offset within this entry type.  Not all strings may be meaningful or even defined on all systems.
Some keywords may return more than one result on some systems (e.g.  processor-version on a multi-processor system).  If KEYWORD is not  pro‐
vided or not valid, a list of all valid keywords is printed and dmidecode exits with an error.  This option cannot be used more than once.

Note:  on Linux, most of these strings can alternatively be read directly from sysfs, typically from files under /sys/devices/virtual/dmi/id.
Most of these files are even readable by regular users.



## go test run specific test function
```
go test . -v -run '^TestDecoder$'
```

```

## dmidecode data type
```
/*
 * Per SMBIOS v2.8.0 and later, all structures assume a little-endian
 * ordering convention.
 */
#if defined(ALIGNMENT_WORKAROUND) || defined(BIGENDIAN)
#define WORD(x) (u16)((x)[0] + ((x)[1] << 8))
#define DWORD(x) (u32)((x)[0] + ((x)[1] << 8) + ((x)[2] << 16) + ((x)[3] << 24))
#define QWORD(x) (U64(DWORD(x), DWORD(x + 4)))
#else /* ALIGNMENT_WORKAROUND || BIGENDIAN */
#define WORD(x) (u16)(*(const u16 *)(x))
#define DWORD(x) (u32)(*(const u32 *)(x))
#define QWORD(x) (*(const u64 *)(x))
#endif /* ALIGNMENT_WORKAROUND || BIGENDIAN */

#endif
```
