package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gen"
	app.Usage = "generate images"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name: "grid",
			Flags: []cli.Flag{
				cli.IntFlag{"width", 100, "image width"},
				cli.IntFlag{"height", 100, "image height"},
				cli.IntFlag{"grid", 10, "image gridSize"},
				cli.StringFlag{"file", "out.png", "image name"},
			},
			Usage: "generate a grid",
			Action: func(c *cli.Context) {
				MakeGrid(c.Int("width"), c.Int("height"), c.Int("grid"), c.String("file"))
			},
		},
		{
			Name:  "crush",
			Usage: "minimize pizel art",
			Action: func(c *cli.Context) {
				Crush(c.Args()...)
			},
		},
	}

	app.Run(os.Args)
}

func MakeGrid(w, h, gridSize int, file string) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			xEven := (i/gridSize)%2 == 0
			yEven := (j/gridSize)%2 == 0
			if xEven == yEven {
				img.SetRGBA(i, j, color.RGBA{0, 0, 0, 255})
			} else {
				img.SetRGBA(i, j, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = png.Encode(f, img)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Crush(files ...string) {
	for _, file := range files {
		fmt.Printf("Crushing %s\n", file)

		f, err := os.OpenFile(file, os.O_RDONLY, 0660)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			break
		}

		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
		bound := img.Bounds()

		column := []int{}
		for i := 0; i < bound.Dx(); i++ {
			last := img.At(i, 0)
			sum := 0
			for j := 0; j < bound.Dy(); j++ {
				next := img.At(i, j)
				if !equalColors(next, last) {
					column = append(column, j-sum)
					sum = j
					last = next
				}
			}
			column = append(column, bound.Dy()-sum)
		}
		verticalCrush := gcf(column...)

		row := []int{}
		for i := 0; i < bound.Dy(); i++ {
			last := img.At(i, 0)
			sum := 0
			for j := 0; j < bound.Dx(); j++ {
				next := img.At(j, i)
				if !equalColors(next, last) {
					row = append(row, j-sum)
					sum = j
					last = next
				}
			}
			row = append(row, bound.Dx()-sum)
		}
		horizontalCrush := gcf(row...)

		cf := gcf(horizontalCrush, verticalCrush)
		fmt.Println("Maximum crush potential: ", cf)

		newBound := image.Rect(bound.Min.X, bound.Min.Y, bound.Max.X/cf, bound.Max.Y/cf)
		outImg := image.NewRGBA(newBound)
		for i := 0; i < newBound.Dx(); i++ {
			for j := 0; j < newBound.Dy(); j++ {
				outImg.Set(i, j, img.At(i*cf, j*cf))
			}
		}

		outf, err := os.OpenFile("crushed-"+file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)

		err = png.Encode(outf, outImg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func equalColors(a, b color.Color) bool {
	R, G, B, A := a.RGBA()
	R2, G2, B2, A2 := b.RGBA()
	if R != R2 {
		return false
	}
	if G != G2 {
		return false
	}
	if B != B2 {
		return false
	}
	if A != A2 {
		return false
	}
	return true
}
