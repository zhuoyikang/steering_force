package steering_force

import (
	"math"
)

type Vector2D struct {
	X, Y float64
}

func (a Vector2D) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

func (a Vector2D) Distance(b Vector2D) float64 {
	return a.Sub(b).Length()
}

func (a Vector2D) LengthSquared() float64 {
	return a.X*a.X + a.Y*a.Y
}

func (a Vector2D) DistanceSquared(b Vector2D) float64 {
	return a.Sub(b).LengthSquared()
}

func (a Vector2D) Dot(b Vector2D) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vector2D) Normalize() Vector2D {
	d := a.Length()
	return Vector2D{a.X / d, a.Y / d}
}

func (a Vector2D) Perp() Vector2D {
	return Vector2D{-a.Y, a.X}
}

func (a Vector2D) Truncate(max float64) Vector2D {
	if a.Length() > max {
		return a.Normalize().MulScalar(max)
	} else {
		return a
	}
}

func (a Vector2D) Add(b Vector2D) Vector2D {
	return Vector2D{a.X + b.X, a.Y + b.Y}
}

func (a Vector2D) Sub(b Vector2D) Vector2D {
	return Vector2D{a.X - b.X, a.Y - b.Y}
}

func (a Vector2D) Mul(b Vector2D) Vector2D {
	return Vector2D{a.X * b.X, a.Y * b.Y}
}

func (a Vector2D) Div(b Vector2D) Vector2D {
	return Vector2D{a.X / b.X, a.Y / b.Y}
}

func (a Vector2D) AddScalar(b float64) Vector2D {
	return Vector2D{a.X + b, a.Y + b}
}

func (a Vector2D) SubScalar(b float64) Vector2D {
	return Vector2D{a.X - b, a.Y - b}
}

func (a Vector2D) MulScalar(b float64) Vector2D {
	return Vector2D{a.X * b, a.Y * b}
}

func (a Vector2D) DivScalar(b float64) Vector2D {
	return Vector2D{a.X / b, a.Y / b}
}

func (a Vector2D) Min(b Vector2D) Vector2D {
	return Vector2D{math.Min(a.X, b.X), math.Min(a.Y, b.Y)}
}

func (a Vector2D) Max(b Vector2D) Vector2D {
	return Vector2D{math.Max(a.X, b.X), math.Max(a.Y, b.Y)}
}
