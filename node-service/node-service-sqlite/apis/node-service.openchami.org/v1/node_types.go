// Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"

	"github.com/openchami/fabrica/pkg/fabrica"
)

// Node represents a node resource
type Node struct {
	APIVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata `json:"metadata" yaml:"metadata"`
	Spec       NodeSpec         `json:"spec" yaml:"spec" validate:"required"`
	Status     NodeStatus       `json:"status,omitempty" yaml:"status,omitempty"`
}

// +fabrica:resource
// +fabrica:storage=dedicated
// +fabrica:index:fields=ID:name=idx_id:unique:type=btree
type NodeSpec struct {
	// +fabrica:field:unique
	// +fabrica:field:index=btree:name=idx_named_name
	ID string `json:"id" yaml:"id"`

	// +fabrica:field:default=0
	Number int `json:"number" yaml:"number"`

	// +fabrica:field:nullable
	SomethingOrNothing string `json:"somethingOrNothing,omitempty" yaml:"somethingOrNothing,omitempty"`
}

// NodeStatus defines the observed state of Node
type NodeStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
	// Add your status fields here
}

// Validate implements custom validation logic for Node
func (r *Node) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}

// GetKind returns the kind of the resource
func (r *Node) GetKind() string {
	return "Node"
}

// GetName returns the name of the resource
func (r *Node) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *Node) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *Node) IsHub() {}
