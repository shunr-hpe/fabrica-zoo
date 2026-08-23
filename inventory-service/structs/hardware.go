package smd_structs

import (
	"context"
	"encoding/json"
)

// HardwareSpec defines the desired state of Hardware
type HardwareSpec struct {
	ID      string `json:"ID" yaml:"ID"`
	Type    string `json:"Type" yaml:"Type"`
	Ordinal int    `json:"Ordinal" yaml:"Ordinal"`
	Status  string `json:"Status" yaml:"Status"`

	// This is used as a descriminator to determine the type of *Info
	// struct that will be included below.
	HWInventoryByLocationType string `json:"HWInventoryByLocationType" yaml:"HWInventoryByLocationType"`

	// One of:var ErrHMSXnameInvalid = errors.New("got HMSTypeInvalid instead of valid type")
	//    HMSType                  Underlying RF Type          How named in json object
	HMSCabinetLocationInfo       *ChassisLocationInfoRF   `json:"CabinetLocationInfo,omitempty" yaml:"CabinetLocationInfo,omitempty"`
	HMSChassisLocationInfo       *ChassisLocationInfoRF   `json:"ChassisLocationInfo,omitempty" yaml:"ChassisLocationInfo,omitempty"` // Mountain chassis
	HMSComputeModuleLocationInfo *ChassisLocationInfoRF   `json:"ComputeModuleLocationInfo,omitempty" yaml:"ComputeModuleLocationInfo,omitempty"`
	HMSRouterModuleLocationInfo  *ChassisLocationInfoRF   `json:"RouterModuleLocationInfo,omitempty" yaml:"RouterModuleLocationInfo,omitempty"`
	HMSNodeEnclosureLocationInfo *ChassisLocationInfoRF   `json:"NodeEnclosureLocationInfo,omitempty" yaml:"NodeEnclosureLocationInfo,omitempty"`
	HMSHSNBoardLocationInfo      *ChassisLocationInfoRF   `json:"HSNBoardLocationInfo,omitempty" yaml:"HSNBoardLocationInfo,omitempty"`
	HMSMgmtSwitchLocationInfo    *ChassisLocationInfoRF   `json:"MgmtSwitchLocationInfo,omitempty" yaml:"MgmtSwitchLocationInfo,omitempty"`
	HMSMgmtHLSwitchLocationInfo  *ChassisLocationInfoRF   `json:"MgmtHLSwitchLocationInfo,omitempty" yaml:"MgmtHLSwitchLocationInfo,omitempty"`
	HMSCDUMgmtSwitchLocationInfo *ChassisLocationInfoRF   `json:"CDUMgmtSwitchLocationInfo,omitempty" yaml:"CDUMgmtSwitchLocationInfo,omitempty"`
	HMSNodeLocationInfo          *SystemLocationInfoRF    `json:"NodeLocationInfo,omitempty" yaml:"NodeLocationInfo,omitempty"`
	HMSProcessorLocationInfo     *ProcessorLocationInfoRF `json:"ProcessorLocationInfo,omitempty" yaml:"ProcessorLocationInfo,omitempty"`
	HMSNodeAccelLocationInfo     *ProcessorLocationInfoRF `json:"NodeAccelLocationInfo,omitempty" yaml:"NodeAccelLocationInfo,omitempty"`
	HMSMemoryLocationInfo        *MemoryLocationInfoRF    `json:"MemoryLocationInfo,omitempty" yaml:"MemoryLocationInfo,omitempty"`
	HMSDriveLocationInfo         *DriveLocationInfoRF     `json:"DriveLocationInfo,omitempty" yaml:"DriveLocationInfo,omitempty"`
	HMSHSNNICLocationInfo        *NALocationInfoRF        `json:"NodeHsnNicLocationInfo,omitempty" yaml:"NodeHsnNicLocationInfo,omitempty"`

	HMSPDULocationInfo                      *PowerDistributionLocationInfo `json:"PDULocationInfo,omitempty" yaml:"PDULocationInfo,omitempty"`
	HMSOutletLocationInfo                   *OutletLocationInfo            `json:"OutletLocationInfo,omitempty" yaml:"OutletLocationInfo,omitempty"`
	HMSCMMRectifierLocationInfo             *PowerSupplyLocationInfoRF     `json:"CMMRectifierLocationInfo,omitempty" yaml:"CMMRectifierLocationInfo,omitempty"`
	HMSNodeEnclosurePowerSupplyLocationInfo *PowerSupplyLocationInfoRF     `json:"NodeEnclosurePowerSupplyLocationInfo,omitempty" yaml:"NodeEnclosurePowerSupplyLocationInfo,omitempty"`
	HMSNodeBMCLocationInfo                  *ManagerLocationInfoRF         `json:"NodeBMCLocationInfo,omitempty" yaml:"NodeBMCLocationInfo,omitempty"`
	HMSRouterBMCLocationInfo                *ManagerLocationInfoRF         `json:"RouterBMCLocationInfo,omitempty" yaml:"RouterBMCLocationInfo,omitempty"`
	HMSNodeAccelRiserLocationInfo           *NodeAccelRiserLocationInfoRF  `json:"NodeAccelRiserLocationInfo,omitempty" yaml:"NodeAccelRiserLocationInfo,omitempty"`
	// TODO: Remaining types in hmsTypeArrays

	// If status != empty, up to one of following, matching above *Info.
	PopulatedFRU *HWInvByFRU `json:"PopulatedFRU,omitempty" yaml:"PopulatedFRU,omitempty"`

	// These are for nested references for subcomponents.
	hmsTypeArrays
}

