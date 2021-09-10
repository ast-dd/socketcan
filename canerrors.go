package socketcan

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// got from Linux /include/uapi/linux/can/error.h
// search regex: #define ([\w]+) +(0x[0-9A-F]+)U? \/\* (.+) \*\/
// replace regex maps: $1: "$3",

// error class (mask) in can_id
var errStringsID = map[uint]string{
	unix.CAN_ERR_TX_TIMEOUT: "TX timeout (by netdevice driver)",
	unix.CAN_ERR_LOSTARB:    "lost arbitration",    // see data[0]
	unix.CAN_ERR_CRTL:       "controller problems", // see data[1]
	unix.CAN_ERR_PROT:       "protocol violations", // see data[2..3]
	unix.CAN_ERR_TRX:        "transceiver status",  // see data[4]
	unix.CAN_ERR_ACK:        "received no ACK on transmission",
	unix.CAN_ERR_BUSOFF:     "bus off",
	unix.CAN_ERR_BUSERROR:   "bus error (may flood!)",
	unix.CAN_ERR_RESTARTED:  "controller restarted",
}

// arbitration lost in bit ... / data[0]
var errStringsLostarb = map[uint]string{
	unix.CAN_ERR_LOSTARB_UNSPEC: "unspecified",
	/* else bit number in bitstream */
}

// error status of CAN-controller / data[1]
var errStringsCtrl = map[uint]string{
	unix.CAN_ERR_CRTL_UNSPEC:      "unspecified",
	unix.CAN_ERR_CRTL_RX_OVERFLOW: "RX buffer overflow",
	unix.CAN_ERR_CRTL_TX_OVERFLOW: "TX buffer overflow",
	unix.CAN_ERR_CRTL_RX_WARNING:  "reached warning level for RX errors",
	unix.CAN_ERR_CRTL_TX_WARNING:  "reached warning level for TX errors",
	unix.CAN_ERR_CRTL_RX_PASSIVE:  "reached error passive status RX (at least one error counter exceeds the protocol-defined level of 127)",
	unix.CAN_ERR_CRTL_TX_PASSIVE:  "reached error passive status TX (at least one error counter exceeds the protocol-defined level of 127)",
	unix.CAN_ERR_CRTL_ACTIVE:      "recovered to error active state",
}

// error in CAN protocol (type) / data[2]
var errStringsProtType = map[uint]string{
	unix.CAN_ERR_PROT_UNSPEC:   "unspecified",
	unix.CAN_ERR_PROT_BIT:      "single bit error",
	unix.CAN_ERR_PROT_FORM:     "frame format error",
	unix.CAN_ERR_PROT_STUFF:    "bit stuffing error",
	unix.CAN_ERR_PROT_BIT0:     "unable to send dominant bit",
	unix.CAN_ERR_PROT_BIT1:     "unable to send recessive bit",
	unix.CAN_ERR_PROT_OVERLOAD: "bus overload",
	unix.CAN_ERR_PROT_ACTIVE:   "active error announcement",
	unix.CAN_ERR_PROT_TX:       "error occurred on transmission",
}

// error in CAN protocol (location) / data[3]
var errStringsProtLocation = map[uint]string{
	unix.CAN_ERR_PROT_LOC_UNSPEC:  "unspecified",
	unix.CAN_ERR_PROT_LOC_SOF:     "start of frame",
	unix.CAN_ERR_PROT_LOC_ID28_21: "ID bits 28 - 21 (SFF: 10 - 3)",
	unix.CAN_ERR_PROT_LOC_ID20_18: "ID bits 20 - 18 (SFF: 2 - 0)",
	unix.CAN_ERR_PROT_LOC_SRTR:    "substitute RTR (SFF: RTR)",
	unix.CAN_ERR_PROT_LOC_IDE:     "identifier extension",
	unix.CAN_ERR_PROT_LOC_ID17_13: "ID bits 17-13",
	unix.CAN_ERR_PROT_LOC_ID12_05: "ID bits 12-5",
	unix.CAN_ERR_PROT_LOC_ID04_00: "ID bits 4-0",
	unix.CAN_ERR_PROT_LOC_RTR:     "RTR",
	unix.CAN_ERR_PROT_LOC_RES1:    "reserved bit 1",
	unix.CAN_ERR_PROT_LOC_RES0:    "reserved bit 0",
	unix.CAN_ERR_PROT_LOC_DLC:     "data length code",
	unix.CAN_ERR_PROT_LOC_DATA:    "data section",
	unix.CAN_ERR_PROT_LOC_CRC_SEQ: "CRC sequence",
	unix.CAN_ERR_PROT_LOC_CRC_DEL: "CRC delimiter",
	unix.CAN_ERR_PROT_LOC_ACK:     "ACK slot",
	unix.CAN_ERR_PROT_LOC_ACK_DEL: "ACK delimiter",
	unix.CAN_ERR_PROT_LOC_EOF:     "end of frame",
	unix.CAN_ERR_PROT_LOC_INTERM:  "intermission",
}

