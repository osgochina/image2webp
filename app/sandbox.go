package app

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/dserver"
	"github.com/osgochina/dmicro/logger"
	"sync"
)

type App struct {
	dserver.ServiceSandbox
	svr            *ghttp.Server
	quality        float32
	allowSizes     *garray.StrArray
	enableTmpFile  bool
	storagePath    string
	tmpStoragePath string
	imagePool      *sync.Pool
}

func (that *App) Name() string {
	return "app"
}

func (that *App) Setup() error {
	storagePath := that.Config.GetString("storagePath")
	if len(storagePath) <= 0 {
		logger.Fatal("please set config storagePath!")
	}
	that.storagePath = gfile.Abs(storagePath)
	if !gfile.IsDir(that.storagePath) {
		logger.Fatalf("storagePath %s not exist!", that.storagePath)
	}
	that.enableTmpFile = that.Config.GetBool("enableTmpFile", false)
	if that.enableTmpFile {
		tmpStoragePath := that.Config.GetString("tmpStoragePath")
		if len(tmpStoragePath) <= 0 {
			logger.Fatal("please set config tmp tmpStoragePath!")
		}
		that.tmpStoragePath = gfile.Abs(tmpStoragePath)
		if !gfile.IsDir(that.tmpStoragePath) {
			err := gfile.Mkdir(that.tmpStoragePath)
			if err != nil {
				logger.Fatal(err)
			}
		}
	}

	addr := that.Config.GetString("addr", "127.0.0.1:8563")
	that.quality = that.Config.GetFloat32("quality", 100)
	that.allowSizes = garray.NewStrArrayFrom(that.Config.GetStrings("allow_sizes"))
	that.imagePool = &sync.Pool{
		// 默认的返回值设置，不写这个参数，默认是nil
		New: func() interface{} {
			return &Image{}
		},
	}
	logger.Infof("Image Storage %s", that.storagePath)
	logger.Infof("Tmp Image Storage %s", that.tmpStoragePath)
	that.svr = g.Server()
	that.svr.BindHandler("/*", that.HandHandler)
	that.svr.SetAddr(addr)
	return that.svr.Start()
}

func (that *App) Shutdown() error {
	return that.svr.Shutdown()
}

func (that *App) HandHandler(req *ghttp.Request) {
	url := req.Request.RequestURI
	req.Header.Add("Content-Type", "image/webp")
	var tmpFile string
	//是否开启缓存文件
	if that.enableTmpFile {
		tmpFile = fmt.Sprintf("%s%s", that.tmpStoragePath, url)
		if gfile.IsFile(tmpFile) {
			req.Response.ServeFile(tmpFile)
			return
		}
	}
	//未匹配则返回404
	if !IsMatch(url) {
		req.Response.WriteHeader(404)
		return
	}
	//更进一步匹配未成功
	match, err := NewMatchUrl(url)
	if err != nil {
		req.Response.WriteHeader(404)
		return
	}
	// 缩略图，需要匹配使用的长宽是否符合配置
	if match.IsThump && !match.AllowSizes(that.allowSizes) {
		req.Response.WriteHeader(404)
		return
	}
	// 获取原始图片的绝对路径，并判断图片是否存在
	absFile := match.GetRealFilePath(that.storagePath)
	if !gfile.Exists(absFile) {
		req.Response.WriteHeader(404)
		return
	}
	// 从内存池中获取内存
	imageObj := that.imagePool.Get().(*Image)
	defer that.imagePool.Put(imageObj)

	// 打开图片文件准备转换
	if err = imageObj.Open(absFile); err != nil {
		req.Response.WriteHeader(404)
		return
	}
	//如果要转换的图片不是缩略图
	if !match.IsThump {
		//小于96的png图片不要转换,没有解决边缘毛刺的问题
		if imageObj.ImageType == ImagePng && imageObj.Width <= 96 && imageObj.Height <= 96 {
			req.Header.Add("Content-Type", "image/png")
			req.Header.Add("Content-Length", gconv.String(len(imageObj.Data)))
			req.Response.Write(imageObj.Data)
			return
		}
		//转换图片
		buf, err := imageObj.ToWebP(that.quality)
		if err != nil {
			req.Response.WriteHeader(404)
			return
		}
		req.Header.Add("Content-Length", gconv.String(len(buf)))
		req.Response.Write(buf)
		if that.enableTmpFile {
			err = gfile.PutBytes(tmpFile, buf)
			if err != nil {
				glog.Error(err)
			}
		}
		return
	}
	//如果是要转换缩略图
	buf, err := imageObj.MakeThumbnail(match.Weight, match.Height, that.quality)
	if err != nil {
		req.Response.WriteHeader(404)
		glog.Warningf("%s error:%v", match.OriginUrl, err)
		return
	}
	req.Header.Add("Content-Length", gconv.String(len(buf)))
	req.Response.Write(buf)
	if that.enableTmpFile {
		err = gfile.PutBytes(tmpFile, buf)
		if err != nil {
			glog.Error(err)
		}
	}
	return
}
