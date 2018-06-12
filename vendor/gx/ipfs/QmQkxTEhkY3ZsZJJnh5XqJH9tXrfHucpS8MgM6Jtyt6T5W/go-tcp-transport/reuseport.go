package tcp

import (
	"os"
	"strings"

	reuseport "gx/ipfs/QmYkx3VWoYeQJKEpZp13MxQYkWtLd9E6G4fqyyDCS76qcD/go-reuseport"
)

// envReuseport is the env variable name used to turn off reuse port.
// It default to true.
const envReuseport = "IPFS_REUSEPORT"

// envReuseportVal stores the value of envReuseport. defaults to true.
var envReuseportVal = true

func init() {
	v := strings.ToLower(os.Getenv(envReuseport))
	if v == "false" || v == "f" || v == "0" {
		envReuseportVal = false
		log.Infof("REUSEPORT disabled (IPFS_REUSEPORT=%s)", v)
	}
}

// reuseportIsAvailable returns whether reuseport is available to be used. This
// is here because we want to be able to turn reuseport on and off selectively.
// For now we use an ENV variable, as this handles our pressing need:
//
//   IPFS_REUSEPORT=false ipfs daemon
//
// If this becomes a sought after feature, we could add this to the config.
// In the end, reuseport is a stop-gap.
func ReuseportIsAvailable() bool {
	return envReuseportVal && reuseport.Available()
}