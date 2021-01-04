package strip

import (
	"log"
	"strconv"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices/apa102"
	"periph.io/x/periph/host"
)

type Pixel struct {
	R float64
	G float64
	B float64
}

type Strip struct {
	NumPixles int
	strip     *apa102.Dev
	buffer    []byte
}

func NewStrip(numPixels int, mhz int64) (Strip, error) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	s1, err := spireg.Open("/dev/spidev0.0")
	if err != nil {
		panic(err)
	}
	defer s1.Close()
	dd := physic.MegaHertz
	dd.Set(strconv.FormatInt(mhz, 10) + "MHz")

	if err := s1.LimitSpeed(dd); err != nil {
		return Strip{}, err
	}

	if p, ok := s1.(spi.Pins); ok {
		log.Printf("Using pins CLK: %s  MOSI: %s  MISO: %s", p.CLK(), p.MOSI(), p.MISO())
	}

	opts := apa102.PassThruOpts
	opts.NumPixels = numPixels
	opts.Intensity = 50
	a, err := apa102.New(s1, &opts)
	return Strip{
		NumPixles: numPixels,
		strip:     a,
		buffer:    []byte{},
	}, nil
}

func (s *Strip) Render(p []Pixel) {
	s.buffer = s.buffer[:0]
	for i, p := range p {
		s.buffer = append(s.buffer, []byte{
			clamp255(p.R * 255),
			clamp255(p.G * 255),
			clamp255(p.B * 255),
			byte(i),
		}...)
	}
	s.strip.Write(s.buffer)
}

func clamp255(v float64) byte {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}

	return byte(v)
}
