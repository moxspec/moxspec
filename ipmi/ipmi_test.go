package ipmi

import (
	"fmt"
	"testing"
)

func TestParseLanPrint(t *testing.T) {
	var ipaddr, netmask, gateway, mac string
	tests := []struct {
		in      string
		ipaddr  string
		netmask string
		gateway string
		mac     string
	}{
		{`Set in Progress         : Set Complete
IP Address Source       : DHCP Address
IP Address              : 10.17.7.36
Subnet Mask             : 255.0.0.0
MAC Address             : 5c:54:6d:0b:54:6b
SNMP Community String   : TrapAdmin12#$
IP Header               : TTL=0x40 Flags=0x40 Precedence=0x00 TOS=0x10
Default Gateway IP      : 0.0.0.0
802.1q VLAN ID          : Disabled
RMCP+ Cipher Suites     : 0,1,2,3,17
Cipher Suite Priv Max   : XuuaXXXXXXXXXXX
                        :     X=Cipher Suite Unused
                        :     c=CALLBACK
                        :     u=USER
                        :     o=OPERATOR
                        :     a=ADMIN
                        :     O=OEM`, "10.17.7.36", "255.0.0.0", "0.0.0.0", "5c:54:6d:0b:54:6b"},
		{`garbage : garbage 
dummy : dummy 
mox : mox 
248 : 248 
fail : fail`, "", "", "", ""},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			ipaddr, netmask, gateway, mac = parseLanPrint(tt.in)
			if ipaddr != tt.ipaddr || netmask != tt.netmask || gateway != tt.gateway || mac != tt.mac {
				fmtr := "got: ipaddr %s / netmask %s / gw: %s / mac %s, expect: ipaddr %s / netmask %s / gw %s / mac %s"
				t.Errorf(fmtr, ipaddr, netmask, gateway, mac, tt.ipaddr, tt.netmask, tt.gateway, tt.mac)
			}
		})
	}
}
