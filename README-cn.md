## 介绍

`Image2webp`项目是使用golang编写的支持多种图片格式实时转码成`Webp`格式的图片的应用程序。

目前支持`JPG`,`JPEG`, `PNG`, `BMP`, `GIF` 这些格式转码成`Webp`格式。

`Webp`格式能够在不影响图片质量的情况下显著的缩小图片体积，从而提升网络传输的速度，节省带宽流量资源。

本项目还支持`Webp`格式的缩略图功能，也是实时转码，能够无侵入的配置到现有服务器中，只需要在`nginx`中做一些转发，不用更改任何其他逻辑。

## 最佳实践

本项目的最佳实现是在图片服务器中启动本服务，在nginx或者其他web服务器中配置转发，
只需要匹配`^(.+)_webp(_(\d+)_(\d+)(.*))?$` 这个url格式就能转发到image2webp程序中,从而输出webp格式的图片。

如果是前端访问图片，在请求量很大的时候，这种实时转码的方式会造成cpu性能不足，所以针对大请求量的情况，
可以选择开启缓存，设置缓存目录，生成的webp文件会写入该目录。
把cdn图片缓存有效期设置的长一点，本服务只是回源使用，那样就很完美了。

## 编译

编译是非常简单的，确保你的go版本是`1.16`以上，并且打开了`mod`的支持。

执行`make` 命令就自动编译好了.
当然你也可以手动执行全路径编译命令
```shell script
$ go build -v -o image2webp
```
编译成功后，你可以执行 
```shell script
$ ./image2webp start --config=./config.json
```
启动它

## 配置

`config.json `中是json格式的配置.

* addr 监听的地址与端口 "127.0.0.1:8563"
* storagePath 图片存储的目录
* quality 转码成webp时的图片质量百分比 默认是80%
* allowSizes  支持的缩略图格式,_50_50: _width_height
* enableTmpFile  开启图片缓存
* tmpStoragePath  图片缓存的目录

## 部署

`image2webp` 部署也是很方便的，只需要把编译的二进制文件直接放到服务器上去运行，
再配置一下前端代理，具体可以参考`nginx.conf`。

部署好后的效果如下：

* 原图片地址: `http://image.example.com/images/timg.jpeg`
* webp格式图片地址: `http://image.example.com/images/timg.jpeg_webp`
* webp缩略图格式地址: `http://image.example.com/images/timg.jpeg_webp_100_100`

## 感谢

* 本项目是基于[DMicro](https://gitee.com/osgochina/dmicro)框架开发.
* 本项目是基于[Go Frame](https://github.com/gogf/gf)框架开发.
* 参考项目[webp_server_go](https://github.com/webp-sh/webp_server_go)
* 感谢提供了各种图片转换库的项目，大家可以看源码库的引用。





