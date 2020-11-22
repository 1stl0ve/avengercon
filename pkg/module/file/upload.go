package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (m *Admin) setUploadConfig(args []string) error {
	var err error

	m.Config = config{
		Name:      UploadModuleName,
		Operation: upload,
		Local:     args[1],
		Remote:    args[2],
	}
	m.Config.Content, err = ioutil.ReadFile(m.Config.Local)
	if err != nil {
		return err
	}
	return nil
}

// Do prints out the results of the command that would've been written to stdout
// or stderr on the remote system.
func (m *Admin) upload(agent string) {
	fmt.Printf("\n[%s:%s:%s->%s]\n",
		agent,
		m.Name(),
		m.Config.Local,
		m.Config.Remote)
}

func (m *Agent) upload() error {
	err := ioutil.WriteFile(m.Config.Remote, m.Config.Content, os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}