type HWInvByFRU struct {
	FRUID   string `json:"FRUID" yaml:"FRUID"`
	Type    string `json:"Type" yaml:"Type"`
	Subtype string `json:"Subtype" yaml:"Subtype"`

	// This is used as a descriminator to specify the type of *Info
	// struct that will be included below.
	HWInventoryByFRUType string `json:"HWInventoryByFRUType" yaml:"HWInventoryByFRUType"`

	// One of (based on HWFRUInfoType):
	//   HMSType             Underlying RF Type      How named in json object
	HMSCabinetFRUInfo       *ChassisFRUInfoRF   `json:"CabinetFRUInfo,omitempty" yaml:"CabinetFRUInfo,omitempty"`
	HMSChassisFRUInfo       *ChassisFRUInfoRF   `json:"ChassisFRUInfo,omitempty" yaml:"ChassisFRUInfo,omitempty"` // Mountain chassis
	HMSComputeModuleFRUInfo *ChassisFRUInfoRF   `json:"ComputeModuleFRUInfo,omitempty" yaml:"ComputeModuleFRUInfo,omitempty"`
	HMSRouterModuleFRUInfo  *ChassisFRUInfoRF   `json:"RouterModuleFRUInfo,omitempty" yaml:"RouterModuleFRUInfo,omitempty"`
	HMSNodeEnclosureFRUInfo *ChassisFRUInfoRF   `json:"NodeEnclosureFRUInfo,omitempty" yaml:"NodeEnclosureFRUInfo,omitempty"`
	HMSHSNBoardFRUInfo      *ChassisFRUInfoRF   `json:"HSNBoardFRUInfo,omitempty" yaml:"HSNBoardFRUInfo,omitempty"`
	HMSMgmtSwitchFRUInfo    *ChassisFRUInfoRF   `json:"MgmtSwitchFRUInfo,omitempty" yaml:"MgmtSwitchFRUInfo,omitempty"`
	HMSMgmtHLSwitchFRUInfo  *ChassisFRUInfoRF   `json:"MgmtHLSwitchFRUInfo,omitempty" yaml:"MgmtHLSwitchFRUInfo,omitempty"`
	HMSCDUMgmtSwitchFRUInfo *ChassisFRUInfoRF   `json:"CDUMgmtSwitchFRUInfo,omitempty" yaml:"CDUMgmtSwitchFRUInfo,omitempty"`
	HMSNodeFRUInfo          *SystemFRUInfoRF    `json:"NodeFRUInfo,omitempty" yaml:"NodeFRUInfo,omitempty"`
	HMSProcessorFRUInfo     *ProcessorFRUInfoRF `json:"ProcessorFRUInfo,omitempty" yaml:"ProcessorFRUInfo,omitempty"`
	HMSNodeAccelFRUInfo     *ProcessorFRUInfoRF `json:"NodeAccelFRUInfo,omitempty" yaml:"NodeAccelFRUInfo,omitempty"`
	HMSMemoryFRUInfo        *MemoryFRUInfoRF    `json:"MemoryFRUInfo,omitempty" yaml:"MemoryFRUInfo,omitempty"`
	HMSDriveFRUInfo         *DriveFRUInfoRF     `json:"DriveFRUInfo,omitempty" yaml:"DriveFRUInfo,omitempty"`
	HMSHSNNICFRUInfo        *NAFRUInfoRF        `json:"NodeHsnNicFRUInfo,omitempty" yaml:"NodeHsnNicFRUInfo,omitempty"`

	HMSPDUFRUInfo                      *PowerDistributionFRUInfo `json:"PDUFRUInfo,omitempty" yaml:"PDUFRUInfo,omitempty"`
	HMSOutletFRUInfo                   *OutletFRUInfo            `json:"OutletFRUInfo,omitempty" yaml:"OutletFRUInfo,omitempty"`
	HMSCMMRectifierFRUInfo             *PowerSupplyFRUInfoRF     `json:"CMMRectifierFRUInfo,omitempty" yaml:"CMMRectifierFRUInfo,omitempty"`
	HMSNodeEnclosurePowerSupplyFRUInfo *PowerSupplyFRUInfoRF     `json:"NodeEnclosurePowerSupplyFRUInfo,omitempty" yaml:"NodeEnclosurePowerSupplyFRUInfo,omitempty"`
	HMSNodeBMCFRUInfo                  *ManagerFRUInfoRF         `json:"NodeBMCFRUInfo,omitempty" yaml:"NodeBMCFRUInfo,omitempty"`
	HMSRouterBMCFRUInfo                *ManagerFRUInfoRF         `json:"RouterBMCFRUInfo,omitempty" yaml:"RouterBMCFRUInfo,omitempty"`
	HMSNodeAccelRiserFRUInfo           *NodeAccelRiserFRUInfoRF  `json:"NodeAccelRiserFRUInfo,omitempty" yaml:"NodeAccelRiserFRUInfo,omitempty"`

	// TODO: Remaining types in hmsTypeArray
}

