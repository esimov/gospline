# gospline
Gospline is little Go library for transforming straight lines into curves.

The library is written in a modular way which means that it can be plugged in into different rendering methods. The examples provided below renders the resulted curves into image and svg, but since it has been developed using Go's philosophy of inheritance over encapsulation, it should supports other types of outputs until they implements the <a href="https://golang.org/pkg/io/#Writer">io.Writer</a> interface.

To render the output as image, the library implements the <a href="https://en.wikipedia.org/wiki/Xiaolin_Wu's_line_algorithm">Xiaolin Wu</a> antialiasing method, if the provided parameter is true, otherwise it implements the <a href="https://en.wikipedia.org/wiki/Bresenham's_line_algorithm">Bresenham</a> line algorithm. This means the library is not using the default Go <a href="https://github.com/golang/freetype/">github.com/golang/freetype/raster</a> package for rendering.

```go
func (img *Canvas) DrawLine(x1, y1, x2, y2 float64, col color.Color, antialiased bool) *Canvas {
	if antialiased {
		xiaolinWuLine(img, x1, y1, x2, y2, col.(color.NRGBA))
	} else {
		bresenhamLine(img, x1, y1, x2, y2, col.(color.NRGBA))
	}
	return img
}
```

## Installation
```bash
$ go get -u -f github.com/esimov/gospline
```

## Examples
Some test files have been included into the <strong>example</strong> directory. To run them type:
`go run examples/draw.go` for example.

This is how you initialize the base method:

```go
spline := NewBSpline(points, 3, false)
spline.Init()
```
...where the 2nd paramether means the degree of curvature. 

Below is an example to render the spline as svg in the web browser.

```go
svg := &spline.SVG{
  Width: width,
  Height: height,
  Title: "BSpline",
  Lines: []spline.Line{},
  Color: color.NRGBA{R:255,G:0,B:0,A:255},
  Description: "Convert straight lines to curves",
  StrokeWidth: 2,
  StrokeLineCap: "round", //butt, round, square
}

handler := func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "image/svg+xml")
  svg.Draw(w, points, false)
}
http.HandleFunc("/", handler)
http.ListenAndServe("localhost:8000", nil)
```
and the corresponding method to render the output into an image.

```go
raster := &spline.Image{
  Width : width,
  Height : height,
  Color : color.NRGBA{R:255,G:0,B:0,A:255},
}
output, _ := os.Create("./samples/curve.png")
defer output.Close()
raster.Draw(output, points, true)
```

You can even use the debug option to show the original lines.

```go
// Draw original line
if debug {
  for i:=0; i < len(points)-1; i++ {
    raster = canvas.DrawLine(points[i][0], points[i][1], points[i+1][0], points[i+1][1], color.NRGBA{R:155,G:155,B:155,A:70 }, false)
  }
}
```
This will produce an image like this with antialiasing mode set to true.
<img alt="BSPline" title="BSpline" src="https://raw.githubusercontent.com/esimov/gospline/master/samples/curve.png"/>

## License

This software is distributed under the MIT license found in the LICENSE file.
