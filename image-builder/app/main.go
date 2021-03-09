package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"
)

type Circle struct {
	X, Y, R float64
}

func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	} else {
		return 255
	}
}

func main() {
	log.Info("build image")
	outputPath := os.Getenv("IMAGE_OUTPUT_PATH")
	log.Infof("output path: %s", outputPath)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		log.Fatalf("folder %s does not exist", outputPath)
	}

	img := createImage()

	timestamp := time.Now().Format("20060102150405")

	labelText := fmt.Sprintf("timestamp: %v", timestamp)
	addLabel(img, 20, 20, labelText)

	filePath := fmt.Sprintf("%s/image-%v.png", outputPath, timestamp)
	saveImage(filePath, img)
	log.Info("completed")
}

func saveImage(filePath string, img *image.RGBA) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	png.Encode(f, img)
	log.Infof("image saved %s", filePath)
}

func createImage() *image.RGBA {
	var w, h int = 280, 240
	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := 40.0
	θ := 2 * math.Pi / 3
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), 60}
	cg := &Circle{hw - r*math.Sin(θ), hh - r*math.Cos(θ), 60}
	cb := &Circle{hw - r*math.Sin(-θ), hh - r*math.Cos(-θ), 60}

	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := color.RGBA{
				R: cr.Brightness(float64(x), float64(y)),
				G: cg.Brightness(float64(x), float64(y)),
				B: cb.Brightness(float64(x), float64(y)),
				A: 255,
			}
			m.Set(x, y, c)
		}
	}
	log.Infof("created image %vx%v pixels", w, h)
	return m
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{R: 200, G: 100, A: 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
	log.Infof("added label %s", label)
}
