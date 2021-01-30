package ipmi

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/platform"
	"github.com/actapio/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("ipmi")
}

// Device represents a ipmi device
type Device struct {
	Path     string
	Firmware string
	MAC      string
	IPAddr   string
	Netmask  string
	Gateway  string
}

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

	cmdpath, err := exec.LookPath("ipmitool")
	if err != nil {
		return fmt.Errorf("ipmitool is not installed")
	}

	rev, err := util.LoadString("/sys/class/ipmi/ipmi0/device/bmc/firmware_revision")
	if err == nil {
		d.Firmware = rev
	}

	log.Debugf("running %s lan print", cmdpath)
	res, _ := util.Exec(cmdpath, "lan", "print") // ipmitool always return 1
	d.IPAddr, d.Netmask, d.Gateway, d.MAC = parseLanPrint(res)

	return nil
}

func parseLanPrint(res string) (ipaddr, netmask, gateway, mac string) {
	for _, l := range strings.Split(res, "\n") {
		log.Debug(l)

		kv := strings.SplitN(l, ":", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		switch key {
		case "IP Address":
			ipaddr = val
		case "Subnet Mask":
			netmask = val
		case "Default Gateway IP":
			gateway = val
		case "MAC Address":
			mac = val
		}
	}
	return
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
