// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.

// Exercise 1.12: Modify the Lissajous server to read parameter values from the
// URL. For example, you might arrange it so that a URL like
// http://localhost:8000/?cycles=20 sets the number of cycles to 20 instead of
// the default 5. Use the strconv.Atoi function to convert the string parameter
// into an integer. You can see its documentation with go doc strconv.Atoi.

//!+main

// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"fmt"
	"io"

	"image"
	"image/color"
	"image/gif"

	"math"
	"math/rand"
)

//!-main
// Packages not needed by version in book.
import (
	"log"
	"net/http"
	"strconv"
	"time"
)

//!+main

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	//!-main
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.

	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)

	default_par := params{
		cycles:  5,
		res:     0.001,
		size:    100,
		nframes: 64,
		delay:   8,
		freq:    rand.Float64(), // relative frequency of y oscillator
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "error parsing form", 500)
			return
		}

		//copy
		par := default_par

		if err := parseParams(r, &par); err != nil {
			http.Error(w, "error parsing params", 500)
			return
		}

		lissajous(w, &par)
	}
	http.HandleFunc("/lissajous", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//copy
		par := default_par

		if err := parseParams(r, &par); err != nil {
			http.Error(w, "error parsing params", 500)
			return
		}

		fmt.Fprintf(w, `
			<html>
				<body>
					<h2> Examples </h2>

					<ul>
					<li> <a href="/?freq=0.4&cycles=1&size=150"> Heart </a>
					</li>
					<li><a href="/?freq=1.4&cycles=5&size=200"> Multiple</a>
					</li>
					<li><a href="/?freq=1.0&cycles=5&size=200"> Another </a>
					</li>
					<li><a href="/?freq=1.4&cycles=12&size=150&res=0.04"> Pixelated (low <b>res</b>olution </a>
					</li>
					</ul>

					<h2> Params</h2>
					<div> %v </div>

					<img src="/lissajous?%s"/>

				</body>
			</html>
		`, par, r.URL.RawQuery)
	})

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
	//!+main
}

func parseParams(r *http.Request, p *params) error {
	if str := r.FormValue("freq"); str != "" {
		if v, err := strconv.ParseFloat(str, 20); err != nil {
			return err
		} else {
			p.freq = v
		}
	}

	if str := r.FormValue("res"); str != "" {
		if v, err := strconv.ParseFloat(str, 8); err != nil {
			return err
		} else {
			p.res = v
		}
	}

	if str := r.FormValue("size"); str != "" {
		if v, err := strconv.ParseUint(str, 10, 32); err != nil {
			return err
		} else {
			p.size = uint(v)
		}
	}

	if str := r.FormValue("cycles"); str != "" {
		if v, err := strconv.ParseUint(str, 10, 32); err != nil {
			return err
		} else {
			p.cycles = uint(v)
		}
	}

	if str := r.FormValue("nframes"); str != "" {
		if v, err := strconv.ParseUint(str, 10, 32); err != nil {
			return err
		} else {
			p.nframes = uint(v)
		}
	}

	if str := r.FormValue("delay"); str != "" {
		if v, err := strconv.ParseUint(str, 10, 32); err != nil {
			return err
		} else {
			p.delay = uint(v)
		}
	}

	return nil
}

type params struct {
	cycles  uint    // 5     // number of complete x oscillator revolutions
	res     float64 // 0.001 // angular resolution
	size    uint    // 100   // image canvas covers [-size..+size]
	nframes uint    // 64    // number of animation frames
	delay   uint    // 8     // delay between frames in 10ms units
	freq    float64 // rand  // relative frequency of y oscillator

}

func lissajous(out io.Writer, p *params) {
	var (
		cycles  = float32(p.cycles)
		res     = p.res
		size    = int(p.size)
		nframes = int(p.nframes)
		delay   = p.delay
		freq    = p.freq * 3.0
	)

	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += float64(res) {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, int(delay))
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main
