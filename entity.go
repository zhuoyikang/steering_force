package steering_force

import (
	"fmt"
	"time"
)

var (
	DR                    = float64(20)
	SPEED                 = float64(60)
	MINDETECTIONBOXLENGTH = float64(20)
)

type Entity struct {
	pos            Vector2D
	boundingRadius float64
	velocity       Vector2D

	//the mass
	Mass float64

	//a normalized vector pointing in the direction the entity is heading.
	heading Vector2D

	//a vector perpendicular to the heading vector
	side Vector2D

	//the maximum speed this entity may travel at.
	MaxSpeed float64

	//the maximum force this entity can produce to power itself
	//(think rockets and thrust)
	MaxForce float64

	//the maximum rate (radians per second)this vehicle can rotate
	MaxTurnRate float64

	targetPos Vector2D
	targetOn  bool

	dDBoxLength float64

	delete bool

	world *World
}

func NewEntity() *Entity {
	e := &Entity{
		delete:         false,
		velocity:       Vector2D{0, 0},
		pos:            Vector2D{0, 0},
		targetOn:       false,
		targetPos:      Vector2D{0, 0},
		MaxSpeed:       SPEED,
		Mass:           1,
		boundingRadius: DR,
	}
	return e
}

func (e *Entity) SetDelete() {
	e.delete = true
}

func (e *Entity) SetPos(pos Vector2D) {
	e.pos = pos
}

func (e *Entity) GetPos() Vector2D {
	return e.pos
}

func (e *Entity) GetTarget() Vector2D {
	return e.targetPos
}

func (e *Entity) GetBoundingRadius() float64 {
	return e.boundingRadius
}

func (e *Entity) IsTargetOn() bool {
	return e.targetOn
}

// 设置目标
func (e *Entity) SetTarget(target Vector2D) {
	e.targetOn = true
	e.targetPos = target
	// fmt.Printf("SetTarget now %v %v\n", target, time.Now().Unix())
}

func (e *Entity) ClearTarget() {
	e.targetOn = false
	e.pos = e.targetPos
	e.velocity = Vector2D{0, 0}
	fmt.Printf("ClearTarget now %v\n", time.Now().Unix())
}

func (e *Entity) EnforceNonPenetrationConstraint() {
	for _, e2 := range e.world.AllEntities() {
		if e == e2 {
			continue
		}

		toEntity := e.pos.Sub(e2.pos)
		distFromEachOther := toEntity.Length()

		amountOfOverLap := (e.boundingRadius + e2.boundingRadius) - distFromEachOther
		if amountOfOverLap > 0 {
			pos2 := e.pos.Add(toEntity.DivScalar(distFromEachOther).AddScalar(amountOfOverLap))
			fmt.Printf("set pos %v %v %v\n", e.pos, pos2, amountOfOverLap)
			e.SetPos(pos2)
		}
	}
}

func (e *Entity) Update(timeDelta float64) {
	if !e.targetOn {
		return
	}

	sf, canmove := e.Calculate()

	if !canmove {
		return
	}

	// 加速度 = 力/质量
	acceleration := sf.DivScalar(e.Mass)

	// 不进行加速度计算
	// e.velocity = e.velocity.Add(acceleration.MulScalar(timeDelta)).
	// 	Truncate(e.MaxSpeed)

	e.velocity = e.velocity.Add(acceleration).
		Truncate(e.MaxSpeed)

	e.pos = e.pos.Add(e.velocity.MulScalar(timeDelta))
	// fmt.Printf("pos %v %v %v\n", e.pos, e.velocity, sf)

	if e.velocity.LengthSquared() > 0.0001 {
		e.heading = e.velocity.Normalize()
		e.side = e.heading.Perp()
	}

	if e.targetPos.Sub(e.pos).LengthSquared() < 1 {
		e.ClearTarget()
	}

	e.EnforceNonPenetrationConstraint()
}

// 计算合力
func (e *Entity) Calculate() (Vector2D, bool) {
	force := Vector2D{0, 0}
	targetForce := Vector2D{0, 0}
	if e.targetOn {
		targetForce = e.Seek(e.targetPos)
	}

	obForce := e.ObstacleAvoidance()

	// fmt.Printf("targetForce length %v obForce  %v\n",
	// 	targetForce.LengthSquared(),
	// 	obForce.LengthSquared())

	// if obForce.LengthSquared() > 1e-7 {
	// 	val := targetForce.Normalize().Dot(obForce.Normalize())
	// 	val1 := math.Acos(val)
	// 	fmt.Printf("angle %v limit %v\n", val1, (0.95 * math.Pi))
	// }

	force = targetForce.Add(obForce)
	return force.Sub(e.velocity), true
}

// 靠近
func (e *Entity) Seek(targetPos Vector2D) Vector2D {
	diredvelocity := targetPos.Sub(e.pos).Normalize().
		MulScalar(e.MaxSpeed)
	// result := diredvelocity.Sub(e.velocity)
	result := diredvelocity
	return result
}

// 计算Entity之间的阻挡
func (e *Entity) ObstacleAvoidance() Vector2D {
	e.dDBoxLength = MINDETECTIONBOXLENGTH

	steeringForce := Vector2D{0, 0}

	for _, e2 := range e.world.AllEntities() {
		rlen := e2.boundingRadius + MINDETECTIONBOXLENGTH
		if e == e2 || e2.pos.Sub(e.pos).LengthSquared() > rlen*rlen {
			continue
		}

		localPos := PointToLocalSpace(e2.pos, e.heading, e.side, e.pos)
		if localPos.X <= 0 {
			continue
		}

		sf := Vector2D{0, 0}
		multiplier := 1.0 + (e.dDBoxLength-localPos.X)/e.dDBoxLength
		sf.Y = (e2.boundingRadius - localPos.Y) * multiplier

		brakingWeight := 0.2
		sf.X = (e2.boundingRadius - localPos.X) * brakingWeight

		s1 := VectorToWorldSpace(sf, e.heading, e.side, e.pos)
		//fmt.Printf("sf %v s1 %v\n", sf, s1)
		steeringForce = steeringForce.Add(s1)
	}

	if steeringForce.LengthSquared() < 1e-7 {
		return steeringForce
	} else {
		return steeringForce.Normalize().MulScalar(e.MaxSpeed)
	}

}
