package main

import (
	"fmt"
	"log"
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
	app.Version = "0.0.2"
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
			Value: 8888,
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
		log.Printf("Serving %s at http://%s:%d/", dir, addr, port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}

	app.Run(os.Args)
}
