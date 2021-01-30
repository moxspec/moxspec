package nw

import (
	"fmt"
	"net"
	"testing"
)

func TestCalcIPv4Broadcast(t *testing.T) {
	var got string
	tests := []struct {
		in string
		ex string
	}{
		{"10.0.0.0/1", "127.255.255.255"},
		{"172.248.22.23/21", "172.248.23.255"},
		{"248.248.248.0/24", "248.248.248.255"},
		{"248.248.248.248/32", ""},
		{"248.248.248.248/0", ""},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			addr, anet, _ := net.ParseCIDR(tt.in)
			got, _ = calcIPv4Broadcast(addr, anet)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}

}
