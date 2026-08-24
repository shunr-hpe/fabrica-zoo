// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"encoding/json"

	"github.com/openchami/fabrica/pkg/fabrica"
)

// Component represents a component resource
type Component struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       ComponentSpec    `json:"spec" validate:"required"`
	Status     ComponentStatus  `json:"status,omitempty"`
}

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

// ComponentStatus defines the observed state of Component
type ComponentStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
	// Add your status fields here
}

// Validate implements custom validation logic for Component
func (r *Component) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}

// GetKind returns the kind of the resource
func (r *Component) GetKind() string {
	return "Component"
}

// GetName returns the name of the resource
func (r *Component) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *Component) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *Component) IsHub() {}
