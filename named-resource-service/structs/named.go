// Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package structs

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
