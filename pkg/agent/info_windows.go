package agent

import (
	"fmt"

	"github.com/1stl0ve/avengercon/pkg/api"
	"golang.org/x/sys/windows"
)

// windows specific systemInfo()
func (a *Agent) systemInfo(reg *api.Registration) {
	var err error
	reg.OS = "windows"
	major, minor, build := windows.RtlGetNtVersionNumbers()
	reg.Version = fmt.Sprintf("%d.%d.%d", major, minor, build)
	reg.IP, err = getIPAddresses()
	if err != nil {
		fmt.Println(err)
	}
}
