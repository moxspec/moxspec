package netlink

import (
	"context"
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

func (nli *netlinkInterface) recv(ctx context.Context) ([]byte, error) {
	bufCh := make(chan []byte)
	errCh := make(chan error)
	go func() {
		ptr := 0
		buf := make([]byte, bufferSize)
		for {
			n, _, err := syscall.Recvfrom(nli.fd, nil, syscall.MSG_PEEK|syscall.MSG_TRUNC)
			if err != nil {
				errCh <- fmt.Errorf("scan failed: %w", err)
				return
			}

			if n == 0 {
				errCh <- fmt.Errorf("end of file on the socket")
				return
			}

			if n > bufferSize {
				errCh <- fmt.Errorf("not enough buffer (%d > %d)", n, bufferSize)
				return
			}

			log.Debugf("%d bytes in the queue", n)

			n, _, err = syscall.Recvfrom(nli.fd, buf, 0)
			if err != nil {
				errCh <- fmt.Errorf("recv failed: %w", err)
				return
			}

			log.Debugf("read %d bytes", n)
			ptr += n
			if n > 0 {
				break
			}
		}
		bufCh <- buf[:ptr]
		return
	}()

	select {
	case buf := <-bufCh:
		return buf, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
	}

	return nil, fmt.Errorf("recv timeout")
}

func (nli *netlinkInterface) close() error {
	log.Debugf("close socket (fd = %d)", nli.fd)
	return syscall.Close(nli.fd)
}
