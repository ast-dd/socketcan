package socketcan_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/ast-dd/socketcan"
)

type cErr struct {
	class    uint
	partS    string
	detail   uint
	location uint
}

var parseCanErrorsTests = []struct {
	name     string
	class    int
	data     []byte
	wantErrs []cErr
}{
	{"no errors", 0x00000002, []byte{0, 0, 0, 0, 0, 0, 0}, []cErr{}},
	{"wrong dlc", 0x20000002, []byte{0, 0, 0, 0, 0, 0, 0}, []cErr{
		{0, "length", 7, 0}},
	},
	{"arbitration: unspecified", 0x20000002, []byte{0, 0, 0, 0, 0, 0, 0, 0}, []cErr{
		{0x2, "lost arbitration: unspecified", 0, 0},
	}},
	{"arbitration: @17", 0x20000002, []byte{17, 0, 0, 0, 0, 0, 0, 0}, []cErr{
		{0x2, "lost arbitration: bit 17", 0, 17},
	}},
	{"ctrl: warning", 0x20000004, []byte{0, 0x0c, 0, 0, 0, 0, 0, 0}, []cErr{
		{0x4, "controller problems: reached warning", 0x4, 0},
		{0x4, "warning level for TX", 0x8, 0},
	}},
	{"ctrl: passive", 0x20000004, []byte{0, 0x30, 0, 0, 0, 0, 0, 0}, []cErr{
		{0x4, "controller problems: reached error passive status RX", 0x10, 0},
		{0x4, "passive status TX", 0x20, 0},
	}},
	{"prot: stuff@dlc", 0x20000008, []byte{0, 0, 0x04, 0x0b, 0, 0, 0, 0}, []cErr{
		{0x8, "protocol violations: bit stuffing error @ data length code", 0x4, 0xb},
	}},
	{"trx: no canl wire", 0x20000010, []byte{0, 0, 0, 0, 0x40, 0, 0, 0}, []cErr{
		{0x10, "transceiver status: CANL no wire", 0x40, 0},
	}},
	{"trx: CANH_SHORT_TO_GND", 0x20000010, []byte{0, 0, 0, 0, 0x07, 0, 0, 0}, []cErr{
		{0x10, "transceiver status: CANH short to GND", 0x7, 0},
	}},
	{"bus off", 0x20000040, []byte{0, 0, 0, 0, 0, 0, 0, 0}, []cErr{
		{0x40, "bus off", 0, 0},
	}},
	{"combination: tx timeout & prot overload", 0x20000009, []byte{0, 0, 0x20, 0x19, 0, 0, 0, 0}, []cErr{
		{0x1, "TX timeout", 0, 0},
		{0x8, "protocol violations: bus overload @ ACK slot", 0x20, 0x19},
	}},
}

func TestParseCanErrors(t *testing.T) {
	for _, tt := range parseCanErrorsTests {
		t.Run(tt.name, func(t *testing.T) {
			gotErrs := socketcan.ParseCanErrors(tt.class, tt.data)
			if g, w := len(gotErrs), len(tt.wantErrs); g != w {
				t.Errorf("ParseCanErrors() returned %v errors, want %v", g, w)
				return
			}

			// ensure determined order
			sort.Slice(gotErrs, func(i int, j int) bool { return gotErrs[i].Error() < gotErrs[j].Error() })

			for i, want := range tt.wantErrs {
				got := gotErrs[i]
				if got.Class != want.class {
					t.Errorf("ParseCanErrors() Class = %v, want %v", got.Class, want.class)
				}
				if got.Detail != want.detail {
					t.Errorf("ParseCanErrors() Detail = %v, want %v", got.Detail, want.detail)
				}
				if got.Location != want.location {
					t.Errorf("ParseCanErrors() Location = %v, want %v", got.Location, want.location)
				}
				if !strings.Contains(got.Error(), want.partS) {
					t.Errorf("ParseCanErrors() Error() = %v, want contain \"%v\"", got.Error(), want.partS)
				}
			}
		})
	}
}
