// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"

	"github.com/openchami/fabrica/pkg/fabrica"
)

// EthernetInterface represents a ethernetinterface resource
type EthernetInterface struct {
	APIVersion string                  `json:"apiVersion"`
	Kind       string                  `json:"kind"`
	Metadata   fabrica.Metadata        `json:"metadata"`
	Spec       EthernetInterfaceSpec   `json:"spec" validate:"required"`
	Status     EthernetInterfaceStatus `json:"status,omitempty"`
}

type EthernetInterfaceSpec struct {
	// Description            string              `json:"Description" yaml:"Description"`
	ID          string      `json:"ID" yaml:"ID"`
	MACAddr     string      `json:"MACAddress" yaml:"MACAddress"`
	LastUpdate  string      `json:"LastUpdate" yaml:"LastUpdate"`
	CompID      string      `json:"ComponentID" yaml:"ComponentID"`
	Type        string      `json:"Type" yaml:"Type"`
	IPAddresses []IPAddress `json:"IPAddresses" yaml:"IPAddresses"`
}

// EthernetInterfaceStatus defines the observed state of EthernetInterface
type EthernetInterfaceStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
	// Add your status fields here
}

// Validate implements custom validation logic for EthernetInterface
func (r *EthernetInterface) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}

// GetKind returns the kind of the resource
func (r *EthernetInterface) GetKind() string {
	return "EthernetInterface"
}

// GetName returns the name of the resource
func (r *EthernetInterface) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *EthernetInterface) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *EthernetInterface) IsHub() {}

type IPAddress struct {
	IPAddress string `json:"IPAddress" yaml:"IPAddress"`
	Network   string `json:"Network,omitempty" yaml:"Network,omitempty"`
}
