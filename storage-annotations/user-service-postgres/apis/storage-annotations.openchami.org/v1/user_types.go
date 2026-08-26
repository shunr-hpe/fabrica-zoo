// Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"

	"github.com/openchami/fabrica/pkg/fabrica"
)

// User represents a user resource
type User struct {
	APIVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata `json:"metadata" yaml:"metadata"`
	Spec       UserSpec         `json:"spec" yaml:"spec" validate:"required"`
	Status     UserStatus       `json:"status,omitempty" yaml:"status,omitempty"`
}

// UserSpec is a system user with authentication credentials. It exercises the
// PR-106 storage annotations: dedicated (hybrid) storage where each annotated
// scalar spec field is promoted to a real spec_* column, plus a composite
// index over username+email.
//
// +fabrica:resource
// +fabrica:storage=dedicated
// +fabrica:index:fields=Username,Email:name=idx_user_login:unique:type=btree
type UserSpec struct {
	// Username is the unique identifier for this user; immutable after creation.
	// +fabrica:field:immutable
	// +fabrica:field:unique
	// +fabrica:field:index
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`

	// Email must be unique across all users.
	// +fabrica:field:unique
	// +fabrica:field:index
	Email string `json:"email" validate:"required,email"`

	// Password is plaintext in the request and bcrypt-hashed before storage.
	// +fabrica:field:storage=hashed:bcrypt:cost=12
	// +fabrica:field:sensitive
	// +fabrica:field:immutable
	Password string `json:"password" validate:"required,min=8"`

	// FullName is the display name; not annotated, so it stays in spec_data JSON.
	FullName string `json:"fullName" validate:"required,min=1,max=128"`

	// Role defines the permission level.
	// +fabrica:field:default=user
	// +fabrica:field:index
	Role string `json:"role" validate:"required,oneof=admin user readonly"`

	// Active indicates whether the account is enabled.
	// +fabrica:field:default=true
	Active bool `json:"active"`
}

// UserStatus defines the observed state of User
type UserStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
	// Add your status fields here
}

// Validate implements custom validation logic for User
func (r *User) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}

// GetKind returns the kind of the resource
func (r *User) GetKind() string {
	return "User"
}

// GetName returns the name of the resource
func (r *User) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *User) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *User) IsHub() {}
