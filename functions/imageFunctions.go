package functions

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
)

func transformToNRGBA(pathImage string) (newImg *image.NRGBA) {
	imgFile, err := os.Open(pathImage)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	im, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Cannot decode image:", err)
		return
	}

	b := im.Bounds()
	// fmt.Println("Type: ", reflect.TypeOf(b.Dx()))
	// fmt.Println("b.Dx: ", b.Dx())
	// fmt.Println("b.Dy: ", b.Dy())
	newImg = image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newImg, newImg.Bounds(), im, b.Min, draw.Src)
	return
}

func GetImageInfo(pathImage string) (dx int, dy int, pix []uint8) {
	imgFile, err := os.Open(pathImage)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	im, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Cannot decode image:", err)
		return
	}

	b := im.Bounds()
	newImg := transformToNRGBA(pathImage)

	dx = b.Dx()
	dy = b.Dy()
	pix = newImg.Pix

	fmt.Println("Dx: ", dx)
	fmt.Println("Dy: ", dy)

	return
}

// build an image with the given parameters
func BuildImage(rect image.Rectangle, pix []uint8) (created *image.NRGBA) {

	created = &image.NRGBA{
		Pix:    pix,
		Stride: rect.Dx() * 4,
		Rect:   rect,
	}
	return
}

//save an image
func Save(filePath string, img *image.NRGBA) {
	imgFile, err := os.Create(filePath)
	defer imgFile.Close()
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	png.Encode(imgFile, img.SubImage(img.Rect))
}
