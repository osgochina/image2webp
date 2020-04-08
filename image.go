package main

import (
	"bytes"
	"errors"
	"github.com/chai2010/webp"
	"github.com/gogf/gf/os/gfile"
	"github.com/nfnt/resize"
	giftowebp "github.com/sizeofint/gif-to-webp"
	"golang.org/x/image/bmp"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"strings"
)

const ImageJpeg int = 1
const ImagePng int = 2
const ImageBmp int = 3
const ImageGif int = 4
const DefaultMaxWidth float64 = 960
const DefaultMaxHeight float64 = 3000

type Image struct {
	FilePath  string
	Data      []byte
	ImageType int
	Ext       string
	Width     int
	Height    int
}

/**
打开图片文件
*/
func (that *Image) Open(filePath string) (err error) {
	that.FilePath = filePath
	that.Data = gfile.GetBytes(filePath)
	that.Ext = gfile.ExtName(filePath)
	contentType := http.DetectContentType(that.Data[:512])
	if strings.Contains(contentType, "jpeg") {
		that.ImageType = ImageJpeg
	} else if strings.Contains(contentType, "png") {
		that.ImageType = ImagePng
	} else if strings.Contains(contentType, "bmp") {
		that.ImageType = ImageBmp
	} else if strings.Contains(contentType, "gif") {
		that.ImageType = ImageGif
	}
	reader := bytes.NewReader(that.Data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	b := img.Bounds()
	that.Width = b.Max.X
	that.Height = b.Max.Y
	return nil
}

/**
图片转webP
*/
func (that *Image) ToWebP(quality float32) (out []byte, err error) {
	var img image.Image
	reader := bytes.NewReader(that.Data)
	switch that.ImageType {
	case ImageJpeg:
		img, _ = jpeg.Decode(reader)
		break
	case ImagePng:
		img, _ = png.Decode(reader)
		break
	case ImageBmp:
		img, _ = bmp.Decode(reader)
		break
	case ImageGif:
		return that.gitToWebP(that.Data, quality)
	}
	if img == nil {
		msg := "image file " + that.FilePath + " is corrupted or not supported"
		err = errors.New(msg)
		return nil, err
	}
	var buf bytes.Buffer
	if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

/**
git 转webP
*/
func (that *Image) gitToWebP(gifBin []byte, quality float32) (webPBin []byte, err error) {
	converter := giftowebp.NewConverter()
	converter.LoopCompatibility = false
	//0 有损压缩  1无损压缩
	converter.WebPConfig.SetLossless(0)
	//压缩速度  0-6  0最快 6质量最好
	converter.WebPConfig.SetMethod(0)
	converter.WebPConfig.SetQuality(quality)
	//搞不懂什么意思,例子是这样用的
	converter.WebPAnimEncoderOptions.SetKmin(9)
	converter.WebPAnimEncoderOptions.SetKmax(17)

	return converter.Convert(gifBin)
}

// 计算图片缩放后的尺寸
func (that *Image) calculateRatioFit(srcWidth, srcHeight int) (int, int) {
	ratio := math.Min(DefaultMaxWidth/float64(srcWidth), DefaultMaxHeight/float64(srcHeight))
	return int(math.Ceil(float64(srcWidth) * ratio)), int(math.Ceil(float64(srcHeight) * ratio))
}

/**
创建缩略图
*/
func (that *Image) MakeThumbnail(width int, height int) (out []byte, err error) {
	w, h := that.calculateRatioFit(width, height)
	var img image.Image
	reader := bytes.NewReader(that.Data)
	switch that.ImageType {
	case ImageJpeg:
		img, _ = jpeg.Decode(reader)
		break
	case ImagePng:
		img, _ = png.Decode(reader)
		break
	case ImageBmp:
		img, _ = bmp.Decode(reader)
		break
	case ImageGif:
		gifData, err := that.resizeGif(width, height)
		if err != nil {
			return nil, err
		}
		return that.gitToWebP(gifData, 100)
	}
	if img == nil {
		msg := "image file " + that.FilePath + " is corrupted or not supported"
		err = errors.New(msg)
		return nil, err
	}
	var buf bytes.Buffer
	m := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	if err = webp.Encode(&buf, m, &webp.Options{Lossless: false, Quality: 100}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

/**
改变gif的长宽
*/
func (that *Image) resizeGif(width int, height int) (out []byte, err error) {
	reader := bytes.NewReader(that.Data)
	im, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, err
	}
	// reset the gif width and height
	im.Config.Width = width
	im.Config.Height = height

	firstFrame := im.Image[0].Bounds()
	img := image.NewRGBA(image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy()))

	// resize frame by frame
	for index, frame := range im.Image {
		b := frame.Bounds()
		draw.Draw(img, b, frame, b.Min, draw.Over)
		im.Image[index] = that.imageToPaletted(resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor))
	}
	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, im)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (that *Image) imageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}