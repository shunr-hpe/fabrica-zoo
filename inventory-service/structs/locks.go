package smd_structs

// LockSpec is the source-of-truth definition for the SMD-compatible lock /
// reservation record. There is at most one Lock per Component, keyed by the
// component xname (ID). Admin lock flags (Locked, ReservationDisabled) live on
// the Component itself; this record only tracks an active reservation.
type LockSpec struct {
	ID             string `json:"ID" yaml:"ID"`
	ReservationKey string `json:"ReservationKey,omitempty" yaml:"ReservationKey,omitempty"`
	DeputyKey      string `json:"DeputyKey,omitempty" yaml:"DeputyKey,omitempty"`
	CreatedTime    string `json:"CreatedTime,omitempty" yaml:"CreatedTime,omitempty"`
	// ExpirationTime is an RFC3339 timestamp. An empty value denotes an
	// indefinite (admin) reservation.
	ExpirationTime string `json:"ExpirationTime,omitempty" yaml:"ExpirationTime,omitempty"`
}
