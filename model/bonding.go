package model

import (
	"fmt"
	"net"
)

// BondInterface represents a bonding device
type BondInterface struct {
	Name      string    `json:"name"`
	Slaves    []string  `json:"slaves,omitempty"`
	LinkAttrs LinkAttrs `json:"link_attrs,omitempty"`
	BondAttrs BondAttrs `json:"bond_attrs,omitempty"`
}

// LinkAttrs represents link attributes
type LinkAttrs struct {
	State   string       `json:"state,omitempty"`
	HWAddr  string       `json:"hwaddr,omitempty"`
	MTU     int          `json:"mtu"`
	TxQLen  int          `json:"tx_qlen"`
	IPAddrs []*IPAddress `json:"ipaddrs,omitempty"`
}

// BondAttrs represents bonding attributes
type BondAttrs struct {
	Mode            string   `json:"mode,omitempty"`
	ActiveSlave     string   `json:"active_slaves,omitempty"`
	Miimon          int      `json:"miimon,omitempty"`
	UpDelay         int      `json:"updelay,omitempty"`
	DownDelay       int      `json:"downdelay,omitempty"`
	UseCarrier      int      `json:"use_carrier,omitempty"`
	ArpInterval     int      `json:"arp_interval,omitempty"`
	ArpIPTargets    []net.IP `json:"arp_ip_target,omitempty"`
	ArpValidate     string   `json:"arp_validate,omitempty"`
	ArpAllTargets   string   `json:"arp_all_targets,omitempty"`
	Primary         string   `json:"primary,omitempty"`
	PrimaryReselect string   `json:"primary_reselect,omitempty"`
	FailOverMac     string   `json:"fail_over_mac,omitempty"`
	XmitHashPolicy  string   `json:"xmit_hash_policy,omitempty"`
	LacpRate        string   `json:"lacp_rate,omitempty"`
}

// Summary returns summarized string
func (n BondInterface) Summary() string {
	return fmt.Sprintf("Bond Device %s", n.Name)
}
