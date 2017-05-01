package steering_force

var (
	DR    = float64(20)
	SPEED = float64(60)
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
}

func NewEntity() *Entity {
	e := &Entity{
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

	if e.targetPos.Sub(e.pos).LengthSquared() < 0.1 {
		e.ClearTarget()
	}
}

// 计算合力
func (e *Entity) Calculate() Vector2D {
	force := Vector2D{0, 0}
	if e.targetOn {
		force = force.Add(e.Seek(e.targetPos))
	}
	return force
}

// 靠近
func (e *Entity) Seek(targetPos Vector2D) Vector2D {
	diredvelocity := targetPos.Sub(e.pos).Normalize().
		MulScalar(e.MaxSpeed)
	result := diredvelocity.Sub(e.velocity)
	return result
}
