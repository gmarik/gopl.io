package main

import (
	"bufio"
	"golang.org/x/net/html"
	"io"
	"math"
	"strings"
	"testing"
)

func countWords(rd io.Reader) (int, error) {
	freq := make(map[string]int)
	input := bufio.NewScanner(rd)
	input.Split(bufio.ScanWords)

	for input.Scan() {
		if err := input.Err(); err != nil {
			return 0, err
		}

		freq[strings.ToLower(input.Text())] += 1
	}
	return len(freq), nil
}

func countWordsAndImages(n *html.Node) (nwords, nimages int) {
	if n.Type == html.ElementNode && n.Data == "img" {
		nimages += 1
	}

	if n.Type == html.TextNode {
		nw, err := countWords(strings.NewReader(n.Data))
		if err != nil {
			panic(err)
		}
		nwords += nw
	}

	// log.Println(n)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nw, ni := countWordsAndImages(c)
		nwords += nw
		nimages += ni
	}

	return nwords, nimages
}

func corner(i, j int) (sx, sy float64) {

	const (
		width, height = 600, 320            // canvas size in pixels
		cells         = 100                 // number of grid cells
		xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
		xyscale       = width / 2 / xyrange // pixels per x or y unit
		zscale        = height * 0.001      // pixels per z unit
		angle         = math.Pi / 9         // angle of x, y axes (=30°)
	)
	var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := x*x - y*y

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx = width/2 + (x-y)*cos30*xyscale
	sy = height/2 + (x+y)*sin30*xyscale - z*zscale
	return
}

func Test_Corner(t *testing.T) {

	expx, expy := 297.18092213764226, 57.67853829223041
	gotx, goty := corner(1, 2)
	if gotx != expx || goty != expy {
		t.Errorf("\nExp: %v\nGot: %v", expx, gotx)
		t.Errorf("\nExp: %v\nGot: %v", expy, goty)
	}

}

func Test_countWordsAndImages(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(snippet))
	if err != nil {
		t.Fatal(err)
	}

	nwords, nimages := countWordsAndImages(doc)
	if !(nwords == 121 && nimages == 1) {
		t.Errorf("Unexpected results: words: %d, images: %d", nwords, nimages)
	}
}

const snippet = `
<html>
<body>
  <div class="wrapper">
    <a class="site-title" href="/">
      <img src="/images/head.jpg" height="55" width="50">
      Mike Perham
    </a>
		<article class="post-content">
			<p>I’ve been exploring a few new (to me!) technologies recently and <a href="http://smarden.org/runit/">runit</a> is one that I’ve come away really impressed with. Linux distros have a few competing init services available: Upstart, systemd, runit or creaky old sysvinit. Having researched all of them and having built lots of server-side systems over the last two decades, I can firmly recommend runit if you want a server-focused, reliable init system based on the traditional Unix philosophy.</p>
			<p>The point of an init system is to start and supervise processes when the machine boots up. If you’re building a modern web site, you want memcached, redis, postgresql, mysql and other daemons to start up immediately when the machine boots. Supervision means the init system will restart the process immediately if it disappears for some reason, e.g. a crash. Reliability of the init system is paramount so simplicity is a key attribute.</p>
		</article>
	</div>
</body>
</html>
`
