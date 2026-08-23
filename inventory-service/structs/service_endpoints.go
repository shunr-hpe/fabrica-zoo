package smd_structs

import "encoding/json"

type ServiceEndpointSpec struct {
	// Embedded struct
	ServiceDescription

	// These are read-only, derived from associated RfEndpointId in
	// rf.ServiceDescription
	RfEndpointFQDN string `json:"RedfishEndpointFQDN" yaml:"RedfishEndpointFQDN"`
	URL            string `json:"RedfishURL" yaml:"RedfishURL"`

	// These are all stored in the same JSON blob, only one of these
	// one of these will be set, based on the value of RedfishType
	ServiceInfo json.RawMessage `json:"ServiceInfo,omitempty" yaml:"ServiceInfo,omitempty"`
}

type ServiceDescription struct {
	RfEndpointID   string `json:"RedfishEndpointID" yaml:"RedfishEndpointID"` // Key
	RedfishType    string `json:"RedfishType" yaml:"RedfishType"`       // Key
	RedfishSubtype string `json:"RedfishSubtype,omitempty" yaml:"RedfishSubtype,omitempty"`
	UUID           string `json:"UUID" yaml:"UUID"`
	OdataID        string `json:"OdataID" yaml:"OdataID"`
}
