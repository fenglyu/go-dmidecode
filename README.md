
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

BIOS DMIType = iota
System
Baseboard
Chassis
Processor
Memory_Controller
Memory_Module
Cache
Port_Connector
System_Slots
On_Board_Devices
OEM_Strings
System_Configuration_Options
BIOS_Language
Group_Associations
System_Event_Log
Physical_Memory_Array
Memory_Device
32-bit_Memory_Error
Memory_Array_Mapped_Address
Memory_Device_Mapped_Address
Built-in_Pointing_Device
Portable_Battery
System_Reset
Hardware_Security
System_Power_Controls
Voltage_Probe
Cooling_Device
Temperature_Probe
Electrical_Current_Probe
Out-of-band_Remote_Access
Boot_Integrity_Services
System_Boot
64-bit_Memory_Error
Management_Device
Management_Device_Component
Management_Device_Threshold_Data
Memory_Channel
IPMI_Device
Power_Supply
Additional_Information
Onboard_Devices_Extended_Information
Management_Controller_Host_Interface

*/

	//	fmt.Println(ss)
	//	for _, s := range ss {
	/*
		// Only look at memory devices.
		if s.Header.Type != 17 {
			continue
		}

		// Formatted section contains a variety of data, but only parse the DIMM size.
		size := int(binary.LittleEndian.Uint16(s.Formatted[8:10]))
		// String 0 is the DIMM slot's identifier.
		name := s.Strings[0]

		// If 0, no DIMM present in this slot.
		if size == 0 {
			fmt.Printf("[% 3s] empty\n", name)
			continue
		}

		// An extended uint32 DIMM size field appears if 0x7fff is present in size.
		if size == 0x7fff {
			size = int(binary.LittleEndian.Uint32(s.Formatted[24:28]))
		}

		// Size units depend on MSB.  Little endian MSB for uint16 is in second byte.
		// 0 means megabytes, 1 means kilobytes.
		unit := "KB"
		if s.Formatted[9]&0x80 == 0 {
			unit = "MB"
		}

		fmt.Printf("[% 3s] DIMM: %d %s\n", name, size, unit)
		//		fmt.Println(s)
	*/
	// Assume we find BIOS and end of table types.
