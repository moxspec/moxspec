package netlink

import (
	"syscall"
)

const (
	bufferSize = 4096
)

type netlinkInterface struct {
	fd int
}

func newNetlinkInterface() (*netlinkInterface, error) {
	log.Debug("opening netlink socket")
	sock, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_DGRAM, syscall.NETLINK_ROUTE)
	if err != nil {
		log.Debugf("failed: %s", err)
		return nil, err
	}

	log.Debugf("fd = %d", sock)
	nli := netlinkInterface{sock}

	return &nli, nil
}

func (nli *netlinkInterface) post(req []byte) ([]syscall.NetlinkMessage, error) {
	sa := syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pad:    0,
		Pid:    0,
		Groups: 0,
	}

	err := syscall.Sendto(nli.fd, req, 0, &sa)
	if err != nil {
		log.Debugf("failed: %s", err)
		return nil, err
	}

	data := []byte{}
	buf := [bufferSize]byte{}
	for {
		n, err := syscall.Read(nli.fd, buf[:])
		if err != nil {
			log.Debugf("failed: %s", err)
			return nil, err
		}

		data = append(data, buf[:n]...)
		if n != bufferSize {
			break
		}
	}

	if err != nil {
		log.Debugf("failed: %s", err)
		return nil, err
	}

	nlms, err := syscall.ParseNetlinkMessage(data)
	if err != nil {
		log.Debugf("failed: %s", err)
		return nil, err
	}

	return nlms, nil
}

func (nli *netlinkInterface) close() error {
	log.Debugf("close socket (fd = %d)", nli.fd)
	return syscall.Close(nli.fd)
}
