package smd_structs

type EthernetInterfaceSpec struct {
	// Description            string              `json:"Description" yaml:"Description"`

	// ID is the ethernet interface identifier and serves as the resource identifier.
	//
	// +fabrica:field:index
	// +fabrica:field:unique
	ID          string      `json:"ID" yaml:"ID"`
	MACAddr     string      `json:"MACAddress" yaml:"MACAddress"`
	LastUpdate  string      `json:"LastUpdate" yaml:"LastUpdate"`
	CompID      string      `json:"ComponentID" yaml:"ComponentID"`
	Type        string      `json:"Type" yaml:"Type"`
	IPAddresses []IPAddress `json:"IPAddresses" yaml:"IPAddresses"`
}

type IPAddress struct {
	IPAddress string `json:"IPAddress" yaml:"IPAddress"`
	Network   string `json:"Network,omitempty" yaml:"Network,omitempty"`
}
