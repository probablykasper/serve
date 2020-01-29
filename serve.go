package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

var helpText = `
Serve {{.Version}}
  
USAGE:
    serve [dir] [options...]

EXAMPLE:
    serve -p 80 ./website

OPTIONS:
    {{range .VisibleFlags}}{{.}}
    {{end}}
`

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = helpText
	app.Version = "1.0.0"
	app.HideHelp = true // so `serve help` works
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address, a",
			Usage: "The IP address or hostname of the interface",
			Value: "0.0.0.0",
		},
		cli.IntFlag{
			Name:  "port, p",
			Usage: "The port to listen on",
			Value: 2233,
		},
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "Log requests",
		},
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "Show help menu",
		},
	}
	app.Action = func(c *cli.Context) error {

		addr := c.String("address")
		port := c.Int("port")
		dir := c.Args().Get(0)
		if len(c.Args()) == 0 {
			dir = "./"
		}
		verbose := c.Bool("verbose")
		if c.Bool("help") {
			cli.ShowAppHelpAndExit(c, 0)
		}

		loggingHandler := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if verbose {
					log.Println(r.RemoteAddr, r.Method, r.URL.Path, r.URL.RawQuery)
				}
				// When you press reload in Chrome, the page may be cached. Setting
				// Cache-Control no-store fixes this, but with the side-effect of
				// disabling cache for the back button. It's probably better to keep
				// this commented out for that reason?
				// w.Header().Set("Cache-Control", "no-store")
				w.Header().Set("Cache-Control", "no-cache")
				h.ServeHTTP(w, r)
			})
		}

		http.Handle("/", loggingHandler(http.FileServer(http.Dir(dir))))
		fmt.Println("")
		if addr == "0.0.0.0" {
			externalAddr, err := externalIp()
			if err != nil {
				externalAddr = addr
			}
			fmt.Printf("  Serving %s at:\n", dir)
			fmt.Printf("  - Local:   http://%s:%d\n", addr, port)
			fmt.Printf("  - Network: http://%s:%d\n", externalAddr, port)
		} else {
			fmt.Printf("  Serving %s at %s:%d\n", dir, addr, port)
		}
		fmt.Println("")
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}

	app.Run(os.Args)
}

func externalIp() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", errors.New("Could not get local ip")
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", errors.New("Could not get local ip")
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("Could not get local ip")
}
