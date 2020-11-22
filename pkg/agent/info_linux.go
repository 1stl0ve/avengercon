package agent

import (
	"fmt"

	"github.com/1stl0ve/avengercon/pkg/api"
)

// linux specific systemInfo()
func (a *Agent) systemInfo(reg *api.Registration) {
	var err error
	reg.OS = "linux"

	reg.IP, err = getIPAddresses()
	if err != nil {
		fmt.Println(err)
	}
}
