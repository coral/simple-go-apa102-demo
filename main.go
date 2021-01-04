package main

import (
	"flag"
	"time"

	"github.com/coral/simple-go-apa102-demo/strip"
)

var mhz = flag.Int64("megahertz", 6, "what mhz to clock SPI at")

func main() {

	numPixels := 32
	pixels := make([]strip.Pixel, numPixels)

	ls, err := strip.NewStrip(numPixels, *mhz)
	if err != nil {
		panic(err)
	}

	index := 0
	direction := false

	for {
		if index >= numPixels {
			direction = false
		}
		if index <= 0 {
			direction = true
		}

		if direction {
			index++
		} else {
			index--
		}
		//Rita pixlar
		for i, _ := range pixels {
			//Om vi är samma som index, tänd pixeln
			if i == index {
				pixels[i] = strip.Pixel{
					R: 1.0,
					G: 1.0,
					B: 1.0,
				}
			} else {
				pixels[i] = strip.Pixel{
					R: 0.0,
					G: 0.0,
					B: 0.0,
				}
			}
		}

		//Rendera pixlar
		ls.Render(pixels)
		time.Sleep(10 * time.Millisecond)
	}

}
