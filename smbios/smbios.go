package smbios

import (
	"fmt"

	"github.com/actapio/moxspec/loglet"
	gosmbios "github.com/digitalocean/go-smbios/smbios"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("smbios")
}

// Spec represents a smbios spec
type Spec struct {
	Major   int
	Minor   int
	Rev     int
	Records map[uint8][]*Structure
}

// Version returns the version string of smbios
func (s Spec) Version() string {
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Rev)
}

// GetBIOS returns BIOS
func (s Spec) GetBIOS() *BIOS {
	r, ok := s.Records[biosInformation]
	if ok && len(r) == 1 { // it should be only one
		return r[0].Data.(*BIOS)
	}
	return nil
}

// GetSystem returns System
func (s Spec) GetSystem() *System {
	r, ok := s.Records[systemInformation]
	if ok && len(r) == 1 { // it should be only one
		return r[0].Data.(*System)
	}
	return nil
}

// GetBaseboard returns Baseboard
func (s Spec) GetBaseboard() *Baseboard {
	r, ok := s.Records[baseboardInformation]
	if ok && len(r) == 1 { // it should be only one
		return r[0].Data.(*Baseboard)
	}
	return nil
}

// GetChassis returns Chassis
func (s Spec) GetChassis() *Chassis {
	r, ok := s.Records[systemEnclosure]
	if ok && len(r) == 1 { // it should be only one
		return r[0].Data.(*Chassis)
	}
	return nil
}

// GetProcessor returns Processor(s)
func (s Spec) GetProcessor() []*Processor {

	list := []*Processor{}
	rs, ok := s.Records[processorInformation]
	if !ok {
		return nil
	}

	// build the dictionary to bind a cache to a processor
	cDict := make(map[uint16]*Cache)
	for _, c := range s.Records[cacheInformation] {
		cDict[c.Header.Handle] = c.Data.(*Cache)
	}

	for _, r := range rs {
		pc := r.Data.(*Processor)
		if cache, ok := cDict[pc.L1CacheHandle]; ok {
			pc.L1Cache = cache
		}
		if cache, ok := cDict[pc.L2CacheHandle]; ok {
			pc.L2Cache = cache
		}
		if cache, ok := cDict[pc.L3CacheHandle]; ok {
			pc.L3Cache = cache
		}
		list = append(list, pc)
	}
	return list
}

// GetMemoryDevice returns MemoryDevice(s)
func (s Spec) GetMemoryDevice() []*MemoryDevice {
	list := []*MemoryDevice{}
	if rs, ok := s.Records[memoryDevice]; ok {
		for _, r := range rs {
			list = append(list, r.Data.(*MemoryDevice))
		}
	}
	return list
}

// GetPowerSupply returns PowerSuppl(y|ies)
func (s Spec) GetPowerSupply() []*PowerSupply {
	list := []*PowerSupply{}
	if rs, ok := s.Records[systemPowerSupply]; ok {
		for _, r := range rs {
			list = append(list, r.Data.(*PowerSupply))
		}
	}
	return list
}

// NewDecoder creates and initializes a Spec
func NewDecoder() *Spec {
	return new(Spec)
}

// Decode makes Spec satisfy the mox.Decoder interface
func (s *Spec) Decode() error {
	rc, ep, err := gosmbios.Stream()
	if err != nil {
		return fmt.Errorf("failed to open stream: %v", err)
	}
	defer rc.Close()

	d := gosmbios.NewDecoder(rc)
	tbls, err := d.Decode()
	if err != nil {
		return fmt.Errorf("failed to decode structures: %v", err)
	}

	s.Major, s.Minor, s.Rev = ep.Version()
	s.Records = make(map[uint8][]*Structure)

	log.Debugf("found smbios v%s", s.Version())

	for _, tbl := range tbls {
		st := new(Structure)

		st.Header = new(Header)
		st.Header.Type = tbl.Header.Type
		st.Header.Handle = tbl.Header.Handle

		log.Debugf("type id: %x, handle id: %x", st.Header.Type, st.Header.Handle)

		switch st.Header.Type {
		case biosInformation:
			log.Debug("type: bios_information")
			st.Data = parseBIOS(tbl)
		case systemInformation:
			log.Debug("type: system_information")
			st.Data = parseSystem(tbl)
		case baseboardInformation:
			log.Debug("type: baseboard_information")
			st.Data = parseBaseboard(tbl)
		case systemEnclosure:
			log.Debug("type: system_enclosure")
			st.Data = parseChassis(tbl)
		case processorInformation:
			log.Debug("type: processor_information")
			st.Data = parseProcessor(tbl)
		case cacheInformation:
			log.Debug("type: cache_information")
			st.Data = parseCache(tbl)
		case memoryDevice:
			log.Debug("type: memory_device")
			st.Data = parseMemoryDevice(tbl)
		case systemPowerSupply:
			log.Debug("type: system_power_supply")
			st.Data = parsePowerSupply(tbl)
		default:
			log.Debug("non-supported record type")
			continue
		}

		if st.Data != nil {
			s.Records[st.Header.Type] = append(s.Records[st.Header.Type], st)
		} else {
			log.Debug("could not decode record")
		}
	}

	return nil
}
