## Introduce
[中文文档](https://github.com/osgochina/image2webp/blob/master/README-cn.md)

The Image2webp project is an application written with golang that supports real-time transcoding of multiple image formats into webp images.

Currently support `JPG`, `JPEG`, `PNG`, `BMP`, `GIF` these formats into `Webp` format.

'Webp' format can significantly reduce the size of the picture without affecting the picture quality, so as to improve the speed of network transmission, save bandwidth traffic resources.

This project also supports webp format thumbnail function, is also real-time transcoding, can be configured to the existing server without invasion, only need to do some forwarding in nginx, do not change any other logic.

## Best practice

The best realization of this project is to start this service in the picture server and configure forwarding in nginx or other web servers,
Just match the url format `^(.+)_webp(_(\d+)_(\d+)(.*))?` to forward to the image2webp program and output the webp image.

If it is the front end to access the picture, when the request volume is large, this real-time transcoding method will cause CPU performance, so for the large request volume, please see, before you must use the CDN file.

Set the CDN image cache expiration date a little longer, the service is only used back to the source, that is perfect.

## Compile

Compilation is very simple, make sure your go version is `1.14` and above, and turn on support for `mod`.

Execute the `make` command and it will compile automatically.
Of course, you can also manually execute the full path compiler command.
```shell script
 go build -v -o builds/image2webp
```
After compiling successfully, you can execute
```shell script
./builds/image2webp -f config.json
```
To start it.

## Configuration

`config.json `Is in json format.

* addr The listening address and port. link "127.0.0.1:8563"
* storagePath The directory where images are stored
* quality Image quality when transcoding to webp.  The default percentage is 80%
* allowSizes  Supported thumbnail format. _50_50: _width_height

## Deploy

`image2webp` deployment is also very convenient, just need to put the compiled binaries directly on the server to run,
Then configure the front - end proxy, you can refer to `nginx.conf`.

### Thanks

* This project is based on the [Go Frame](https://github.com/gogf/gf) framework.
* Refer to the project [webp_server_go](https://github.com/webp-sh/webp_server_go)
* Thank you for providing a variety of image conversion library project, you can see the source library reference.





