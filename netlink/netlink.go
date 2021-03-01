package netlink

import (
	"github.com/moxspec/moxspec/loglet"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("netlink")
}

// Interface represents the ethernet interface
type Interface struct {
	Name  string
	Stats RtnlLinkStats64
}

// RtnlLinkStats64 represents interaface stats
// cf. include/uapi/linux/if_link.h
type RtnlLinkStats64 struct {
	RxPackets         uint64
	TxPackets         uint64
	RxBytes           uint64
	TxBytes           uint64
	RxErrors          uint64
	TxErrors          uint64
	RxDropped         uint64
	TxDropped         uint64
	Multicast         uint64
	Collisions        uint64
	RxLengthErrors    uint64
	RxOverErrors      uint64
	RxCrcErrors       uint64
	RxFrameErrors     uint64
	RxFifoErrors      uint64
	RxMissedErrors    uint64
	TxAbortedErrors   uint64
	TxCarrierErrors   uint64
	TxFifoErrors      uint64
	TxHeartbeatErrors uint64
	TxWindowErrors    uint64
	RxCompressed      uint64
	TxCompressed      uint64
	RxNohandler       uint64
}

// Decode make Interface satisfy the mox.Decoder interface
func (intf *Interface) Decode() error {
	nli, err := newNetlinkInterface()
	if err != nil {
		return err
	}
	defer nli.close()

	stats, err := getStats(nli, intf.Name)
	if err != nil {
		return err
	}

	if stats != nil {
		intf.Stats = *stats
	}

	return nil
}

// NewDecoder creates and initializes a Interface as Decoder
func NewDecoder(name string) *Interface {
	intf := new(Interface)
	intf.Name = name
	return intf
}
