package scsi

type pg83Descriptor struct {
	CodeSet   byte
	PID       byte
	DesigType byte
	Assoc     byte
	PIV       byte
	Len       int
	Body      []byte
}