type ChassisLocationInfoRF struct {
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
	Hostname    string `json:"HostName" yaml:"HostName"`
}

// Redfish ProcessorSummary struct - Sub-struct of ComputerSystem
type ComputerSystemProcessorSummary struct {
	Count json.Number `json:"Count" yaml:"Count"`
	Model string      `json:"Model" yaml:"Model"`
	//Status StatusRF    `json:"Status" yaml:"Status"`
}

// Redfish MemorySummary struct - Sub-struct of ComputerSystem
type ComputerSystemMemorySummary struct {
	TotalSystemMemoryGiB json.Number `json:"TotalSystemMemoryGiB" yaml:"TotalSystemMemoryGiB"`
	//Status               rf.StatusRF    `json:"Status" yaml:"Status"`
}

// Location-specific Redfish properties to be stored in hardware inventory
// These are only relevant to the currently installed location of the FRU
// TODO: How to version these (as HMS structures).
type SystemLocationInfoRF struct {
	// Redfish pass-through from Redfish ComputerSystem
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
	Hostname    string `json:"HostName" yaml:"HostName"`

	ProcessorSummary ComputerSystemProcessorSummary `json:"ProcessorSummary" yaml:"ProcessorSummary"`

	MemorySummary ComputerSystemMemorySummary `json:"MemorySummary" yaml:"MemorySummary"`
}

// Location-specific Redfish properties to be stored in hardware inventory
// These are only relevant to the currently installed location of the FRU
// TODO: How to version these (as HMS structures).
type ProcessorLocationInfoRF struct {
	// Redfish pass-through from rf.Processor
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
	Socket      string `json:"Socket" yaml:"Socket"`
}

// Location-specific Redfish properties to be stored in hardware inventory
// These are only relevant to the currently installed location of the FRU
// TODO: How to version these (as HMS structures)
type MemoryLocationInfoRF struct {
	// Redfish pass-through from rf.Memory
	Id             string           `json:"Id" yaml:"Id"`
	Name           string           `json:"Name" yaml:"Name"`
	Description    string           `json:"Description" yaml:"Description"`
	MemoryLocation MemoryLocationRF `json:"MemoryLocation" yaml:"MemoryLocation"`
}

type MemoryLocationRF struct {
	Socket           json.Number `json:"Socket" yaml:"Socket"`
	MemoryController json.Number `json:"MemoryController" yaml:"MemoryController"`
	Channel          json.Number `json:"Channel" yaml:"Channel"`
	Slot             json.Number `json:"Slot" yaml:"Slot"`
}

type ManagerLocationInfoRF struct {
	DateTime            string `json:"DateTime" yaml:"DateTime"`
	DateTimeLocalOffset string `json:"DateTimeLocalOffset" yaml:"DateTimeLocalOffset"`
	Description         string `json:"Description" yaml:"Description"`
	FirmwareVersion     string `json:"FirmwareVersion" yaml:"FirmwareVersion"`
	Id                  string `json:"Id" yaml:"Id"`
	Name                string `json:"Name" yaml:"Name"`
}

// Location-specific Redfish properties to be stored in hardware inventory
// These are only relevant to the currently installed location of the FRU
// TODO: How to version these (as HMS structures).
type NodeAccelRiserLocationInfoRF struct {
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
}

// Location-specific Redfish properties to be stored in hardware inventory
// These are only relevant to the currently installed location of the FRU
type PowerSupplyLocationInfoRF struct {
	Name            string `json:"Name" yaml:"Name"`
	FirmwareVersion string `json:"FirmwareVersion" yaml:"FirmwareVersion"`
}

// Location-specific Redfish properties to be stored in hardware inventory
// These are only relevant to the currently installed location of the FRU
// TODO: How to version these (as HMS structures).
type DriveLocationInfoRF struct {
	// Redfish pass-through from rf.Drive
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
}

// Location-specific Redfish properties to be stored in hardware inventory
// These are only relevant to the currently installed location of the FRU
type NALocationInfoRF struct {
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
}

// Redfish fields from the PowerDistributionFRUInfo schema that go into
// HWInventoryByLocation.  We capture them as an embedded struct within the
// full schema during inventory discovery.
type PowerDistributionLocationInfo struct {
	Id          string    `json:"Id" yaml:"Id"`
	Description string    `json:"Description" yaml:"Description"`
	Name        string    `json:"Name" yaml:"Name"`
	UUID        string    `json:"UUID" yaml:"UUID"`
	Location    *Location `json:"Location,omitempty" yaml:"Location,omitempty"`
}

