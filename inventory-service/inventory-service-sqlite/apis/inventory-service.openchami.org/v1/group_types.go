// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"

	"github.com/openchami/fabrica/pkg/fabrica"
)

// Group represents a group resource
type Group struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       GroupSpec        `json:"spec" validate:"required"`
	Status     GroupStatus      `json:"status,omitempty"`
}

type GroupSpec struct {
	// Description    string   `json:"description" yaml:"description"`
	Label          string   `json:"label" yaml:"label"`
	ExclusiveGroup string   `json:"exclusiveGroup,omitempty" yaml:"exclusiveGroup,omitempty"`
	Tags           []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Members        Members  `json:"members" yaml:"members"` // List of xnames, required.

	// Private
	// normalized bool
	// verified   bool
}

// GroupStatus defines the observed state of Group
type GroupStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
	// Add your status fields here
}

// Validate implements custom validation logic for Group
func (r *Group) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}

// GetKind returns the kind of the resource
func (r *Group) GetKind() string {
	return "Group"
}

// GetName returns the name of the resource
func (r *Group) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *Group) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *Group) IsHub() {}

type Members struct {
	IDs []string `json:"ids" yaml:"ids"` // xname array

	// Private
	// normalized bool
	// verified   bool
}

type Membership struct {
	ID            string   `json:"id" yaml:"id"`
	GroupLabels   []string `json:"groupLabels" yaml:"groupLabels"`
	PartitionName string   `json:"partitionName" yaml:"partitionName"`
}
