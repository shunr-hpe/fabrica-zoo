// Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"github.com/openchami/fabrica/pkg/fabrica"
)

// Named represents a named resource
type Named struct {
	APIVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata `json:"metadata" yaml:"metadata"`
	Spec       NamedSpec        `json:"spec" yaml:"spec" validate:"required"`
	Status     NamedStatus      `json:"status,omitempty" yaml:"status,omitempty"`
}

// +fabrica:resource
// +fabrica:storage=dedicated
// +fabrica:index:fields=AltName,Name:name=idx_named_altname_name:unique:type=btree
type NamedSpec struct {
	// +fabrica:field:unique
	// +fabrica:field:index=btree:name=idx_named_name
	Name string `json:"name" yaml:"name"`

	// +fabrica:field:notnull
	// +fabrica:field:size=64
	AltName string `json:"altName" yaml:"altName"`

	// +fabrica:field:default=0
	Number int `json:"number" yaml:"number"`

	// +fabrica:field:nullable
	SomethingOrNothing string `json:"somethingOrNothing,omitempty" yaml:"somethingOrNothing,omitempty"`
}

// NamedStatus defines the observed state of Named
type NamedStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
	// Add your status fields here
}

// Validate implements custom validation logic for Named
func (r *Named) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}

// GetKind returns the kind of the resource
func (r *Named) GetKind() string {
	return "Named"
}

// GetName returns the name of the resource
func (r *Named) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *Named) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *Named) IsHub() {}