// Location
//
// Resource type.  Appears under Chassis, PowerDistribution, etc.
type Location struct {
	ContactInfo   *ContactInfo   `json:"ContactInfo,omitempty" yaml:"ContactInfo,omitempty"`
	Latitude      json.Number    `json:"Latitude,omitempty" yaml:"Latitude,omitempty"`
	Longitude     json.Number    `json:"Longitude,omitempty" yaml:"Longitude,omitempty"`
	PartLocation  *PartLocation  `json:"PartLocation,omitempty" yaml:"PartLocation,omitempty"`
	Placement     *Placement     `json:"Placement,omitempty" yaml:"Placement,omitempty"`
	PostalAddress *PostalAddress `json:"PostalAddress,omitempty" yaml:"PostalAddress,omitempty"`
}

// Within Location - ContactInfo
type ContactInfo struct {
	ContactName  string `json:"ContactName" yaml:"ContactName"`
	EmailAddress string `json:"EmailAddress" yaml:"EmailAddress"`
	PhoneNumber  string `json:"PhoneNumber,omitempty" yaml:"PhoneNumber,omitempty"`
}

// Within Location - PartLocation
type PartLocation struct {
	LocationOrdinalValue json.Number `json:"LocationOrdinalValue,omitempty" yaml:"LocationOrdinalValue,omitempty"`
	LocationType         string      `json:"LocationType" yaml:"LocationType"` //enum
	Orientation          string      `json:"Orientation" yaml:"Orientation"`   //enum
	Reference            string      `json:"Reference" yaml:"Reference"`       //enum
	ServiceLabel         string      `json:"ServiceLabel" yaml:"ServiceLabel"`
}

// Within Location - PostalAddress
type PostalAddress struct {
	Country    string `json:"Country" yaml:"Country"`
	Territory  string `json:"Territory" yaml:"Territory"`
	City       string `json:"City" yaml:"City"`
	Street     string `json:"Street" yaml:"Street"`
	Name       string `json:"Name" yaml:"Name"`
	PostalCode string `json:"PostalCode" yaml:"PostalCode"`
	Building   string `json:"Building" yaml:"Building"`
	Floor      string `json:"Floor" yaml:"Floor"`
	Room       string `json:"Room" yaml:"Room"`
}

type Placement struct {
	AdditionalInfo  string      `json:"AdditionalInfo,omitempty" yaml:"AdditionalInfo,omitempty"`
	Rack            string      `json:"Rack,omitempty" yaml:"Rack,omitempty"`
	RackOffset      json.Number `json:"RackOffset,omitempty" yaml:"RackOffset,omitempty"`
	RackOffsetUnits string      `json:"RackOffsetUnits,omitempty" yaml:"RackOffsetUnits,omitempty"`
	Row             string      `json:"Row,omitempty" yaml:"Row,omitempty"`
}

// Outlets do not have individual FRUs, PDUs do, but their properties are
// potentially important.  This is location-dependent data for HwInventory
type OutletLocationInfo struct {
	Id          string `json:"Id" yaml:"Id"`
	Description string `json:"Description" yaml:"Description"`
	Name        string `json:"Name" yaml:"Name"`
}

// Validate implements custom validation logic for Hardware
func (r *Hardware) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}

// Durable Redfish properties to be stored in hardware inventory as
// a specific FRU, which is then link with it's current location
// i.e. an x-name.  These properties should follow the hardware and
// allow it to be tracked even when it is removed from the system.
// TODO: How to version these (as HMS structures)
type ChassisFRUInfoRF struct {
	AssetTag     string `json:"AssetTag" yaml:"AssetTag"`
	ChassisType  string `json:"ChassisType" yaml:"ChassisType"`
	Model        string `json:"Model" yaml:"Model"`
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
	SKU          string `json:"SKU" yaml:"SKU"`
}

// Durable Redfish properties to be stored in hardware inventory as
// a specific FRU, which is then (typically) associated with a location
// i.e. an x-name in HMS terms and the ProcessorLocationInfo fields
// in Redfish terms on the controller.  These properties should
// follow the hardware and allow it to be tracked even when it is removed
// from the system.
// TODO: How to version these (as HMS structures).
type SystemFRUInfoRF struct {
	// Redfish pass-through from Redfish ComputerSystem
	AssetTag     string `json:"AssetTag" yaml:"AssetTag"`
	BiosVersion  string `json:"BiosVersion" yaml:"BiosVersion"`
	Model        string `json:"Model" yaml:"Model"`
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
	SKU          string `json:"SKU" yaml:"SKU"`
	SystemType   string `json:"SystemType" yaml:"SystemType"`
	UUID         string `json:"UUID" yaml:"UUID"`
}