// error status of CAN-transceiver / data[4]                    CANH CANL
var errStringsTrx = map[uint]string{
	unix.CAN_ERR_TRX_UNSPEC:             "unspecified",           // 0000 0000
	unix.CAN_ERR_TRX_CANH_NO_WIRE:       "CANH no wire",          // 0000 0100
	unix.CAN_ERR_TRX_CANH_SHORT_TO_BAT:  "CANH short to battery", // 0000 0101
	unix.CAN_ERR_TRX_CANH_SHORT_TO_VCC:  "CANH short to VCC",     // 0000 0110
	unix.CAN_ERR_TRX_CANH_SHORT_TO_GND:  "CANH short to GND",     // 0000 0111
	unix.CAN_ERR_TRX_CANL_NO_WIRE:       "CANL no wire",          // 0100 0000
	unix.CAN_ERR_TRX_CANL_SHORT_TO_BAT:  "CANL short to BAT",     // 0101 0000
	unix.CAN_ERR_TRX_CANL_SHORT_TO_VCC:  "CANL short to VCC",     // 0110 0000
	unix.CAN_ERR_TRX_CANL_SHORT_TO_GND:  "CANL short to GND",     // 0111 0000
	unix.CAN_ERR_TRX_CANL_SHORT_TO_CANH: "CANL short to CANH",    // 1000 0000
}

// CanError contains the error information parsed from CAN error frames
type CanError struct {
	s        string
	Class    uint
	Detail   uint
	Location uint
}

// Error implements the error interface
func (e CanError) Error() (s string) {
	return e.s
}

// ParseCanErrors returns the CAN errors from given CAN frame
func ParseCanErrors(id int, data []byte) (errs []CanError) {
	errs = make([]CanError, 0)
	uid := uint(id)

	if uid&unix.CAN_ERR_FLAG == 0 {
		return
	}

	if len(data) != unix.CAN_ERR_DLC {
		err := CanError{
			s:      fmt.Sprintf("error frame data invalid length: %d", len(data)),
			Class:  0,
			Detail: uint(len(data)),
		}
		errs = append(errs, err)
		return
	}

	for classMask, classS := range errStringsID {
		if uid&classMask == 0 {
			continue
		}

		switch classMask {
		case unix.CAN_ERR_LOSTARB:
			errs = append(errs, canErrLostarb(data)...)
		case unix.CAN_ERR_CRTL:
			errs = append(errs, canErrCrtl(data)...)
		case unix.CAN_ERR_PROT:
			errs = append(errs, canErrProt(data)...)
		case unix.CAN_ERR_TRX:
			errs = append(errs, canErrTrx(data)...)
		default:
			errs = append(errs, CanError{
				Class: classMask,
				s:     classS,
			})
		}
	}

	return
}

func canErrLostarb(data []byte) (errs []CanError) {
	errs = make([]CanError, 0)
	classMask := uint(unix.CAN_ERR_LOSTARB)
	var locationS string
	if data[0] == unix.CAN_ERR_LOSTARB_UNSPEC {
		locationS = errStringsLostarb[0]
	} else {
		locationS = fmt.Sprintf("bit %d", data[0])
	}
	errs = append(errs, CanError{
		s:        errStringsID[classMask] + ": " + locationS,
		Class:    classMask,
		Location: uint(data[0]),
	})
	return
}

func canErrCrtl(data []byte) (errs []CanError) {
	errs = make([]CanError, 0)
	classMask := uint(unix.CAN_ERR_CRTL)
	for ctrlBit, ctrlS := range errStringsCtrl {
		if byte(ctrlBit)&data[1] != 0 {
			errs = append(errs, CanError{
				s:      errStringsID[classMask] + ": " + ctrlS,
				Class:  classMask,
				Detail: ctrlBit,
			})
		}
	}
	return
}

func canErrProt(data []byte) (errs []CanError) {
	errs = make([]CanError, 0)
	classMask := uint(unix.CAN_ERR_PROT)
	for protMask, protS := range errStringsProtType {
		if byte(protMask)&data[2] != 0 {
			locationS := errStringsProtLocation[uint(data[3])]
			errs = append(errs, CanError{
				Class:    classMask,
				s:        errStringsID[classMask] + ": " + protS + " @ " + locationS,
				Detail:   protMask,
				Location: uint(data[3]),
			})
		}
	}
	return
}

func canErrTrx(data []byte) (errs []CanError) {
	errs = make([]CanError, 0)
	classMask := uint(unix.CAN_ERR_TRX)
	if data[4] == unix.CAN_ERR_TRX_UNSPEC {
		errs = append(errs, CanError{
			Class:  classMask,
			s:      errStringsID[classMask] + ": " + errStringsTrx[unix.CAN_ERR_TRX_UNSPEC],
			Detail: unix.CAN_ERR_TRX_UNSPEC,
		})
	}
	for statusMask, statusS := range errStringsTrx {
		if statusMask == unix.CAN_ERR_TRX_UNSPEC {
			continue
		}
		if byte(statusMask) == data[4]&0x07 {
			errs = append(errs, CanError{
				Class:  classMask,
				s:      errStringsID[classMask] + ": " + statusS,
				Detail: statusMask,
			})
		}
		if byte(statusMask) == data[4]&0x70 {
			errs = append(errs, CanError{
				Class:  classMask,
				s:      errStringsID[classMask] + ": " + statusS,
				Detail: statusMask,
			})
		}
	}
	return
}
