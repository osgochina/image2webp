package app

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

type MatchUrl struct {
	OriginUrl     string //原始url
	FilePath      string // 匹配到的文件名
	IsThump       bool   //是否是缩略图
	ThumpParams   string // 最终匹配出的缩略图参数
	Weight        int    // 缩略图高度
	Height        int    // 缩略图宽度
	SpecialParams string // 特殊参数，目前未使用到
}

// NewMatchUrl match 解析后的数组
/**
	//原图
	0=>/pics/png.jpg_webp //整个url
	1=>/pics/png.jpg	  //匹配到的文件名
	//缩略图
	0=>/pics/png.jpg_100_100_xy.jpg_webp //整个url
	1=>/pics/png.jpg					 //匹配到的文件名
	2=>_100_100_xy						 //最终匹配出的缩略图参数
	3=>100								 // 宽度
	4=>100								 // 高度
	5=>_xy								 //特殊参数，目前未使用到
**/
func NewMatchUrl(url string) (*MatchUrl, error) {
	match, err := gregex.MatchString(`^(.+?\..+)(_(\d+)_(\d+)(_.+)?)\.jpg_webp$`, url)
	if err != nil {
		return nil, err
	}
	m := &MatchUrl{
		OriginUrl: url,
		IsThump:   false,
	}
	if len(match) > 2 && len(match[2]) > 0 {
		m.IsThump = true
		m.FilePath = match[1]
		m.ThumpParams = match[2]
		m.Weight = gconv.Int(match[3])
		m.Height = gconv.Int(match[4])
		m.SpecialParams = match[5]
		return m, nil
	}
	match, err = gregex.MatchString(`^(.*)_webp$`, url)
	if err != nil {
		return nil, err
	}
	if len(match) <= 1 {
		return nil, errors.New("not found")
	}
	m.FilePath = match[1]
	return m, nil
}

// AllowSizes 判断缩略图参数是否符合规范
func (that *MatchUrl) AllowSizes(allowSizes *garray.StrArray) bool {
	return allowSizes.Contains(that.ThumpParams)
}

// GetRealFilePath 获取图片的真实文件地址
func (that *MatchUrl) GetRealFilePath(storagePath string) string {
	absFile := gfile.Abs(fmt.Sprintf("%s%s", storagePath, that.FilePath))
	return absFile
}

// IsMatch 判断传入的url是否符合规则
func IsMatch(url string) bool {
	match, _ := gregex.MatchString(`^(.*)_webp$`, url)
	if len(match) <= 0 {
		return false
	}
	return true
}
