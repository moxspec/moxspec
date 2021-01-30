package spc

import (
	"math"
)

func makeReadLogSenseCmd(page, subpage, length byte) []byte {
	// SCSI Commands Reference Manual
	// 3.8 LOG SENSE command
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	return []byte{
		0x4D,
		0x00,
		page,
		subpage,
		0x00,
		0x00,
		0x00,
		0x00,
		length,
		0x00,
	}
}

func makeInquiryCmd(page byte, evpd bool, length int) []byte {
	if length > math.MaxUint16 || length < 0 {
		log.Debugf("allocate length of inquiry cmd is invalid (%d)", length)
		return nil
	}

	cdb02h := byte(0x00)
	if evpd {
		cdb02h = byte(0x01)
	}

	// SCSI Commands Reference Manual
	// 3.6.1 INQUIRY command introduction
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	return []byte{
		0x12,
		cdb02h,
		page,
		byte(length & 0xFF00 >> 8),
		byte(length & 0x00FF),
		0x00,
	}
}

func makeATAPThruCmd(feature, lbaLow, lbaMid, lbaHigh, cmd byte) []byte {
	// ATA Command Pass-Through
	// 13.2.3 ATA PASS-THROUGH (16) command overview
	// http://ftp.t10.org/ftp/t10/document.04/04-262r8.pdf
	return []byte{
		0x85, // Opecode
		0x08, // Protocul
		0x0E, // Flags
		0x00,
		feature, // feature
		0x00,
		0x01, // sector cound
		0x00,
		lbaLow, //lba low
		0x00,
		lbaMid, // lba mid
		0x00,
		lbaHigh, // lba high
		0x00,    //device
		cmd,     // command
		0x00,    // controll
	}
}
