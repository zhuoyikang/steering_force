package steering_force

import (
	"fmt"
	"math"
)

var (
	DR                    = float64(20)
	SPEED                 = float64(60)
	MINDETECTIONBOXLENGTH = float64(40)
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
	// fmt.Printf("ClearTarget now %v\n", time.Now().Unix())
}

func (e *Entity) Update(timeDelta float64) {
	if !e.targetOn {
		return
	}

	sf := e.Calculate()

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

	if e.targetPos.Sub(e.pos).LengthSquared() < 0.3 {
		e.ClearTarget()
	}
}

// 计算合力
func (e *Entity) Calculate() Vector2D {
	force := Vector2D{0, 0}
	if e.targetOn {
		force = force.Add(e.Seek(e.targetPos))
	}

	obforce := e.ObstacleAvoidance()
	fmt.Printf("obforce %v\n", obforce)

	force = force.Add(obforce)

	return force
}

// 靠近
func (e *Entity) Seek(targetPos Vector2D) Vector2D {
	diredvelocity := targetPos.Sub(e.pos).Normalize().
		MulScalar(e.MaxSpeed)
	result := diredvelocity.Sub(e.velocity)
	return result
}

// 计算Entity之间的阻挡
func (e *Entity) ObstacleAvoidance() Vector2D {
	e.dDBoxLength = 2 * MINDETECTIONBOXLENGTH

	var closestIntersectingObstacle *Entity
	distToClosestIP := math.MaxFloat64
	var localPosOfClosestObstacle Vector2D
	steeringForce := Vector2D{0, 0}

	for _, e2 := range e.world.AllEntities() {
		rlen := e2.boundingRadius + MINDETECTIONBOXLENGTH
		fmt.Printf("e2.pos.Sub(e.pos).LengthSquared() %v %v\n", e2.pos.Sub(e.pos).LengthSquared(), rlen*rlen)
		if e == e2 || e2.pos.Sub(e.pos).LengthSquared() > rlen*rlen {
			continue
		}

		localPos := PointToLocalSpace(e2.pos, e.heading, e.side, e.pos)
		if localPos.X <= 0 {
			continue
		}

		expandedRadius := e.boundingRadius + e2.boundingRadius
		if math.Abs(localPos.Y) >= expandedRadius {
			continue
		}

		cX := localPos.X
		cY := localPos.Y

		sqrtPart := math.Sqrt(expandedRadius*expandedRadius - cY*cY)
		ip := cX - sqrtPart

		if ip <= 0.0 {
			ip = cX + sqrtPart
		}

		if ip < distToClosestIP {
			distToClosestIP = ip
			closestIntersectingObstacle = e2
			localPosOfClosestObstacle = localPos
		}

		multiplier := 3.0 + (e.dDBoxLength-localPosOfClosestObstacle.X)/e.dDBoxLength

		steeringForce.Y = (closestIntersectingObstacle.boundingRadius -
			localPosOfClosestObstacle.Y) * multiplier

		brakingWeight := 0.0

		steeringForce.X = (closestIntersectingObstacle.boundingRadius -
			localPosOfClosestObstacle.X) *
			brakingWeight

		fmt.Printf("obforce2 %v\n", steeringForce)

		s1 := VectorToWorldSpace(steeringForce, e.heading, e.side, e.pos)
		steeringForce = steeringForce.Add(s1)
	}

	fmt.Printf("obforce1 \n")

	return steeringForce
}
