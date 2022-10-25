package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/esimov/stackblur-go"
)

func main() {
	width, height := uint(1600), uint(800)

	r_shift, g_shift, b_shift := uint(3), uint(4), uint(5)
	blur := uint32(1)

	interesting_params := []JuliaParameters{
		{-1.5, 0, r_shift, g_shift, b_shift},
		{-0.5, 0.75, r_shift, g_shift, b_shift},
		{-0.25, 0.75, r_shift, g_shift, b_shift},
		{0.0, 1.0, r_shift, g_shift, b_shift},
		{0.0, 0.75, r_shift, g_shift, b_shift},
		{-0.75, 0.25, r_shift, g_shift, b_shift},
	}

	for i := 0; i < len(interesting_params); i++ {
		p := interesting_params[i]

		fmt.Println("Generating image with ", p)

		image := gen_img(width, height, p, blur)
		name := fmt.Sprintf("julia-cx_%.2f-cy_%.2f-r_%d-g_%d-b_%d.png", p.cx, p.cy, p.r_shift, p.g_shift, p.b_shift)

		fmt.Println("Writing result as ", name)

		write_png(name, image)
	}
}

type JuliaParameters struct {
	cx      float64
	cy      float64
	r_shift uint
	g_shift uint
	b_shift uint
}

func gen_img(width, height uint, julia_params JuliaParameters, blur uint32) image.Image {
	rectangle := image.Rect(0, 0, int(width), int(height))
	image := image.NewNRGBA(rectangle)

	bounds := image.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			image.Set(x, y, julia_color(x, y, width, height, &julia_params))
		}
	}

	blurred, err := stackblur.Process(image, blur)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return blurred
}

func julia_color(img_x, img_y int, img_width, img_height uint, parameters *JuliaParameters) color.NRGBA {
	width, height := float64(img_width), float64(img_height)
	x, y := float64(img_x), float64(img_y)

	zx := 3.0 * (x - 0.5*width) / width
	zy := 2.0 * (y - 0.5*height) / height

	i, iterations := 0, 1024

	for ; i < iterations && zx*zx+zy*zy < 4.0; i++ {
		t := zx*zx - zy*zy + parameters.cx
		zy = 2.0*zx*zy + parameters.cy
		zx = t
	}

	// Weight the color by bit shifting. More bit shifting reduces the red,
	// blue, or green of the pixel.
	r := uint8(i << int(parameters.r_shift))
	g := uint8(i << int(parameters.g_shift))
	b := uint8(i << int(parameters.b_shift))

	return color.NRGBA{r, g, b, math.MaxUint8}
}

func write_png(name string, image image.Image) {
	file, err := os.Create(name)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = png.Encode(file, image)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
