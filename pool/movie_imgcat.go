package pool

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"golang.org/x/image/bmp"
)

/**
裁减图片
* param		string		bigImg    	//大图路径
* param		string		smallImg		//小图路径
* param		int			width		//裁减结果的宽度
* param		int			height		//裁减结果的高度
*/
func CutImg(bigImg, smallImg string, width, height int) {

}

/**
*只裁减图片
 */
func Cutter(in image.Image, out io.Writer, fm string, width, height, quality int) error {

	//获取原始图片的宽度和高度
	//oW := in.Bounds().Max.X
	//oH := in.Bounds().Max.Y

	//配置切图参数
	cfg := cutter.Config{}
	cfg.Width = width
	cfg.Height = height
	cfg.Anchor = image.Point{0, 0}
	cfg.Options = cutter.Copy

	//切割图片
	cImg, err2 := cutter.Crop(in, cfg)

	if err2 != nil {
		return err2
	}

	//创建固定长宽的背景图
	bgImg := image.NewRGBA(image.Rect(0, 0, width, height))
	//将背景图涂抹成黑色
	for x := 0; x < bgImg.Bounds().Dx(); x++ {
		for y := 0; y < bgImg.Bounds().Dy(); y++ {
			bgImg.Set(x, y, color.Black)
		}
	}

	//将切割后的新图，画到背景图上
	draw.Draw(bgImg, bgImg.Bounds(), cImg, image.Pt(0, 0), draw.Over)

	//保存新的图片
	switch fm {
	case "jpeg":
		return jpeg.Encode(out, bgImg, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, bgImg)
	case "gif":
		return gif.Encode(out, bgImg, &gif.Options{})
	case "bmp":
		return bmp.Encode(out, bgImg)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}

/*
* 只生成缩略图
* 入参:
* 规则: 如果width 或 hight其中有一个为0，则大小不变 如果精度为0则精度保持不变
* 矩形坐标系起点是左上
* 返回:error
 */
func Scale(in io.Reader, width, height, quality int) (image.Image, string, error) {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return nil, fm, err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	//保存临时文件
	out, _ := os.Create("big_test.jpg")
	defer out.Close()
	switch fm {
	case "jpeg":
		jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		png.Encode(out, canvas)
	case "gif":
		gif.Encode(out, canvas, &gif.Options{})
	case "bmp":
		bmp.Encode(out, canvas)
	default:
		errors.New("ERROR FORMAT")
	}

	return canvas, fm, nil
}
