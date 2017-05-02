package steering_force

import (
	"time"
)

type World struct {
	entities []*Entity
	maxX     float64
	maxY     float64
	tick     time.Duration
	preTick  int64
}

func NewWorld(MaxX, MaxY float64, tick time.Duration) *World {
	w := &World{
		maxX:     MaxX,
		maxY:     MaxY,
		entities: make([]*Entity, 0),
		tick:     tick,
	}
	return w
}

// 每一帧执行
func (w *World) update(timeDelta float64) {
	for _, v := range w.entities {
		v.Update(timeDelta)
	}
}

func (w *World) run() {
	c := time.Tick(w.tick)
	w.preTick = time.Now().UnixNano()
	for now := range c {
		nowTick := now.UnixNano()
		delta := float64(nowTick-w.preTick) / float64(time.Second)
		// fmt.Printf("deltal %v\n", delta)
		w.update(delta)
		w.preTick = nowTick
	}
}

func (w *World) Run() {
	go w.run()
}

func (w *World) AddEntity(e *Entity) {
	w.entities = append(w.entities, e)
	e.world = w
}

func (w *World) AllEntities() []*Entity {
	list := make([]*Entity, 0)
	for _, v := range w.entities {
		if !v.delete {
			list = append(list, v)
		}
	}
	return list
}

func (w *World) PosConflict(pos Vector2D) *Entity {
	for _, e := range w.entities {
		if e.delete == false && e.pos.Sub(pos).Length() < (e.boundingRadius+DR) {
			return e
		}
	}
	return nil
}
