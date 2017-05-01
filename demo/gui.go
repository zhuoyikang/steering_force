package main

import (
	"fmt"
	"oops/common"
	"oops/common/massmap"

	"github.com/google/gxui"
	"github.com/google/gxui/gxfont"
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

type Gui struct {
	size  math.Size
	win   gxui.Window
	drv   gxui.Driver
	theme gxui.Theme

	Scale float64
	chart massmap.MassChart

	bk       gxui.Canvas
	dy       gxui.Image
	posLab   gxui.Label
	charLab  gxui.Label
	linerLab gxui.LinearLayout
}

func fy(y float64) float64 {
	return massmap.HEIGHT - y
}

func fp(pv []gxui.PolygonVertex) []gxui.PolygonVertex {
	for idx, _ := range pv {
		pv[idx].Position.Y = massmap.HEIGHT - pv[idx].Position.Y
	}

	return pv
}

// 因为浮点在绘图库中强制变成整数会导致精度丢失，这里将坐标放大，通过缩放gui的方式减小误差。
func ssr(i float64) int {
	return int(i * RATE)
}

func NewGui(scale, width, height float64) *Gui {
	return &Gui{
		Scale: scale,
		size: math.Size{
			W: ssr(width),
			H: ssr(height),
		},
	}
}

func (w *Gui) WinMain(driver gxui.Driver) {
	w.theme = dark.CreateTheme(driver)
	w.win = w.theme.CreateWindow(massmap.WIDTH*5, massmap.HEIGHT*5, "gui")
	w.win.SetScale(float32(w.Scale))
	w.drv = driver

	w.dy = w.theme.CreateImage()
	f, _ := w.drv.CreateFont(gxfont.Default, 40)
	w.posLab = w.theme.CreateLabel()
	w.posLab.SetFont(f)
	w.charLab = w.theme.CreateLabel()
	w.charLab.SetFont(f)

	table := w.theme.CreateLinearLayout()
	table.AddChild(w.posLab)
	table.AddChild(w.charLab)
	w.linerLab = table

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

func charStr(c massmap.CharChart) string {
	if c.CharType == pb.CharType_Pass {
		return fmt.Sprintf("PASS [%d]\nHP [%d] <%.0f, %.0f> \nAL [%d] \nTroop %v",
			c.CharId, c.HP, c.Pos.X, c.Pos.Y, c.AllianceId, c.Troops)
	} else {
		if c.PlayerId == -2 {
			return fmt.Sprintf("NPC CHAR[%d]\nHP [%d] <%.0f, %.0f> \nState %v",
				c.CharId, c.HP, c.Pos.X, c.Pos.Y, c.Stat)
		}
		return fmt.Sprintf("CHAR [%d]\nHP [%d] <%.0f, %.0f> \nAL [%d] \nTroop %v",
			c.CharId, c.HP, c.Pos.X, c.Pos.Y, c.AllianceId, c.Troops)
	}
}

func castleStr(c massmap.CastleChart) string {
	s := ""
	if c.Fallen {
		s = "Fallen"
	}
	return fmt.Sprintf("[%s] %d\nHP [%d] Warror[%d] Catapult[%d] %s\nAL [%d]",
		c.Name, c.CharId, c.HP, c.WarriorHP, c.CatapultHP, s, c.AllianceId)
}

func (w *Gui) Click(x, y float64) {
	fmt.Printf("click x %v y %v \n", x, y)
	// p := ln.Vector{X: x, Y: y}
	// for _, c := range w.chart.Chars {
	// 	if c.Pos.Distance(p) < 5 {
	// 		w.SetActiveChar(c.CharId)
	// 		w.SetCharLable(charStr(c))
	// 		return
	// 	}
	// }

	// for _, c := range w.chart.Castles {
	// 	if c.Pos.Distance(p) < 5 {
	// 		w.SetActiveChar(c.CharId)
	// 		w.SetCharLable(castleStr(c))
	// 		return
	// 	}
	// }

	// w.SetCharLable("")
	// w.SetActiveChar(0)
}

func rec(p ln.Vector, r float64) math.Rect {
	return math.Rect{
		Min: math.Point{
			W: ssr(p.X - r),
			H: ssr(fy(p.Y) - r),
		},
		Max: math.Point{
			W: ssr(p.X + r),
			H: ssr(fy(p.Y) + r),
		},
	}
}

func (w *Gui) DrawMassChart() {
	chart := w.chart
	if w.bk == nil {
		w.bk = w.drv.CreateCanvas(w.size)
		w.drawMesh(w.bk)
		w.bk.Complete()

		image := w.theme.CreateImage()
		image.SetCanvas(w.bk)
		w.win.AddChild(image)
		w.win.AddChild(w.dy)
		//w.win.AddChild(w.posLab)
		//w.win.AddChild(w.charLab)
		w.win.AddChild(w.linerLab)

		w.posLab.SetColor(gxui.Green70)

		w.charLab.SetColor(gxui.Red90)
		w.charLab.SetMultiline(true)
		//w.charLab.SetMargin(math.Spacing{0, 10, 0, 0})
	}

	canvas := w.drv.CreateCanvas(w.size)

	for _, m := range chart.Moves {
		canvas.DrawLines(pathPolygon(m), p2)
	}

	for _, char := range chart.Chars {
		if char.CharId == w.activeChar {
			w.SetCharLable(charStr(char))
		}
		w.drawChar(canvas, &char)
	}

	for _, castle := range chart.Castles {
		if castle.CharId == w.activeChar {
			w.SetCharLable(castleStr(castle))
		}
		w.drawCastle(canvas, &castle)
	}

	canvas.Complete()
	w.dy.SetCanvas(canvas)
}

// func (w *Gui) drawChar(canvas gxui.Canvas, char *massmap.CharChart) {
// 	b := b4

// 	if char.MoveSt == pb.MoveState_Move {
// 		b = b2
// 	}

// 	if char.CharType == pb.CharType_Pass {
// 		b = b1
// 	}

// 	canvas.DrawRoundedRect(rec(char.Pos, char.R), 50, 50, 50, 50,
// 		gxui.TransparentPen, b)
// }

// func (w *Gui) drawCastle(canvas gxui.Canvas, castle *massmap.CastleChart) {
// 	pl := common.CalcSlots(castle.Pos, 12, 0.612+castle.R)
// 	canvas.DrawRoundedRect(rec(castle.Pos, castle.R), 50, 50, 50, 50, p01, b3)
// 	for _, p := range pl {
// 		canvas.DrawRoundedRect(rec(p, 0.612), 50, 50, 50, 50, p01,
// 			gxui.TransparentBrush)
// 	}
// }

func (w *Gui) drawMesh(canvas gxui.Canvas) {
	for k := 0; k < len(w.nm.Triangles); k++ {
		w.drawTri(canvas, w.nm.Triangles[k])
	}

	for i := 0; i < ssr(massmap.WIDTH); i += ssr(massmap.GRID_WIDTH) {
		canvas.DrawLines(p1p2(i, 0, i, ssr(massmap.HEIGHT)), p3)
	}

	for i := 0; i < ssr(massmap.HEIGHT); i += ssr(massmap.GRID_HEIGHT) {
		canvas.DrawLines(p1p2(0, i, ssr(massmap.WIDTH), i), p3)
	}
}

func (w *Gui) drawTri(canvas gxui.Canvas, tri common.TriangleCfg) {
	p := tri.Vertices
	poly := []gxui.PolygonVertex{
		gxui.PolygonVertex{
			Position: math.Point{
				W: ssr(w.nm.Vertices[p[0]].X),
				H: ssr(fy(w.nm.Vertices[p[0]].Y)),
			}},

		gxui.PolygonVertex{
			Position: math.Point{
				W: ssr(w.nm.Vertices[p[1]].X),
				H: ssr(fy(w.nm.Vertices[p[1]].Y)),
			}},

		gxui.PolygonVertex{
			Position: math.Point{
				W: ssr(w.nm.Vertices[p[2]].X),
				H: ssr(fy(w.nm.Vertices[p[2]].Y)),
			}},
	}
	b := b1
	switch tri.Type {
	case 1:
		b = b11
	case 2:
		b = b12
	}
	canvas.DrawPolygon(poly, p3, b)
}

func (w *Gui) Fy(y float64) float64 {
	return fy(y)
}

func (w *Gui) SetActiveChar(charId int64) {
	w.activeChar = charId
}

func (w *Gui) SetChart() {
}

func pathPolygon(m massmap.MoveChart) []gxui.PolygonVertex {
	var poly []gxui.PolygonVertex
	poly = append(poly,
		gxui.PolygonVertex{
			Position: math.Point{
				W: ssr(m.Src.X),
				H: ssr(fy(m.Src.Y)),
			}})

	for k := range m.Path {
		poly = append(poly,
			gxui.PolygonVertex{
				Position: math.Point{
					W: ssr(m.Path[k].X),
					H: ssr(fy(m.Path[k].Y)),
				}})
	}
	poly = append(poly,
		gxui.PolygonVertex{
			Position: math.Point{
				W: ssr(m.Dst.X),
				H: ssr(fy(m.Dst.Y)),
			}})

	return poly
}

func p1p2(x1, y1, x2, y2 int) []gxui.PolygonVertex {
	return []gxui.PolygonVertex{
		gxui.PolygonVertex{Position: math.Point{W: x1, H: y1}},
		gxui.PolygonVertex{Position: math.Point{W: x2, H: y2}},
	}
}
