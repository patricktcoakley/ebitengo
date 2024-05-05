package main

type paddle struct {
	pos vec2
	vel vec2
}

func newPaddle(x, y float32) paddle {
	return paddle{
		pos: vec2{x: x, y: y},
		vel: vec2{x: 0, y: paddleWidth},
	}
}

func (p paddle) didHit(b ball) bool {
	isBallInFront := b.pos.x < p.pos.x+paddleWidth
	isBallBehind := b.pos.x+ballWidth > p.pos.x
	isBallAbove := b.pos.y+ballWidth > p.pos.y
	isBallBelow := b.pos.y < p.pos.y+paddleHeight

	return isBallInFront && isBallBehind && isBallBelow && isBallAbove
}
