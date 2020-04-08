package main

import (
	"errors"
	"fmt"
	"github.com/gogf/gf-cli/library/mlog"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"reflect"
	"strconv"
)

const (
	VERSION = "v0.0.1"
)

func main() {
	for k := range gcmd.GetOptAll() {
		switch k {
		case "?", "h":
			fmt.Println("image2web -f  </path/to/config.json>")
			return
		case "i", "v":
			fmt.Println("image2web:" + VERSION)
			return
		}
	}
	configPath := gcmd.GetOpt("f")
	if configPath == "" {
		configPath = "./config.json"
	}
	if !gfile.Exists(configPath) {
		mlog.Print("image2web -f  </path/to/config.json>")
		return
	}
	g.Cfg().AddPath(gfile.Dir(configPath))
	g.Cfg().SetFileName(gfile.Basename(configPath))
	storagePath := g.Cfg().Get("storagePath")
	if storagePath == nil {
		fmt.Println("please set config storage path!")
		return
	}
	if !gfile.IsDir(storagePath.(string)) {
		fmt.Printf("The storage path [%s] does not exist!", storagePath)
		return
	}
	Addr := g.Cfg().GetString("addr", "127.0.0.1:8563")
	s := g.Server()
	s.BindHandler("/*", func(req *ghttp.Request) {
		url := req.Request.RequestURI
		quality := g.Cfg().GetFloat32("quality", 80)
		match, _ := gregex.MatchString(`^(.+)_webp(_(\d+)_(\d+)(.*))?$`, url)
		if len(match) <= 0 {
			req.Response.WriteHeader(404)
			req.Exit()
			return
		}
		var thump = false
		if len(match[2]) > 0 {
			//判断缩放是否在可用列表中
			allow_sizes := g.Cfg().GetArray("allowSizes")
			if b, _ := Contain(match[2], allow_sizes); !b {
				req.Response.WriteHeader(404)
				req.Exit()
			}
			thump = true
		}
		file := match[1]
		absFile := gfile.Abs(fmt.Sprintf("%s%s", storagePath, file))
		if !gfile.Exists(absFile) {
			req.Response.WriteHeader(404)
			req.Exit()
		}
		imageObj := Image{}
		err := imageObj.Open(absFile)
		if err != nil {
			req.Response.WriteHeader(404)
			req.Exit()
		}
		req.Header.Add("Content-Type", "image/webp")
		if !thump {
			buf, err := imageObj.ToWebP(quality)
			if err != nil {
				req.Response.WriteHeader(404)
				req.Exit()
			}
			req.Header.Add("Content-Length", string(len(buf)))
			req.Response.WriteExit(buf)
			return
		}
		width, _ := strconv.Atoi(match[3])
		height, _ := strconv.Atoi(match[4])
		buf, err := imageObj.MakeThumbnail(width, height)
		if err != nil {
			req.Response.WriteHeader(404)
			req.Exit()
		}
		req.Header.Add("Content-Length", string(len(buf)))
		req.Response.WriteExit(buf)
	})
	s.SetAddr(Addr)
	s.Run()
}

func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in array")
}
