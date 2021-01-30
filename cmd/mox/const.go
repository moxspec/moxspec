package main

// Transport is used to indicate transport protocol
type Transport string

// Transports
const (
	SATA Transport = "SATA"
	SAS  Transport = "SAS"
	NVMe Transport = "NVMe"
)
