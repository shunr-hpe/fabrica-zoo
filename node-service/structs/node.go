// Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package structs

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
