package main

import (
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
	"image2webp/app"
	_ "net/http/pprof"
)

const (
	VERSION = "v1.0.1"
)

func main() {
	dserver.CloseCtl()
	dserver.Authors = "ClownFish"
	dserver.BuildVersion = VERSION
	dserver.SetName("image2webp")
	dserver.Setup(func(svr *dserver.DServer) {
		err := svr.AddSandBox(new(app.App))
		if err != nil {
			logger.Fatal(err)
		}
	})
}
