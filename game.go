package main

import (
	"bytes"
	"ebitenpong/assets/fonts"
	"ebitenpong/assets/sounds"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand/v2"
)

const (
	screenWidth  = 1200
	screenHeight = 800

	paddleWidth  = 18
	paddleHeight = 80

	ballSpeed = 12
	ballWidth = 18
)

var mainFont font.Face

type game struct {
	player      paddle
	cpu         paddle
	ball        ball
	score       struct{ player, cpu int }
	audioPlayer *audio.Player
}

func newGame() *game {
	g := game{}
	g.player = newPaddle(100, screenHeight/2-paddleHeight)
	g.cpu = newPaddle(screenWidth-paddleWidth-100, screenHeight/2-paddleHeight)
	g.ball = newBall(randDir())

	ttf, err := opentype.Parse(fonts.ArcadeFont)
	if err != nil {
		log.Fatal(err)
	}

	mainFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    16,
		DPI:     300,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	beep, err := wav.DecodeWithSampleRate(48000, bytes.NewReader(sounds.BeepWAV))
	if err != nil {
		log.Fatal(err)
	}

	audioContext := audio.NewContext(48000)

	audioPlayer, err := audioContext.NewPlayer(beep)
	if err != nil {
		log.Fatal(err)
	}

	g.audioPlayer = audioPlayer

	return &g
}

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.pos = g.player.pos.sub(g.player.vel)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.pos = g.player.pos.add(g.player.vel)
	} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	} else if ebiten.IsKeyPressed(ebiten.KeyF1) {
		g.reset()
		return nil
	}

	g.player.pos.y = clamp(g.player.pos.y, paddleWidth, screenHeight-paddleHeight-paddleWidth)

	g.cpu.pos.y = lerp(g.cpu.pos.y, g.ball.pos.y, float32(rand.NormFloat64()))
	g.cpu.pos.y = clamp(g.cpu.pos.y, paddleWidth, screenHeight-paddleHeight-paddleWidth)

	if g.ball.pos.y <= paddleWidth || g.ball.pos.y >= screenHeight-paddleWidth-ballWidth {
		g.ball.vel.y *= -1
	}

	g.ball.pos = g.ball.pos.add(g.ball.vel.mul(ballSpeed))

	if g.ball.pos.x <= 0 {
		g.score.cpu += 1
		g.ball = newBall(-1)
	} else if g.ball.pos.x >= screenWidth {
		g.score.player += 1
		g.ball = newBall(1)
	} else if g.player.didHit(g.ball) || g.cpu.didHit(g.ball) {
		g.playBeep()
		g.ball.vel = vec2{x: g.ball.vel.x * -1, y: randDir()}
	}

	if g.score.player >= 10 || g.score.cpu >= 10 {
		g.reset()
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	score := fmt.Sprintf("%d %d", g.score.player, g.score.cpu)
	text.Draw(screen, score, mainFont, screenWidth/2-90, 80, color.White)

	for i := range screenHeight {
		if paddleWidth-i%50 == 0 {
			vector.DrawFilledRect(screen, screenWidth/2, float32(i), paddleWidth, 20, color.White, false)
		}
	}

	for i := range screenWidth {
		vector.DrawFilledRect(screen, float32(i), 0, paddleWidth, 20, color.White, false)
		vector.DrawFilledRect(screen, float32(i), screenHeight-paddleWidth-1, paddleWidth, paddleWidth, color.White, false)
	}

	vector.DrawFilledRect(screen, g.ball.pos.x, g.ball.pos.y, ballWidth, ballWidth, color.White, false)
	vector.DrawFilledRect(screen, g.player.pos.x, g.player.pos.y, paddleWidth, paddleHeight, color.White, false)
	vector.DrawFilledRect(screen, g.cpu.pos.x, g.cpu.pos.y, paddleWidth, paddleHeight, color.White, false)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *game) playBeep() {
	err := g.audioPlayer.Rewind()
	if err != nil {
		fmt.Println(err)
		return
	}

	g.audioPlayer.Play()
}

func (g *game) reset() {
	g.player = paddle{pos: vec2{x: 100, y: screenHeight/2 - paddleHeight}}
	g.cpu = paddle{pos: vec2{x: screenWidth - paddleWidth - 100, y: screenHeight/2 - paddleHeight}}
	g.ball = newBall(randDir())
	g.score.player, g.score.cpu = 0, 0
}
