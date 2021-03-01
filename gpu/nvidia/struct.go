package nvidia

import (
	"encoding/xml"

	"github.com/moxspec/moxspec/pci"
	"github.com/moxspec/moxspec/util"
)

// Devices represents GPUs
type Devices struct {
	dict map[string]*GPU
}

// Decode makes Device satisfy the mox.Decoder interface
func (d *Devices) Decode() error {
	log.Debug("executing nvidia-smi")
	res, err := util.Exec("nvidia-smi", "-q", "-x")
	if err != nil {
		return err
	}

	smlog, err := decodeLog(res)
	if err != nil {
		return err
	}

	for _, g := range smlog.GPUs {
		log.Debugf("found %s %s, sn:%s, uuid:%s, vbios:%s", g.ID, g.ProductName, g.Serial, g.UUID, g.VBIOSVersion)

		dom, bus, dev, fun, err := pci.ParseLocater(g.ID)
		if err != nil {
			log.Debugf("invalid gpu id: %s", g.ID)
			continue
		}
		ids := pci.IDString(dom, bus, dev, fun)
		log.Debugf("pci-id: %s", ids)

		d.dict[ids] = g
	}

	return nil
}

// GetGPU returns the GPU which has given pci-id
func (d Devices) GetGPU(pciid string) *GPU {
	if g, ok := d.dict[pciid]; ok {
		return g
	}
	return nil
}

// SmiLog represents nvidia-smi log format
type SmiLog struct {
	XMLName       xml.Name `xml:"nvidia_smi_log"`
	DriverVersion string   `xml:"driver_version"`
	GPUs          []*GPU   `xml:"gpu"`
}

// GPU represents the GPU appears in nvidia-smi log
type GPU struct {
	baseSpec
	utilSpec
	eccErrorsSpec
	tempSpec
	powerSpec
}

// UnmarshalXML imprements xml.Unmarshaler interface
func (g *GPU) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	raw := struct {
		baseSpec
		utilSpecRaw
		eccErrorsSpecRaw
		tempSpecRaw
		powerSpecRaw
	}{}

	err := d.DecodeElement(&raw, &start)
	if err != nil {
		return err
	}

	*g = GPU{}
	g.baseSpec = raw.baseSpec
	g.utilSpec = raw.utilSpecRaw.convert()
	g.eccErrorsSpec = raw.eccErrorsSpecRaw.convert()
	g.tempSpec = raw.tempSpecRaw.convert()
	g.powerSpec = raw.powerSpecRaw.convert()
	return nil
}

type baseSpec struct {
	ID           string `xml:"id,attr"`
	ProductName  string `xml:"product_name"`
	ProductBrand string `xml:"product_brand"`
	Serial       string `xml:"serial"`
	UUID         string `xml:"uuid"`
	VBIOSVersion string `xml:"vbios_version"`
}

type utilSpec struct {
	Util struct {
		GPU    float32
		Memory float32
	}
}

type utilSpecRaw struct {
	Util struct {
		GPU    string `xml:"gpu_util"`
		Memory string `xml:"memory_util"`
	} `xml:"utilization"`
}

func (u utilSpecRaw) convert() utilSpec {
	c := utilSpec{}
	c.Util.GPU = convUtilString(u.Util.GPU)
	c.Util.Memory = convUtilString(u.Util.Memory)
	return c
}

type eccErrorsSpec struct {
	ECCErrors struct {
		Volatile  ECCError
		Aggregate ECCError
	}
}

type eccErrorsSpecRaw struct {
	ECCErrors struct {
		Volatile  eccErrorRaw `xml:"volatile"`
		Aggregate eccErrorRaw `xml:"aggregate"`
	} `xml:"ecc_errors"`
}

func (e eccErrorsSpecRaw) convert() eccErrorsSpec {
	c := eccErrorsSpec{}
	c.ECCErrors.Volatile = e.ECCErrors.Volatile.convert()
	c.ECCErrors.Aggregate = e.ECCErrors.Aggregate.convert()
	return c
}

// ECCError represents ecc error counters
// single bit errors are corrected
// double bit errors are uncorrectable
type ECCError struct {
	SingleBit ECCCounter
	DoubleBit ECCCounter
}

type eccErrorRaw struct {
	SingleBit eccCounterRaw `xml:"single_bit"`
	DoubleBit eccCounterRaw `xml:"double_bit"`
}

func (e eccErrorRaw) convert() ECCError {
	c := ECCError{}
	c.SingleBit = e.SingleBit.convert()
	c.DoubleBit = e.DoubleBit.convert()
	return c
}

// ECCCounter represents ECC error counters
type ECCCounter struct {
	DeviceMemory int
	RegisterFile int
	L1Cache      int
	L2Cache      int
	Total        int
}

type eccCounterRaw struct {
	DeviceMemory string `xml:"device_memory"`
	RegisterFile string `xml:"register_file"`
	L1Cache      string `xml:"l1_cache"`
	L2Cache      string `xml:"l2_cache"`
	Total        string `xml:"total"`
}

func (e eccCounterRaw) convert() ECCCounter {
	c := ECCCounter{}
	c.DeviceMemory = convCountString(e.DeviceMemory)
	c.RegisterFile = convCountString(e.RegisterFile)
	c.L1Cache = convCountString(e.L1Cache)
	c.L2Cache = convCountString(e.L2Cache)
	c.Total = convCountString(e.Total)
	return c
}

type tempSpec struct {
	Temp struct {
		GPU    float32
		Memory float32
	}
}

type tempSpecRaw struct {
	Temp struct {
		GPU    string `xml:"gpu_temp"`
		Memory string `xml:"memory_temp"`
	} `xml:"temperature"`
}

func (t tempSpecRaw) convert() tempSpec {
	c := tempSpec{}
	c.Temp.GPU = convTempString(t.Temp.GPU)
	c.Temp.Memory = convTempString(t.Temp.Memory)
	return c
}

type powerSpec struct {
	Power struct {
		Draw  float32
		Limit float32
	}
}

type powerSpecRaw struct {
	Power struct {
		Draw  string `xml:"power_draw"`
		Limit string `xml:"power_limit"`
	} `xml:"power_readings"`
}

func (p powerSpecRaw) convert() powerSpec {
	c := powerSpec{}
	c.Power.Draw = convWattString(p.Power.Draw)
	c.Power.Limit = convWattString(p.Power.Limit)
	return c
}
