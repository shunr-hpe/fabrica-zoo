// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"encoding/json"

	"github.com/openchami/fabrica/pkg/fabrica"
)

// ServiceEndpoint represents a serviceendpoint resource
type ServiceEndpoint struct {
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   fabrica.Metadata      `json:"metadata"`
	Spec       ServiceEndpointSpec   `json:"spec" validate:"required"`
	Status     ServiceEndpointStatus `json:"status,omitempty"`
}

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

// ServiceEndpointStatus defines the observed state of ServiceEndpoint
type ServiceEndpointStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
	// Add your status fields here
}

// Validate implements custom validation logic for ServiceEndpoint
func (r *ServiceEndpoint) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}

// GetKind returns the kind of the resource
func (r *ServiceEndpoint) GetKind() string {
	return "ServiceEndpoint"
}

// GetName returns the name of the resource
func (r *ServiceEndpoint) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *ServiceEndpoint) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *ServiceEndpoint) IsHub() {}

type ServiceDescription struct {
	RfEndpointID   string `json:"RedfishEndpointID" yaml:"RedfishEndpointID"` // Key
	RedfishType    string `json:"RedfishType" yaml:"RedfishType"`             // Key
	RedfishSubtype string `json:"RedfishSubtype,omitempty" yaml:"RedfishSubtype,omitempty"`
	UUID           string `json:"UUID" yaml:"UUID"`
	OdataID        string `json:"OdataID" yaml:"OdataID"`
}
