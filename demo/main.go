package main

import (
	"time"

	sf "github.com/indie21/steering_force"
)

func runWorld() {
	world := sf.NewWorld(100, 100, time.Second/60)
	world.Run()

	entity := sf.NewEntity()
	entity.SetTarget(sf.Vector2D{0, 60})
	world.AddEntity(entity)
}

func runGui() {
	gui := NewGui(1, 1000, 1000, time.Second/2)
	gui.Run()
}

func main() {
	closeChain := make(chan bool)
	runWorld()
	runGui()
	<-closeChain
}
