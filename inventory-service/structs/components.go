package smd_structs

import "encoding/json"

type ComponentSpec struct {
	ID                  string      `json:"ID" yaml:"ID"`
	Type                string      `json:"Type" yaml:"Type"`
	State               string      `json:"State,omitempty" yaml:"State,omitempty"`
	Flag                string      `json:"Flag,omitempty" yaml:"Flag,omitempty"`
	Enabled             *bool       `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`
	SwStatus            string      `json:"SoftwareStatus,omitempty" yaml:"SoftwareStatus,omitempty"`
	Role                string      `json:"Role,omitempty" yaml:"Role,omitempty"`
	SubRole             string      `json:"SubRole,omitempty" yaml:"SubRole,omitempty"`
	NID                 json.Number `json:"NID,omitempty" yaml:"NID,omitempty"`
	Subtype             string      `json:"Subtype,omitempty" yaml:"Subtype,omitempty"`
	NetType             string      `json:"NetType,omitempty" yaml:"NetType,omitempty"`
	Arch                string      `json:"Arch,omitempty" yaml:"Arch,omitempty"`
	Class               string      `json:"Class,omitempty" yaml:"Class,omitempty"`
	ReservationDisabled bool        `json:"ReservationDisabled,omitempty" yaml:"ReservationDisabled,omitempty"`
	Locked              bool        `json:"Locked,omitempty" yaml:"Locked,omitempty"`
}
