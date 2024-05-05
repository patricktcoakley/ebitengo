package main

type ball struct {
	pos vec2
	vel vec2
}

func newBall(dx float32) ball {
	b := ball{
		pos: vec2{x: screenWidth/2 - ballWidth, y: screenHeight/2 - ballWidth},
		vel: vec2{x: dx},
	}
	return b
}
