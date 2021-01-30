package platform

import "github.com/actapio/moxspec/loglet"

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("platform")
}
