package netlink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"syscall"
)

const (
	iflaStats64 = 23
)

func getStats(nli *netlinkInterface, ifname string) (*RtnlLinkStats64, error) {
	req := statsRequest()
	err := nli.post(req)
	if err != nil {
		return nil, err
	}

	var nlms []syscall.NetlinkMessage

	for {
		buf, err := nli.recv()
		if err != nil {
			return nil, err
		}

		nlmsTemp, err := syscall.ParseNetlinkMessage(buf)
		if err != nil {
			log.Debugf("netlink message parse failed: %s", err)
			return nil, err
		}
		nlms = append(nlms, nlmsTemp...)

		nlmDone := false
		for _, nlm := range nlmsTemp {
			log.Debugf("NLM header: 0x%0X", nlm.Header.Type)
			if nlm.Header.Type == syscall.NLMSG_DONE {
				log.Debug("NLMSG_DONE found")
				nlmDone = true
				break
			}
		}

		if nlmDone {
			break
		}
	}

	stats, err := parseStats(nlms, ifname)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func statsRequest() []byte {
	req := struct {
		nh  syscall.NlMsghdr
		ifi syscall.IfInfomsg
	}{
		syscall.NlMsghdr{
			Len:   syscall.SizeofNlMsghdr + syscall.SizeofIfInfomsg,
			Type:  syscall.RTM_GETLINK,
			Flags: syscall.NLM_F_REQUEST | syscall.NLM_F_ROOT,
			Seq:   0,
			Pid:   0,
		},
		syscall.IfInfomsg{
			Family:     syscall.AF_UNSPEC,
			X__ifi_pad: 0,
			Type:       0,
			Index:      0,
			Flags:      0,
			Change:     0,
		},
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &req)

	return buf.Bytes()
}

func parseStats(nlms []syscall.NetlinkMessage, target string) (*RtnlLinkStats64, error) {
	attrs := getTargetAttrs(nlms, target)
	if attrs == nil {
		return nil, fmt.Errorf("target interface is not found (%s)", target)
	}
	for _, attr := range attrs {
		if attr.Attr.Type == iflaStats64 {
			return parseRtnlLinkStats64(attr.Value)
		}
	}

	return nil, fmt.Errorf("iflaStats64 is not found")
}

func getTargetAttrs(nlms []syscall.NetlinkMessage, target string) []syscall.NetlinkRouteAttr {
	for _, nlm := range nlms {
		if nlm.Header.Type != syscall.RTM_NEWLINK {
			continue
		}

		attrs, err := syscall.ParseNetlinkRouteAttr(&nlm)
		if err != nil {
			log.Debugf("failed: %s", err)
			return nil
		}

		for _, attr := range attrs {
			if attr.Attr.Type == syscall.IFLA_IFNAME {
				// remove null charactor ('\0')
				ifname := strings.Trim(string(attr.Value), "\x00")
				log.Debugf("ifname: %s", ifname)

				if ifname == target {
					return attrs
				}
			}
		}
	}

	return nil
}

func parseRtnlLinkStats64(val []byte) (*RtnlLinkStats64, error) {
	stats := new(RtnlLinkStats64)
	reader := bytes.NewReader(val)
	err := binary.Read(reader, binary.LittleEndian, stats)

	if err != nil {
		log.Debugf("failed: %s", err)
		return nil, err
	}

	return stats, nil
}
