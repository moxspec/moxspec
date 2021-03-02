package netlink

import (
	"fmt"
	"syscall"
)

const (
	bufferSize = 32768
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

func (nli *netlinkInterface) post(req []byte) error {
	sa := syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pad:    0,
		Pid:    0,
		Groups: 0,
	}

	return syscall.Sendto(nli.fd, req, 0, &sa)
}

func (nli *netlinkInterface) recv() ([]byte, error) {
	ptr := 0
	buf := make([]byte, bufferSize)
	for {
		n, _, err := syscall.Recvfrom(nli.fd, nil, syscall.MSG_PEEK|syscall.MSG_TRUNC)
		if err != nil {
			log.Debugf("scan failed: %s", err)
			return nil, err
		}

		if n == 0 {
			return nil, fmt.Errorf("end of file on the socket")
		}

		if n > bufferSize {
			return nil, fmt.Errorf("not enough buffer (%d > %d)", n, bufferSize)
		}

		log.Debugf("%d bytes in the queue", n)

		n, _, err = syscall.Recvfrom(nli.fd, buf, 0)
		if err != nil {
			log.Debugf("recv failed: %s", err)
			return nil, err
		}
		log.Debugf("read %d bytes", n)
		ptr += n
		if n > 0 {
			break
		}
	}
	return buf[:ptr], nil
}

func (nli *netlinkInterface) close() error {
	log.Debugf("close socket (fd = %d)", nli.fd)
	return syscall.Close(nli.fd)
}
