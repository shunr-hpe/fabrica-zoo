package smd_structs

type RedfishEndpointSpec struct {
	ID                 string        `json:"ID" yaml:"ID"`
	Type               string        `json:"Type" yaml:"Type"`
	Name               string        `json:"Name,omitempty" yaml:"Name,omitempty"` // user supplied descriptive name
	Hostname           string        `json:"Hostname" yaml:"Hostname"`
	Domain             string        `json:"Domain" yaml:"Domain"`
	FQDN               string        `json:"FQDN" yaml:"FQDN"`
	Enabled            bool          `json:"Enabled" yaml:"Enabled"`
	UUID               string        `json:"UUID,omitempty" yaml:"UUID,omitempty"`
	User               string        `json:"User" yaml:"User"`
	Password           string        `json:"Password" yaml:"Password"` // Temporary until more secure method
	UseSSDP            bool          `json:"UseSSDP,omitempty" yaml:"UseSSDP,omitempty"`
	MACRequired        bool          `json:"MACRequired,omitempty" yaml:"MACRequired,omitempty"`
	MACAddr            string        `json:"MACAddr,omitempty" yaml:"MACAddr,omitempty"`
	IPAddress          string        `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty"`
	RedsicoverOnUpdate bool          `json:"RediscoverOnUpdate" yaml:"RediscoverOnUpdate"`
	TemplateID         string        `json:"TemplateID,omitempty" yaml:"TemplateID,omitempty"`
	DiscoveryInfo      DiscoveryInfo `json:"DiscoveryInfo" yaml:"DiscoveryInfo"`
}

type DiscoveryInfo struct {
	LastAttempt    string `json:"LastDiscoveryAttempt,omitempty" yaml:"LastDiscoveryAttempt,omitempty"`
	LastStatus     string `json:"LastDiscoveryStatus" yaml:"LastDiscoveryStatus"`
	RedfishVersion string `json:"RedfishVersion,omitempty" yaml:"RedfishVersion,omitempty"`
}
