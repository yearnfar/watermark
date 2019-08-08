package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"time"
	"unicode/utf8"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
)

func main() {
	srcImg, err := imaging.Open("src.png")
	if err != nil {
		log.Fatal(err)
	}

	//workpath := time.Now().Unix()

	text := fmt.Sprintf("您好，世界！今天日期：%s", time.Now().Format("2006-01-02"))
	textImg := makeTextImage(text, 32, 60)

	// 调试
	imaging.Save(textImg, "text_img.png")

	offset := image.Pt(srcImg.Bounds().Dx()-textImg.Bounds().Dx(), srcImg.Bounds().Dy()-textImg.Bounds().Dy())

	newImg := image.NewRGBA(srcImg.Bounds())
	draw.Draw(newImg, srcImg.Bounds(), srcImg, image.ZP, draw.Src)
	draw.Draw(newImg, textImg.Bounds().Add(offset), textImg, image.ZP, draw.Over)

	output, err := os.Create("watermark.jpg")
	//output, err := os.Create(fmt.Sprintf("watermark_%d.jpg", timestamp))
	if err != nil {
		log.Fatal(err)
	}

	defer output.Close()
	err = jpeg.Encode(output, newImg, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
	}
	return

}

func makeTextImage(text string, fontsize, rotate float64) *image.NRGBA {
	fontBytes, err := ioutil.ReadFile("msyh.ttc")
	if err != nil {
		log.Println(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(font)
	ctx.SetFontSize(fontsize)
	ctx.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}))

	txtImg := image.NewNRGBA(image.Rect(0, 0, int(fontsize)*utf8.RuneCountInString(text), int(fontsize)))
	ctx.SetClip(txtImg.Bounds())
	ctx.SetDst(txtImg)

	pt := freetype.Pt(0, int(-fontsize/6)+ctx.PointToFixed(fontsize).Ceil())
	ctx.DrawString(text, pt)

	if rotate > 0 {
		return imaging.Rotate(txtImg, rotate, color.Transparent)
	} else {
		return txtImg
	}
}
