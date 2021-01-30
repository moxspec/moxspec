package nw

import (
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("nw")
}

// Controller represents a network controller
type Controller struct {
	Path   string
	Driver string
	netDir string
	Port   Port
}

// Port represents a network port
type Port struct {
	Name    string
	Path    string
	HWAddr  string
	IPAddrs []*IPAddr
	MTU     uint32
	Speed   uint32
	State   string
	Carrier bool
}

// IPAddr represents an ip address
type IPAddr struct {
	Ver       byte
	Addr      string
	Netmask   string
	MaskSize  int
	Network   string
	Broadcast string
}

// NewDecoder creates and initializes Controller as Decoder
func NewDecoder(path, driver string) *Controller {
	c := new(Controller)
	c.Path = path
	c.Driver = driver

	switch c.Driver {
	case "virtio-pci":
		c.netDir = filepath.Join(c.Path, "virtio0", "net")
	default:
		c.netDir = filepath.Join(c.Path, "net")
	}

	return c
}

// Decode makes Controller satisfy the mox.Decoder interface
func (c *Controller) Decode() error {
	log.Debugf("scanning: %s", c.Path)

	ifname := getIntfName(c.netDir)
	if ifname == "" {
		return fmt.Errorf("could not find interface")
	}

	log.Debugf("found: %s", ifname)

	p := new(Port)
	p.Name = ifname

	intf, err := net.InterfaceByName(p.Name)
	if err != nil {
		return err
	}

	p.MTU = uint32(intf.MTU)

	addrs, err := intf.Addrs()
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		i := new(IPAddr)

		a, anet, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}

		if a.To4() == nil {
			i.Ver = 6
		} else {
			i.Ver = 4
		}

		i.Addr = a.String()
		i.Netmask = net.IP(anet.Mask).String()
		i.MaskSize, _ = anet.Mask.Size()
		i.Network = anet.IP.String()

		if i.Ver == 4 {
			bcast, err := calcIPv4Broadcast(a, anet)
			if err == nil {
				i.Broadcast = bcast
			}
		}

		log.Debugf("addr: %+v", i)

		p.IPAddrs = append(p.IPAddrs, i)
	}

	p.Path = filepath.Join(c.netDir, p.Name)
	p.HWAddr, _ = util.LoadString(filepath.Join(p.Path, "address"))
	p.Speed, _ = util.LoadUint32(filepath.Join(p.Path, "speed"))
	p.State, _ = util.LoadString(filepath.Join(p.Path, "operstate"))

	cr, _ := util.LoadByte(filepath.Join(p.Path, "carrier"))
	if cr == 1 {
		p.Carrier = true
	}

	c.Port = *p

	return nil
}

func getIntfName(ndir string) string {
	files, err := ioutil.ReadDir(ndir)
	if err != nil {
		log.Debug(err)
		return ""
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		return file.Name()
	}

	return ""
}

func calcIPv4Broadcast(addr net.IP, anet *net.IPNet) (string, error) {
	var bs []byte
	if addr.To4() == nil {
		bs = []byte(addr.To16())
	} else {
		bs = []byte(addr.To4())
	}

	bits := len(bs) * 8
	size, _ := anet.Mask.Size()
	asize := bits - size // address size

	if size == 32 || size == 0 {
		return "", fmt.Errorf("invalid mask size")
	}

	var masked int
	for i := len(bs) - 1; i >= 0; i-- {
		var j byte
		for j = 0; j < 8; j++ {
			masked++
			bs[i] |= (1 << j)
			if asize == masked {
				break
			}
		}
		if asize == masked {
			break
		}
	}

	return net.IP(bs).String(), nil
}
