package platform

import "github.com/moxspec/moxspec/loglet"

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("platform")
}