// Durable Redfish properties to be stored in hardware inventory as
// a specific FRU, which is then link with it's current location
// i.e. an x-name.  These properties should follow the hardware and
// allow it to be tracked even when it is removed from the system.
// TODO: How to version these (as HMS structures)
type ProcessorFRUInfoRF struct {
	// Redfish pass-through from rf.Processor
	InstructionSet        string        `json:"InstructionSet" yaml:"InstructionSet"`
	Manufacturer          string        `json:"Manufacturer" yaml:"Manufacturer"`
	MaxSpeedMHz           json.Number   `json:"MaxSpeedMHz" yaml:"MaxSpeedMHz"`
	Model                 string        `json:"Model" yaml:"Model"`
	SerialNumber          string        `json:"SerialNumber" yaml:"SerialNumber"`
	PartNumber            string        `json:"PartNumber" yaml:"PartNumber"`
	ProcessorArchitecture string        `json:"ProcessorArchitecture" yaml:"ProcessorArchitecture"`
	ProcessorId           ProcessorIdRF `json:"ProcessorId" yaml:"ProcessorId"`
	ProcessorType         string        `json:"ProcessorType" yaml:"ProcessorType"`
	TotalCores            json.Number   `json:"TotalCores" yaml:"TotalCores"`
	TotalThreads          json.Number   `json:"TotalThreads" yaml:"TotalThreads"`
	Oem                   *ProcessorOEM `json:"Oem" yaml:"Oem"`
}

type ProcessorOEM struct {
	GBTProcessorOemProperty *GBTProcessorOem `json:"GBTProcessorOemProperty,omitempty" yaml:"GBTProcessorOemProperty,omitempty"`
}

type GBTProcessorOem struct {
	ProcessorSerialNumber string `json:"Processor Serial Number,omitempty" yaml:"Processor Serial Number,omitempty"`
}

type ProcessorIdRF struct {
	EffectiveFamily         string `json:"EffectiveFamily" yaml:"EffectiveFamily"`
	EffectiveModel          string `json:"EffectiveModel" yaml:"EffectiveModel"`
	IdentificationRegisters string `json:"IdentificationRegisters" yaml:"IdentificationRegisters"`
	MicrocodeInfo           string `json:"MicrocodeInfo" yaml:"MicrocodeInfo"`
	Step                    string `json:"Step" yaml:"Step"`
	VendorID                string `json:"VendorID" yaml:"VendorID"`
}

// Durable Redfish properties to be stored in hardware inventory as
// a specific FRU, which is then link with it's current location
// i.e. an x-name.  These properties should follow the hardware and
// allow it to be tracked even when it is removed from the system.
// TODO: How to version these (as HMS structures)
type MemoryFRUInfoRF struct {
	// Redfish pass-through from rf.Memory
	BaseModuleType    string      `json:"BaseModuleType,omitempty" yaml:"BaseModuleType,omitempty"`
	BusWidthBits      json.Number `json:"BusWidthBits,omitempty" yaml:"BusWidthBits,omitempty"`
	CapacityMiB       json.Number `json:"CapacityMiB" yaml:"CapacityMiB"`
	DataWidthBits     json.Number `json:"DataWidthBits,omitempty" yaml:"DataWidthBits,omitempty"`
	ErrorCorrection   string      `json:"ErrorCorrection,omitempty" yaml:"ErrorCorrection,omitempty"`
	Manufacturer      string      `json:"Manufacturer,omitempty" yaml:"Manufacturer,omitempty"`
	MemoryType        string      `json:"MemoryType,omitempty" yaml:"MemoryType,omitempty"`
	MemoryDeviceType  string      `json:"MemoryDeviceType,omitempty" yaml:"MemoryDeviceType,omitempty"`
	OperatingSpeedMhz json.Number `json:"OperatingSpeedMhz" yaml:"OperatingSpeedMhz"`
	PartNumber        string      `json:"PartNumber,omitempty" yaml:"PartNumber,omitempty"`
	RankCount         json.Number `json:"RankCount,omitempty" yaml:"RankCount,omitempty"`
	SerialNumber      string      `json:"SerialNumber" yaml:"SerialNumber"`
}

