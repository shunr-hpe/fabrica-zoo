package smd_structs

type GroupSpec struct {
	// Description    string   `json:"description" yaml:"description"`

	// Label is the group label and serves as the resource identifier.
	//
	// +fabrica:field:index
	// +fabrica:field:unique
	Label          string   `json:"label" yaml:"label"`
	ExclusiveGroup string   `json:"exclusiveGroup,omitempty" yaml:"exclusiveGroup,omitempty"`
	Tags           []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Members        Members  `json:"members" yaml:"members"` // List of xnames, required.

	// Private
	// normalized bool
	// verified   bool
}

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
