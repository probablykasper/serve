package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "serve"
	app.Version = "0.0.2"
	app.Usage = ""
	app.Description = `A simple HTTP server for serving static files.

   See https://github.com/sandeepraju/serve for more details!

EXAMPLES:
   * Serve files in the current directory

			serve

   * Serve at port 8000 on 192.168.0.3

	    serve -a 192.168.0.3 -p 8000

	 * Serve /tmp

			serve -d /tmp`
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "directory, dir, d",
			Usage: "The directory to serve",
			Value: "./",
		},
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
			Name:  "log-requests, verbose, V",
			Usage: "Don't print access logs",
		},
	}
	app.Action = func(c *cli.Context) error {

		addr := c.String("address")
		port := c.Int("port")
		dir := c.String("directory")
		verbose := c.Bool("verbose")

		var epoch = time.Unix(0, 0).Format(time.RFC1123)
		loggingHandler := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if verbose {
					log.Println(r.RemoteAddr, r.Method, r.URL.Path, r.URL.RawQuery)
				}
				w.Header().Set("Expires", epoch)
				w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("X-Accel-Expires", "0")
				w.Header().Del("ETag")
				w.Header().Del("If-Modified-Since")
				w.Header().Del("If-Match")
				w.Header().Del("If-None-Match")
				w.Header().Del("If-Range")
				w.Header().Del("If-Unmodified-Since")
				h.ServeHTTP(w, r)
			})
		}

		fmt.Println(http.Dir(dir))
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