// This is an embedded structure for HW inventory.  There should be one
// array for every hms type tracked in the inventory.  This structure
// is also reused to allow individual HWInvByLoc structures to represent
// child components for nested inventory structures.
type hmsTypeArrays struct {
	Nodes          *[]*HardwareSpec `json:"Nodes,omitempty" yaml:"Nodes,omitempty"`
	Cabinets       *[]*HardwareSpec `json:"Cabinets,omitempty" yaml:"Cabinets,omitempty"`
	Chassis        *[]*HardwareSpec `json:"Chassis,omitempty" yaml:"Chassis,omitempty"`
	ComputeModules *[]*HardwareSpec `json:"ComputeModules,omitempty" yaml:"ComputeModules,omitempty"`
	RouterModules  *[]*HardwareSpec `json:"RouterModules,omitempty" yaml:"RouterModules,omitempty"`
	NodeEnclosures *[]*HardwareSpec `json:"NodeEnclosures,omitempty" yaml:"NodeEnclosures,omitempty"`
	HSNBoards      *[]*HardwareSpec `json:"HSNBoards,omitempty" yaml:"HSNBoards,omitempty"`

	Processors *[]*HardwareSpec `json:"Processors,omitempty" yaml:"Processors,omitempty"`
	Memory     *[]*HardwareSpec `json:"Memory,omitempty" yaml:"Memory,omitempty"`
	Drives     *[]*HardwareSpec `json:"Drives,omitempty" yaml:"Drives,omitempty"`

	CabinetPDUs                *[]*HardwareSpec `json:"CabinetPDUs,omitempty" yaml:"CabinetPDUs,omitempty"`
	CabinetPDUOutlets          *[]*HardwareSpec `json:"CabinetPDUPowerConnectors,omitempty" yaml:"CabinetPDUPowerConnectors,omitempty"`
	CMMRectifiers              *[]*HardwareSpec `json:"CMMRectifiers,omitempty" yaml:"CMMRectifiers,omitempty"`
	NodeAccels                 *[]*HardwareSpec `json:"NodeAccels,omitempty" yaml:"NodeAccels,omitempty"`
	NodeAccelRisers            *[]*HardwareSpec `json:"NodeAccelRisers,omitempty" yaml:"NodeAccelRisers,omitempty"`
	NodeEnclosurePowerSupplies *[]*HardwareSpec `json:"NodeEnclosurePowerSupplies,omitempty" yaml:"NodeEnclosurePowerSupplies,omitempty"`
	NodeHsnNICs                *[]*HardwareSpec `json:"NodeHsnNics,omitempty" yaml:"NodeHsnNics,omitempty"`

	// These don't have hardware inventory location/FRU info yet,
	// either because they aren't known yet or because they are manager
	// types.  Each manager (e.g. BMC) should have some kind of physical
	// enclosure, and for the purposes of HW inventory we might not need
	// both (but probably will).
	CECs           *[]*HardwareSpec `json:"CECs,omitempty" yaml:"CECs,omitempty"`
	CDUs           *[]*HardwareSpec `json:"CDUs,omitempty" yaml:"CDUs,omitempty"`
	CabinetCDUs    *[]*HardwareSpec `json:"CabinetCDUs,omitempty" yaml:"CabinetCDUs,omitempty"`
	CMMFpgas       *[]*HardwareSpec `json:"CMMFpgas,omitempty" yaml:"CMMFpgas,omitempty"`
	NodeFpgas      *[]*HardwareSpec `json:"NodeFpgas,omitempty" yaml:"NodeFpgas,omitempty"`
	RouterFpgas    *[]*HardwareSpec `json:"RouterFpgas,omitempty" yaml:"RouterFpgas,omitempty"`
	RouterTORFpgas *[]*HardwareSpec `json:"RouterTORFpgas,omitempty" yaml:"RouterTORFpgas,omitempty"`
	HSNAsics       *[]*HardwareSpec `json:"HSNAsics,omitempty" yaml:"HSNAsics,omitempty"`

	CabinetBMCs           *[]*HardwareSpec `json:"CabinetBMCs,omitempty" yaml:"CabinetBMCs,omitempty"`
	CabinetPDUControllers *[]*HardwareSpec `json:"CabinetPDUControllers,omitempty" yaml:"CabinetPDUControllers,omitempty"`
	ChassisBMCs           *[]*HardwareSpec `json:"ChassisBMCs,omitempty" yaml:"ChassisBMCs,omitempty"`
	NodeBMCs              *[]*HardwareSpec `json:"NodeBMCs,omitempty" yaml:"NodeBMCs,omitempty"`
	RouterBMCs            *[]*HardwareSpec `json:"RouterBMCs,omitempty" yaml:"RouterBMCs,omitempty"`

	CabinetPDUNics      *[]*HardwareSpec `json:"CabinetPDUNics,omitempty" yaml:"CabinetPDUNics,omitempty"`
	NodePowerConnectors *[]*HardwareSpec `json:"NodePowerConnectors,omitempty" yaml:"NodePowerConnectors,omitempty"`
	NodeBMCNics         *[]*HardwareSpec `json:"NodeBMCNics,omitempty" yaml:"NodeBMCNics,omitempty"`
	NodeNICs            *[]*HardwareSpec `json:"NodeNICs,omitempty" yaml:"NodeNICs,omitempty"`
	RouterBMCNics       *[]*HardwareSpec `json:"RouterBMCNics,omitempty" yaml:"RouterBMCNics,omitempty"`

	MgmtSwitches    *[]*HardwareSpec `json:"MgmtSwitches,omitempty" yaml:"MgmtSwitches,omitempty"`
	MgmtHLSwitches  *[]*HardwareSpec `json:"MgmtHLSwitches,omitempty" yaml:"MgmtHLSwitches,omitempty"`
	CDUMgmtSwitches *[]*HardwareSpec `json:"CDUMgmtSwitches,omitempty" yaml:"CDUMgmtSwitches,omitempty"`

	// Also not implemented yet.  Not clear if these will have any interesting
	// info, so they may never be,
	SMSBoxes             *[]*HardwareSpec `json:"SMSBoxes,omitempty" yaml:"SMSBoxes,omitempty"`
	HSNLinks             *[]*HardwareSpec `json:"HSNLinks,omitempty" yaml:"HSNLinks,omitempty"`
	HSNConnectors        *[]*HardwareSpec `json:"HSNConnectors,omitempty" yaml:"HSNConnectors,omitempty"`
	HSNConnectorPorts    *[]*HardwareSpec `json:"HSNConnectorPorts,omitempty" yaml:"HSNConnectorPorts,omitempty"`
	MgmtSwitchConnectors *[]*HardwareSpec `json:"MgmtSwitchConnectors,omitempty" yaml:"MgmtSwitchConnectors,omitempty"`
}

