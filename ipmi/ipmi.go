package ipmi

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/platform"
	"github.com/moxspec/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("ipmi")
}

// Device represents a ipmi device
type Device struct {
	Path          string
	Firmware      string
	MAC           string
	IPAddr        string
	Netmask       string
	Gateway       string
	AddressSource addrSrcType
	VLANID        uint16
}

type addrSrcType string

const (
	srcUnspecified addrSrcType = "Unspecified"                   // 0x00
	srcStatic      addrSrcType = "Static"                        // 0x01
	srcDHCP        addrSrcType = "DHCP"                          // 0x02
	srcBIOS        addrSrcType = "BIOS or System software"       // 0x03
	srcOther       addrSrcType = "Other address assignment pool" // 0x04
)

// NewDecoder creates and initializes a Device as Decoder
func NewDecoder() *Device {
	d := new(Device)
	d.Path = "/dev/ipmi0"
	return d
}

// Decode makes Device satisfy the mox.Decoder interface
func (d *Device) Decode() error {
	if !platform.IsLoadedModule("ipmi_devintf") {
		return fmt.Errorf("kernel module for ipmi is not loaded")
	}

	if !util.Exists(d.Path) {
		return fmt.Errorf("ipmi device was not found")
	}

	_, err := exec.LookPath("ipmitool")
	if err != nil {
		return fmt.Errorf("ipmitool is not installed")
	}

	d.Firmware = getFirmwareRev()
	d.MAC = getMACAddress()
	d.IPAddr = getIPAddress()
	d.Netmask = getSubnetMask()
	d.Gateway = getDefaultGateway()
	d.VLANID = getVLANID()
	d.AddressSource = getAddressSource()

	return nil
}

func getFirmwareRev() string {
	res, err := util.Exec("ipmitool", "raw", "0x06", "0x01")
	if err != nil {
		return ""
	}
	return parseMcInfoRaw(res)
}

func parseMcInfoRaw(in string) string {
	codes := strings.Fields(in)
	if len(codes) < 4 {
		return ""
	}
	fr1, _ := strconv.ParseUint(codes[2], 16, 8)
	maj := fr1 & 0x7F // binary encoded
	min := codes[3]   // BCD encoded
	return fmt.Sprintf("%d.%s", maj, min)
}

func getMACAddress() string {
	res, err := util.Exec("ipmitool", "raw", "0x0c", "0x02", "0x01", "0x05", "  0x00", "0x00")
	if err != nil {
		return ""
	}
	return parseMACAddressRaw(res)
}

func parseMACAddressRaw(in string) string {
	codes := strings.Fields(in)
	if len(codes) != 7 {
		return ""
	}
	return strings.Join(codes[1:], ":")
}

func getIPAddress() string {
	res, err := util.Exec("ipmitool", "raw", "0x0c", "0x02", "0x01", "0x03", "  0x00", "0x00")
	if err != nil {
		return ""
	}
	return parseIPAddressRaw(res)
}

func getSubnetMask() string {
	res, err := util.Exec("ipmitool", "raw", "0x0c", "0x02", "0x01", "0x06", "  0x00", "0x00")
	if err != nil {
		return ""
	}
	return parseIPAddressRaw(res)
}

func getDefaultGateway() string {
	res, err := util.Exec("ipmitool", "raw", "0x0c", "0x02", "0x01", "0x0c", "  0x00", "0x00")
	if err != nil {
		return ""
	}
	return parseIPAddressRaw(res)
}

func parseIPAddressRaw(in string) string {
	codes := strings.Fields(in)
	if len(codes) != 5 { // TBC: IPv6
		return ""
	}

	a, _ := strconv.ParseUint(codes[1], 16, 8)
	b, _ := strconv.ParseUint(codes[2], 16, 8)
	c, _ := strconv.ParseUint(codes[3], 16, 8)
	d, _ := strconv.ParseUint(codes[4], 16, 8)
	return net.IPv4(byte(a), byte(b), byte(c), byte(d)).String()
}

func getVLANID() uint16 {
	res, err := util.Exec("ipmitool", "raw", "0x0c", "0x02", "0x01", "0x14", "  0x00", "0x00")
	if err != nil {
		return 0
	}
	return parseVLANID(res)
}

func parseVLANID(in string) uint16 {
	codes := strings.Fields(in)
	if len(codes) != 3 {
		return 0
	}
	d1, _ := strconv.ParseUint(codes[1], 16, 8)
	d2, _ := strconv.ParseUint(codes[2], 16, 8)

	if byte(d2)&0x80 != 0x80 { // 1 == enabled, 0 == disabled
		return 0
	}

	return ((uint16(d2)&0x0F)<<8 | uint16(d1))
}

func getAddressSource() addrSrcType {
	res, err := util.Exec("ipmitool", "raw", "0x0c", "0x02", "0x01", "0x04", "  0x00", "0x00")
	if err != nil {
		return "uns[ecified"
	}
	return parseAddressSource(res)
}

func parseAddressSource(in string) addrSrcType {
	codes := strings.Fields(in)
	if len(codes) != 2 {
		return srcUnspecified
	}

	d1, _ := strconv.ParseUint(codes[1], 16, 8)
	switch byte(d1) & 0x0F {
	case 0x00:
		return srcUnspecified
	case 0x01:
		return srcStatic
	case 0x02:
		return srcDHCP
	case 0x03:
		return srcBIOS
	case 0x04:
		return srcOther
	}
	return srcUnspecified
}

// GetSEL returns system event list
func GetSEL() []string {
	log.Debug("running ipmitool sel list")
	res, err := util.Exec("ipmitool", "sel", "list")
	if err != nil {
		return nil
	}

	return strings.Split(res, "\n")
}
