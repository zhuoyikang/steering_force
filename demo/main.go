package main

import (
	"time"

	sf "github.com/indie21/steering_force"
)

func runWorld() *sf.World {
	world := sf.NewWorld(100, 100, time.Second/60)
	world.Run()

	entity := sf.NewEntity()
	entity.SetPos(sf.Vector2D{250, 250})
	// entity.SetTarget(sf.Vector2D{500, 500})
	world.AddEntity(entity)
	return world
}

func runGui(world *sf.World) {
	gui := NewGui(world, 1, 1000, 1000, time.Second/30)
	gui.Run()
}

func main() {
	closeChain := make(chan bool)
	world := runWorld()
	runGui(world)
	<-closeChain
}
