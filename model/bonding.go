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

type LinkAttrs struct {
	State   string       `json:"state,omitempty"`
	HWAddr  string       `json:"hwaddr,omitempty"`
	MTU     int          `json:"mtu"`
	TxQLen  int          `json:"tx_qlen"`
	IPAddrs []*IPAddress `json:"ipaddrs,omitempty"`
}

type BondAttrs struct {
	Mode            string   `json:"mode,omitempty"`
	ActiveSlave     string   `json:"active_slaves,omitempty"`
	Miimon          int      `json:"miimon,omitempty"`
	UpDelay         int      `json:"updelay,omitempty"`
	DownDelay       int      `json:"downdelay,omitempty"`
	UseCarrier      int      `json:"use_carrier,omitempty"`
	ArpInterval     int      `json:"arp_interval,omitempty"`
	ArpIpTargets    []net.IP `json:"arp_ip_target,omitempty"`
	ArpValidate     string   `json:"arp_validate,omitempty"`
	ArpAllTargets   string   `json:"arp_all_targets,omitempty"`
	Primary         string   `json:"primary,omitempty"`
	PrimaryReselect string   `json:"primary_reselect,omitempty"`
	FailOverMac     string   `json:"fail_over_mac,omitempty"`
	XmitHashPolicy  string   `json:"xmit_hash_policy,omitempty"`
	LacpRate        string   `json:"lacp_rate,omitempty"`

	// We don't currently need the following values, but netlink provides them
	// Just put here as a reference
	// MiiStatus       int      `json:"mii_status,omitempty"`
	// ResendIgmp      int
	// NumPeerNotif    int
	// AllSlavesActive int
	// MinLinks        int
	// LpInterval      int
	// PacketsPerSlave int
	// AdSelect        BondAdSelect
	// AdInfo          *BondAdInfo
	// AdActorSysPrio  int
	// AdUserPortKey   int
	// AdActorSystem   net.HardwareAddr
	// TlbDynamicLb    int
}

// Summary returns summarized string
func (n BondInterface) Summary() string {
	return fmt.Sprintf("Bond Device %s", n.Name)
}