// Durable Redfish properties to be stored in hardware inventory as
// a specific FRU, which is then link with it's current location
// i.e. an x-name.  These properties should follow the hardware and
// allow it to be tracked even when it is removed from the system.
// TODO: How to version these (as HMS structures)
type DriveFRUInfoRF struct {
	// Redfish pass-through from rf.Drive

	//Manufacture Info
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	Model        string `json:"Model" yaml:"Model"`
	SKU          string `json:"SKU" yaml:"SKU"`

	//Capabilities Info
	CapacityBytes    json.Number `json:"CapacityBytes" yaml:"CapacityBytes"`
	Protocol         string      `json:"Protocol" yaml:"Protocol"`
	MediaType        string      `json:"MediaType" yaml:"MediaType"`
	RotationSpeedRPM json.Number `json:"RotationSpeedRPM" yaml:"RotationSpeedRPM"`
	BlockSizeBytes   json.Number `json:"BlockSizeBytes" yaml:"BlockSizeBytes"`
	CapableSpeedGbs  json.Number `json:"CapableSpeedGbs" yaml:"CapableSpeedGbs"`

	//Status Info
	FailurePredicted              bool        `json:"FailurePredicted" yaml:"FailurePredicted"`
	EncryptionAbility             string      `json:"EncryptionAbility" yaml:"EncryptionAbility"`
	EncryptionStatus              string      `json:"EncryptionStatus" yaml:"EncryptionStatus"`
	NegotiatedSpeedGbs            json.Number `json:"NegotiatedSpeedGbs" yaml:"NegotiatedSpeedGbs"`
	PredictedMediaLifeLeftPercent json.Number `json:"PredictedMediaLifeLeftPercent" yaml:"PredictedMediaLifeLeftPercent"`
}

// Durable Redfish properties to be stored in hardware inventory as
// a specific FRU, which is then link with it's current location
// i.e. an x-name.  These properties should follow the hardware and
// allow it to be tracked even when it is removed from the system.
type NAFRUInfoRF struct {
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	Model        string `json:"Model" yaml:"Model"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	SKU          string `json:"SKU,omitempty" yaml:"SKU,omitempty"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
}

// Redfish fields from the PowerDistribution schema that go into
// HWInventoryByFRU.  We capture them as an embedded struct within the
// full schema during inventory discovery.
type PowerDistributionFRUInfo struct {
	AssetTag          string         `json:"AssetTag" yaml:"AssetTag"`
	DateOfManufacture string         `json:"DateOfManufacture,omitempty" yaml:"DateOfManufacture,omitempty"`
	EquipmentType     string         `json:"EquipmentType" yaml:"EquipmentType"`
	FirmwareVersion   string         `json:"FirmwareVersion" yaml:"FirmwareVersion"`
	HardwareRevision  string         `json:"HardwareRevision" yaml:"HardwareRevision"`
	Manufacturer      string         `json:"Manufacturer" yaml:"Manufacturer"`
	Model             string         `json:"Model" yaml:"Model"`
	PartNumber        string         `json:"PartNumber" yaml:"PartNumber"`
	SerialNumber      string         `json:"SerialNumber" yaml:"SerialNumber"`
	CircuitSummary    CircuitSummary `json:"CircuitSummary" yaml:"CircuitSummary"`
}

// CircuitSummary sub-struct of PowerDistribution
// These are all-readonly
type CircuitSummary struct {
	ControlledOutlets json.Number `json:"ControlledOutlets,omitempty" yaml:"ControlledOutlets,omitempty"`
	MonitoredBranches json.Number `json:"MonitoredBranches,omitempty" yaml:"MonitoredBranches,omitempty"`
	MonitoredOutlets  json.Number `json:"MonitoredOutlets,omitempty" yaml:"MonitoredOutlets,omitempty"`
	MonitoredPhases   json.Number `json:"MonitoredPhases,omitempty" yaml:"MonitoredPhases,omitempty"`
	TotalBranches     json.Number `json:"TotalBranches,omitempty" yaml:"TotalBranches,omitempty"`
	TotalOutlets      json.Number `json:"TotalOutlets,omitempty" yaml:"TotalOutlets,omitempty"`
	TotalPhases       json.Number `json:"TotalPhases,omitempty" yaml:"TotalPhases,omitempty"`
}

