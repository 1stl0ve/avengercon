package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (m *Admin) setDownloadConfig(args []string) {
	m.Config = config{
		Name:      DownloadModuleName,
		Operation: download,
		Local:     args[2],
		Remote:    args[1],
	}
}

func (m *Admin) download(agent string, content []byte) error {
	fmt.Printf("\n[%s:%s:%s->%s]\n",
		agent,
		m.Name(),
		m.Config.Remote,
		m.Config.Local)
	err := ioutil.WriteFile(m.Config.Local, content, os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}

func (m *Agent) download() error {
	var err error
	if m.Config.Content, err = ioutil.ReadFile(m.Config.Remote); err != nil {
		return err
	}
	return nil
}
