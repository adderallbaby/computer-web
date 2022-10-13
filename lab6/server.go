package main

import (
	"flag"
	"log"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

func main() {

	root := flag.String("root", "/Users/dmikuivashev/Downloads/vlab6/server", "Root directory to serve")
	// на предыдущей строке вторым приколом путь на папку, которая в роли сервера выступает
	user := flag.String("user", "dmikuivashev", "Username for login")
	pass := flag.String("pass", "1234", "Password for login")
	port := flag.Int("port", 2121, "Port")
	host := flag.String("host", "localhost", "Host")

	flag.Parse()
	if *root == "" {
		log.Fatalf("Please set a root to serve with -root")
	}

	factory := &filedriver.FileDriverFactory{
		RootPath: *root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", *user, *pass)
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
