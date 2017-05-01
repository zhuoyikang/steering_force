package main

import (
	"fmt"
	"time"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/themes/dark"
)

var (
	p1  = gxui.CreatePen(1, gxui.Black)
	p2  = gxui.CreatePen(1.5, gxui.Green50)
	p3  = gxui.CreatePen(1, gxui.Gray10)
	p4  = gxui.CreatePen(1, gxui.Red90)
	p0  = gxui.CreatePen(1.5, gxui.Yellow)
	p01 = gxui.CreatePen(1.5, gxui.Red)
	p02 = gxui.CreatePen(1.5, gxui.Red50)

	b1  = gxui.CreateBrush(gxui.Gray50)
	b11 = gxui.CreateBrush(gxui.Gray90)
	b12 = gxui.CreateBrush(gxui.Gray70)
	b2  = gxui.CreateBrush(gxui.Green30)
	b3  = gxui.CreateBrush(gxui.Red50)
	b4  = gxui.CreateBrush(gxui.Blue30)
)

const (
	RATE = 10
)

var (
	HEIGHT = float64(1000)
	WIDTH  = float64(1000)
)

type Gui struct {
	size  math.Size
	win   gxui.Window
	drv   gxui.Driver
	theme gxui.Theme

	Scale float64
	tick  time.Duration

	bk       gxui.Canvas
	dy       gxui.Image
	posLab   gxui.Label
	charLab  gxui.Label
	linerLab gxui.LinearLayout
}

func fy(y float64) float64 {
	return HEIGHT - y
}

func fp(pv []gxui.PolygonVertex) []gxui.PolygonVertex {
	for idx, _ := range pv {
		pv[idx].Position.Y = int(HEIGHT) - pv[idx].Position.Y
	}

	return pv
}

// 因为浮点在绘图库中强制变成整数会导致精度丢失，这里将坐标放大，通过缩放gui的方式减小误差。
func ssr(i float64) int {
	return int(i * RATE)
}

func NewGui(scale, width, height float64, tick time.Duration) *Gui {
	return &Gui{
		tick:  tick,
		Scale: scale,
		size: math.Size{
			W: ssr(width),
			H: ssr(height),
		},
	}
}

func (w *Gui) run() {
	c := time.Tick(w.tick)
	for range c {
		w.GetDriver().CallSync(w.DrawUI)
	}
}

func (w *Gui) Run() {
	go w.run()
	gl.StartDriver(w.appMain)
}

func (w *Gui) appMain(driver gxui.Driver) {
	w.WinMain(driver)

	win := w.GetWindow()
	win.OnMouseUp(w.onClick)

	win.OnClose(w.GetDriver().Terminate)
}

func (w *Gui) WinMain(driver gxui.Driver) {
	w.theme = dark.CreateTheme(driver)

	w.win = w.theme.CreateWindow(int(WIDTH), int(HEIGHT), "gui")
	w.win.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray50))

	w.drv = driver
	w.dy = w.theme.CreateImage()
}

func (w *Gui) GetWindow() gxui.Window {
	return w.win
}

func (w *Gui) GetDriver() gxui.Driver {
	return w.drv
}

func (w *Gui) SetPosLable(x, y float64) {
	w.posLab.SetText(fmt.Sprintf("%.1f, %.1f", x, y))
}

func (w *Gui) SetCharLable(s string) {
	w.charLab.SetText(s)
}

func (w *Gui) onClick(me gxui.MouseEvent) {
	fmt.Printf("click %v \n", me.Point)
}

func (w *Gui) DrawUI() {
	fmt.Printf("draw \n")
	canvas := w.drv.CreateCanvas(w.size)
	// draw
	canvas.Complete()
	w.dy.SetCanvas(canvas)
}

func (w *Gui) Fy(y float64) float64 {
	return fy(y)
}

func p1p2(x1, y1, x2, y2 int) []gxui.PolygonVertex {
	return []gxui.PolygonVertex{
		gxui.PolygonVertex{Position: math.Point{X: x1, Y: y1}},
		gxui.PolygonVertex{Position: math.Point{X: x2, Y: y2}},
	}
}
