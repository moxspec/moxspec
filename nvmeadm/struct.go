package nvmeadm

type nvmeAdminCmd struct {
	opcode      uint8
	flags       uint8
	rsvd1       uint16
	nsid        uint32
	cdw2        uint32
	cdw3        uint32
	metadata    uint64
	addr        uint64
	metadataLen uint32
	dataLen     uint32
	cdw10       uint32
	cdw11       uint32
	cdw12       uint32
	cdw13       uint32
	cdw14       uint32
	cdw15       uint32
	timeoutMs   uint32
	result      uint32
}

// Device represents NVMe disk
type Device struct {
	path                string
	class               byte
	Size                uint64
	CurTemp             int16
	WarnTemp            int16
	CritTemp            int16
	SpareSpace          byte
	Used                byte
	SerialNumber        string
	ModelNumber         string
	FirmwareRevision    string
	UnitsRead           uint64
	UnitsWritten        uint64
	ByteRead            uint64
	ByteWritten         uint64
	PowerCycleCount     uint64
	PowerOnHours        uint64
	UnsafeShutdownCount uint64
	MaxNamespaces       uint32
	NamespaceSizes      []uint64
	CritWarnings        []string
	UnrecoveredError    uint64 // Media and Data Integrity Errors
}

// GetNamespaceSize returns a total bytes of specified namespace
func (d Device) GetNamespaceSize(id int) uint64 { // id starts from 1
	if id <= 0 {
		return 0
	}

	if id > len(d.NamespaceSizes) {
		return 0
	}

	return d.NamespaceSizes[id-1]
}