// Outlets do not have individual FRUs, PDUs do, but their properties are
// potentially important. This is FRU-dependent data for HwInventory
// Note: omits configurable parameters.
type OutletFRUInfo struct {
	NominalVoltage   string         `json:"NominalVoltage,omitempty" yaml:"NominalVoltage,omitempty"`
	OutletType       string         `json:"OutletType" yaml:"OutletType"` // Enum
	EnergySensor     *SensorExcerpt `json:"EnergySensor,omitempty" yaml:"EnergySensor,omitempty"`
	FrequencySensor  *SensorExcerpt `json:"FrequencySensor,omitempty" yaml:"FrequencySensor,omitempty"`
	PhaseWiringType  string         `json:"PhaseWiringType,omitempty" yaml:"PhaseWiringType,omitempty"` // Enum
	PowerEnabled     *bool          `json:"PowerEnabled,omitempty" yaml:"PowerEnabled,omitempty"`       // Can be powered?
	RatedCurrentAmps json.Number    `json:"RatedCurrentAmps,omitempty" yaml:"RatedCurrentAmps,omitempty"`
	VoltageType      string         `json:"VoltageType,omitempty" yaml:"VoltageType,omitempty"` // Enum
}

// SensorExcerpt -  Substruct of Outlet and other power-related objects
// This is the more general non-power version of SensorPowerExcerpt
type SensorExcerpt struct {
	DataSourceUri      string      `json:"DataSourceUri" yaml:"DataSourceUri"`
	Name               string      `json:"Name" yaml:"Name"`
	PeakReading        json.Number `json:"PeakReading,omitempty" yaml:"PeakReading,omitempty"`
	PhysicalContext    string      `json:"PhysicalContext,omitempty" yaml:"PhysicalContext,omitempty"`       //enum
	PhysicalSubContext string      `json:"PhysicalSubContext,omitempty" yaml:"PhysicalSubContext,omitempty"` //enum
	Reading            json.Number `json:"Reading,omitempty" yaml:"Reading,omitempty"`
	ReadingUnits       string      `json:"ReadingUnits,omitempty" yaml:"ReadingUnits,omitempty"`
	Status             StatusRF    `json:"Status,omitempty" yaml:"Status,omitempty"`
}

// Durable Redfish properties to be stored in hardware inventory as
// a specific FRU, which is then link with it's current location
// i.e. an x-name.  These properties should follow the hardware and
// allow it to be tracked even when it is removed from the system.
// TODO: How to version these (as HMS structures)
type PowerSupplyFRUInfoRF struct {
	//Manufacture Info
	Manufacturer       string      `json:"Manufacturer" yaml:"Manufacturer"`
	SerialNumber       string      `json:"SerialNumber" yaml:"SerialNumber"`
	Model              string      `json:"Model" yaml:"Model"`
	PartNumber         string      `json:"PartNumber" yaml:"PartNumber"`
	PowerCapacityWatts int         `json:"PowerCapacityWatts" yaml:"PowerCapacityWatts"`
	PowerInputWatts    int         `json:"PowerInputWatts" yaml:"PowerInputWatts"`
	PowerOutputWatts   interface{} `json:"PowerOutputWatts" yaml:"PowerOutputWatts"`
	PowerSupplyType    string      `json:"PowerSupplyType" yaml:"PowerSupplyType"`
}

type ManagerFRUInfoRF struct {
	ManagerType  string `json:"ManagerType" yaml:"ManagerType"`
	Model        string `json:"Model" yaml:"Model"`
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
}

// Durable Redfish properties to be stored in hardware inventory as
// a specific FRU, which is then link with it's current location
// i.e. an x-name.  These properties should follow the hardware and
// allow it to be tracked even when it is removed from the system.
type NodeAccelRiserFRUInfoRF struct {
	//Manufacturer Info
	PhysicalContext        string             `json:"PhysicalContext" yaml:"PhysicalContext"`
	Producer               string             `json:"Producer" yaml:"Producer"`
	SerialNumber           string             `json:"SerialNumber" yaml:"SerialNumber"`
	PartNumber             string             `json:"PartNumber" yaml:"PartNumber"`
	Model                  string             `json:"Model" yaml:"Model"`
	ProductionDate         string             `json:"ProductionDate" yaml:"ProductionDate"`
	Version                string             `json:"Version" yaml:"Version"`
	EngineeringChangeLevel string             `json:"EngineeringChangeLevel" yaml:"EngineeringChangeLevel"`
	OEM                    *NodeAccelRiserOEM `json:"Oem,omitempty" yaml:"Oem,omitempty"`
}

type NodeAccelRiserOEM struct {
	PCBSerialNumber string `json:"PCBSerialNumber" yaml:"PCBSerialNumber"`
}
