/**
* sync dir
* author: hidu <duv123+git@gmail.com>
* https://github.com/hidu/hsync
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	hsync "github.com/ketor/hsync/internal"
)

var d = flag.Bool("d", false, "run model,defaul is client")
var host = flag.String("h", "", "sync host name")
var ve = flag.Bool("version", false, "show version:"+hsync.GetVersion())
var demoConf = flag.String("demo_conf", "", "show default conf [client|server]")
var deployOnly = flag.Bool("deploy", false, "deploy all files for server.")

//add by ketor
var s = flag.String("s", "", "run simple server on addr:port")
var c = flag.String("c", "", "run simple client to addr:port")

func init() {
	flag.Lookup("alsologtostderr").DefValue = "true"
	flag.Set("alsologtostderr", "true")

	df := flag.Usage
	flag.Usage = func() {
		df()
		fmt.Fprintln(os.Stderr, "\n  sync dir, https://github.com/ketor/hsync/")
		fmt.Fprintln(os.Stderr, "  as client:", os.Args[0], "   [hsync.json]")
		fmt.Fprintln(os.Stderr, "  as server:", os.Args[0], "-d [hsyncd.json]")
		fmt.Fprintln(os.Stderr, "\nSimple Usage:")
		fmt.Fprintln(os.Stderr, "    Simple Client: hsync -c 127.0.0.1:8700")
		fmt.Fprintln(os.Stderr, "    Simple Server: hsync -s :8700")
	}
}

func main() {
	flag.Parse()
	if *ve {
		fmt.Fprintln(os.Stderr, "version:", hsync.GetVersion())
		os.Exit(0)
	}
	if *demoConf != "" {
		fmt.Println(hsync.DemoConf(*demoConf))
		os.Exit(0)
	}

	if *deployOnly {
		*d = true
	}

	confName := flag.Arg(0)
	if confName == "" {
		if *d {
			confName = "hsyncd.json"
		} else {
			confName = "hsync.json"
		}
	}

	//run simple server
	if *s != "" {
		glog.Info("New simple server start")
		server, err := hsync.NewSimpleHsyncServer(*s)
		if err != nil {
			glog.Exitln("start simple server failed:", err)
		}
		if *deployOnly {
			server.DeployAll()
			return
		}

		server.Start()
	} else if *c != "" {
		glog.Info("New simple client start")
		client, err := hsync.NewSimpleHsyncClient(*c)
		if err != nil {
			glog.Exitln("start hsync client failed:", err)
		}
		client.Connect()
		client.Watch()
	} else {
		confInfo, err := os.Stat(confName)
		if err != nil || confInfo.IsDir() {
			glog.Exitf("hsync conf [%s] not exists!", confName)
		}

		if *d {
			server, err := hsync.NewHsyncServer(confName)
			if err != nil {
				glog.Exitln("start server failed:", err)
			}
			if *deployOnly {
				server.DeployAll()
				return
			}

			server.Start()
		} else {
			client, err := hsync.NewHsyncClient(confName, *host)
			if err != nil {
				glog.Exitln("start hsync client failed:", err)
			}
			client.Connect()
			client.Watch()
		}
	}
}
