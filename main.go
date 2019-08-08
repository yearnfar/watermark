package main

import (
	//"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	img, err := imaging.Open("1.png")
	if err != nil {
		log.Fatal(err)
	}

	txtImg := makeText()
	imaging.Save(txtImg, "txt.png")

	offset := image.Pt(img.Bounds().Dx()-txtImg.Bounds().Dx()-120, img.Bounds().Dy()-txtImg.Bounds().Dy()-40)

	//根据b画布的大小新建一个新图像
	m := image.NewRGBA(img.Bounds())
	draw.Draw(m, img.Bounds(), img, image.ZP, draw.Src)
	draw.Draw(m, txtImg.Bounds().Add(offset), txtImg, image.ZP, draw.Over)
	//draw.Draw(m, txtImg.Bounds(), txtImg, image.ZP, draw.Over)

	imgw, err := os.Create("new.jpg")
	jpeg.Encode(imgw, m, &jpeg.Options{100})
	defer imgw.Close()

	return

}

func makeText() *image.NRGBA {
	txtImg := image.NewNRGBA(image.Rect(0, 0, 120, 40))

	//拷贝一个字体文件到运行目录
	fontBytes, err := ioutil.ReadFile("msyh.ttc")
	if err != nil {
		log.Println(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
	}

	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(font)
	f.SetFontSize(22)
	f.SetClip(txtImg.Bounds())
	f.SetDst(txtImg)
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}))

	pt := freetype.Pt(txtImg.Bounds().Dx()-120, txtImg.Bounds().Dy()-15)
	f.DrawString("你好，世界！", pt)
	return imaging.Rotate(txtImg, 60, color.Transparent)
}
