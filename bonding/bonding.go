package bonding

import (
	"net"

	"github.com/vishvananda/netlink"
)

// BondInterface represents a bonding device
type BondInterface struct {
	Name      string
	Slaves    []string
	AddrList  []netlink.Addr
	LinkAttrs netlink.LinkAttrs
	BondAttrs BondAttrs
}

// BondAttrs represents bonding attributes
type BondAttrs struct {
	Mode            netlink.BondMode
	ActiveSlave     string
	Miimon          int
	UpDelay         int
	DownDelay       int
	UseCarrier      int
	ArpInterval     int
	ArpIPTargets    []net.IP
	ArpValidate     netlink.BondArpValidate
	ArpAllTargets   netlink.BondArpAllTargets
	Primary         string
	PrimaryReselect netlink.BondPrimaryReselect
	FailOverMac     netlink.BondFailOverMac
	XmitHashPolicy  netlink.BondXmitHashPolicy
	ResendIgmp      int
	NumPeerNotif    int
	AllSlavesActive int
	MinLinks        int
	LpInterval      int
	PacketsPerSlave int
	LacpRate        netlink.BondLacpRate
	AdSelect        netlink.BondAdSelect
	AdInfo          *netlink.BondAdInfo
	AdActorSysPrio  int
	AdUserPortKey   int
	AdActorSystem   net.HardwareAddr
	TlbDynamicLb    int
}

// Decode make BondInterface satisfy the mox.Decoder interface
func (intf *BondInterface) Decode() error {

	bli, err := netlink.LinkByName(intf.Name)
	if err != nil {
		return err
	}

	slaves, err := findBondSlaves(bli.Attrs().Index)
	if err != nil {
		return err
	}

	intf.Slaves = slaves

	stats, err := getBondParameters(bli)
	if err != nil {
		return err
	}

	if stats != nil {
		intf.LinkAttrs = stats.LinkAttrs

		intf.BondAttrs.Mode = stats.Mode

		activeslave, err := netlink.LinkByIndex(stats.ActiveSlave)
		if err == nil {
			intf.BondAttrs.ActiveSlave = activeslave.Attrs().Name
		}

		intf.BondAttrs.Miimon = stats.Miimon
		intf.BondAttrs.UpDelay = stats.UpDelay
		intf.BondAttrs.DownDelay = stats.DownDelay
		intf.BondAttrs.UseCarrier = stats.UseCarrier
		intf.BondAttrs.ArpInterval = stats.ArpInterval
		intf.BondAttrs.ArpIPTargets = stats.ArpIpTargets
		intf.BondAttrs.ArpValidate = stats.ArpValidate
		intf.BondAttrs.ArpAllTargets = stats.ArpAllTargets

		primary, err := netlink.LinkByIndex(stats.Primary)
		if err == nil {
			intf.BondAttrs.ActiveSlave = primary.Attrs().Name
		}

		intf.BondAttrs.PrimaryReselect = stats.PrimaryReselect
		intf.BondAttrs.FailOverMac = stats.FailOverMac
		intf.BondAttrs.XmitHashPolicy = stats.XmitHashPolicy
		intf.BondAttrs.ResendIgmp = stats.ResendIgmp
		intf.BondAttrs.NumPeerNotif = stats.NumPeerNotif
		intf.BondAttrs.AllSlavesActive = stats.AllSlavesActive
		intf.BondAttrs.MinLinks = stats.MinLinks
		intf.BondAttrs.LpInterval = stats.LpInterval
		intf.BondAttrs.PacketsPerSlave = stats.PacketsPerSlave
		intf.BondAttrs.LacpRate = stats.LacpRate
		intf.BondAttrs.AdSelect = stats.AdSelect
		intf.BondAttrs.AdInfo = stats.AdInfo
		intf.BondAttrs.AdActorSysPrio = stats.AdActorSysPrio
		intf.BondAttrs.AdUserPortKey = stats.AdUserPortKey
		intf.BondAttrs.AdActorSystem = stats.AdActorSystem
		intf.BondAttrs.TlbDynamicLb = stats.TlbDynamicLb
	}
	ipaddresses, err := netlink.AddrList(bli, netlink.FAMILY_V4)
	if err != nil {
		return err
	}
	if ipaddresses != nil {
		intf.AddrList = ipaddresses
	}

	ipaddresses, err = netlink.AddrList(bli, netlink.FAMILY_V6)
	if err != nil {
		return err
	}
	if ipaddresses != nil {
		for _, addr := range ipaddresses {
			intf.AddrList = append(intf.AddrList, addr)
		}
	}

	return nil
}

// NewDecoder creates and initializes a BondInterface as Decoder
func NewDecoder(name string) *BondInterface {
	intf := new(BondInterface)
	intf.Name = name
	return intf
}
