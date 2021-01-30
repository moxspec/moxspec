package main

import (
	"github.com/actapio/moxspec/model"
	"github.com/actapio/moxspec/platform"
)

func shapeMisc(r *model.Report) {
	sysname, release, hostname, _ := platform.Uname()

	os := new(model.OS)
	os.Distro = platform.GetDistroName()
	if os.Distro == "" {
		os.Distro = sysname
	}
	os.Kernel = release

	r.OS = os
	r.Hostname = hostname
}
