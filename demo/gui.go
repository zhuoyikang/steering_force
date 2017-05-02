package main

import (
	"fmt"
	"os"
	"time"

	sf "github.com/indie21/steering_force"

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
	RATE = 1
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

	currentMouse  sf.Vector2D
	currentEntity *sf.Entity

	world *sf.World
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

// // 因为浮点在绘图库中强制变成整数会导致精度丢失，这里将坐标放大，通过缩放gui的方式减小误差。
func ssr(i float64) int {
	return int(i * RATE)
}

func NewGui(world *sf.World, scale, width, height float64, tick time.Duration) *Gui {
	return &Gui{
		tick:  tick,
		Scale: scale,
		world: world,
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
	gl.StartDriver(w.appMain)
}

func (w *Gui) appMain(driver gxui.Driver) {
	w.WinMain(driver)

	win := w.GetWindow()
	win.OnMouseUp(w.onClick)
	win.OnMouseMove(w.onMove)
	win.OnKeyDown(w.onKeyDown)

	go w.run()
	win.OnClose(w.GetDriver().Terminate)
}

func (w *Gui) WinMain(driver gxui.Driver) {
	w.theme = dark.CreateTheme(driver)

	w.win = w.theme.CreateWindow(int(WIDTH), int(HEIGHT), "gui")
	w.win.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray50))

	w.drv = driver
	w.dy = w.theme.CreateImage()
	w.win.AddChild(w.dy)
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

func (w *Gui) onMove(me gxui.MouseEvent) {
	w.currentMouse = sf.Vector2D{float64(me.Point.X), fy(float64(me.Point.Y))}
}

func (w *Gui) onKeyDown(ev gxui.KeyboardEvent) {
	switch ev.Key {
	case gxui.KeyEnter:
	case gxui.KeyE:
		if w.world.PosConflict(w.currentMouse) != nil {
			return
		}

		entity := sf.NewEntity()
		entity.SetPos(w.currentMouse)
		w.world.AddEntity(entity)
		w.currentEntity = entity

	case gxui.KeyF:
		if w.world.PosConflict(w.currentMouse) != nil {
			return
		}

		if w.currentEntity != nil {
			w.currentEntity.SetTarget(w.currentMouse)
			w.currentEntity = nil
		} else {
			entity := sf.NewEntity()
			entity.SetPos(w.currentMouse)
			w.world.AddEntity(entity)
			w.currentEntity = entity
		}
	case gxui.KeyD:
		entity := w.world.PosConflict(w.currentMouse)
		if entity != nil {
			entity.SetDelete()
		}

	case gxui.KeyQ:
		os.Exit(0)
	}
}

func (w *Gui) DrawUI() {
	canvas := w.drv.CreateCanvas(w.size)

	// draw
	w.DrawEntities(canvas)

	canvas.Complete()
	w.dy.SetCanvas(canvas)
}

func rec(p sf.Vector2D, r float64) math.Rect {
	return math.Rect{
		Min: math.Point{
			X: ssr(p.X - r),
			Y: ssr(fy(p.Y) - r),
		},
		Max: math.Point{
			X: ssr(p.X + r),
			Y: ssr(fy(p.Y) + r),
		},
	}
}

func (w *Gui) DrawEntity(e *sf.Entity, canvas gxui.Canvas) {
	//fmt.Printf("DrawEntity %v\n", e.GetPos())
	canvas.DrawRoundedRect(rec(e.GetPos(), e.GetBoundingRadius()), 50, 50, 50, 50,
		gxui.TransparentPen, b4)

	if e.IsTargetOn() {
		canvas.DrawRoundedRect(rec(e.GetTarget(), e.GetBoundingRadius()), 50, 50, 50, 50,
			gxui.TransparentPen, b12)
	}
}

func (w *Gui) DrawEntities(canvas gxui.Canvas) {
	for _, e := range w.world.AllEntities() {
		w.DrawEntity(e, canvas)
	}
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
