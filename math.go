package main

import "math/rand/v2"

type vec2 struct {
	x, y float32
}

func (v vec2) add(rhs vec2) vec2 {
	return vec2{
		x: v.x + rhs.x,
		y: v.y + rhs.y,
	}
}

func (v vec2) sub(rhs vec2) vec2 {
	return vec2{
		x: v.x - rhs.x,
		y: v.y - rhs.y,
	}
}

func (v vec2) mul(t float32) vec2 {
	return vec2{
		x: v.x * t,
		y: v.y * t,
	}
}

func clamp(t, min, max float32) float32 {
	if t < min {
		return min
	}

	if t > max {
		return max
	}

	return t
}

func randDir() float32 {
	if rand.IntN(2) == 0 {
		return -1
	}
	return 1
}

func lerp(start, end, amount float32) float32 {
	return start*(1-amount) + end*amount
}
