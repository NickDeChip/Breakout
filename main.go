package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	winWidth  = 800
	winHeight = 450
	cols      = 24
	rows      = 8
)

type state struct {
	score    int
	gameOver bool
}

type paddle struct {
	x       float32
	y       int32
	width   int32
	height  int32
	speed   float32
	dirLeft bool
	still   bool
}

type block struct {
	x          int32
	y          int32
	width      int32
	height     int32
	colour     rl.Color
	pointValue int32
}

type ball struct {
	x       float32
	y       float32
	radius  float32
	velX    float32
	velY    float32
	dirLeft bool
}

func main() {
	rl.InitWindow(winWidth, winHeight, "Breakout")
	rl.SetTargetFPS(60)

	state := state{
		score:    0,
		gameOver: false,
	}

	paddle := paddle{
		x:       winWidth/2 - 40,
		y:       winHeight - 40,
		width:   80,
		height:  10,
		speed:   150,
		dirLeft: false,
		still:   false,
	}
	var blocks [][]block = make([][]block, cols)
	for i := range blocks {
		blocks[i] = make([]block, rows)
	}

	for i := range blocks {
		for j := range blocks[i] {
			blocks[i][j].width = 30
			blocks[i][j].height = 10
			blocks[i][j].x = 17 + int32(i)*32
			blocks[i][j].y = 10 + int32(j)*12

			if j == 0 || j == 1 {
				blocks[i][j].colour = rl.Red
				blocks[i][j].pointValue = 40
			} else if j == 2 || j == 3 {
				blocks[i][j].colour = rl.Orange
				blocks[i][j].pointValue = 30
			} else if j == 4 || j == 5 {
				blocks[i][j].colour = rl.Green
				blocks[i][j].pointValue = 20
			} else if j == 6 || j == 7 {
				blocks[i][j].colour = rl.Yellow
				blocks[i][j].pointValue = 10
			}
		}
	}

	ball := ball{
		x:       winWidth / 2,
		y:       winHeight / 2,
		radius:  5,
		velX:    75,
		velY:    150,
		dirLeft: false,
	}

	for !rl.WindowShouldClose() {

		dt := rl.GetFrameTime()
		update(&paddle, &ball, blocks, &state, dt)

		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)
		draw(&paddle, blocks, &ball, &state)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func update(paddle *paddle, ball *ball, blocks [][]block, state *state, dt float32) {
	if rl.IsKeyPressed(rl.KeyR) {
		restart(blocks, state, ball, paddle)
	}

	if state.gameOver {
		return
	}

	paddleUpdate(paddle, dt)
	ballUpdate(ball, paddle, blocks, dt, state)
}

func paddleUpdate(paddle *paddle, dt float32) {
	if rl.IsKeyDown(rl.KeyA) {
		paddle.x -= paddle.speed * dt
		paddle.dirLeft = true
		paddle.still = false
	} else {
		paddle.dirLeft = false
		paddle.still = true
	}

	if rl.IsKeyDown(rl.KeyD) {
		paddle.x += paddle.speed * dt
		paddle.dirLeft = false
		paddle.still = false
	} else if paddle.still {
		paddle.still = true
	}

	if paddle.x <= 0 {
		paddle.x = 0
	}
	if paddle.x >= winWidth-float32(paddle.width) {
		paddle.x = winWidth - float32(paddle.width)
	}
}

func ballUpdate(ball *ball, paddle *paddle, blocks [][]block, dt float32, state *state) {
	ball.x += ball.velX * dt
	ball.y += ball.velY * dt

	if ball.velX > 0 {
		ball.dirLeft = false
	} else if ball.velX < 0 {
		ball.dirLeft = true
	}

	if rl.CheckCollisionCircleRec(rl.NewVector2(ball.x, ball.y), ball.radius, rl.NewRectangle(paddle.x, float32(paddle.y), float32(paddle.width), float32(paddle.height))) {
		ball.y = float32(paddle.y) - ball.radius
		ball.velY *= -1
		if paddle.still {
			return
		}
		if paddle.dirLeft && !ball.dirLeft {
			ball.velX *= -1
		}
		if !paddle.dirLeft && ball.dirLeft {
			ball.velX *= -1
		}
	}

	if ball.x < ball.radius {
		ball.x = ball.radius
		ball.velX *= -1
	}
	if ball.x > winWidth-ball.radius {
		ball.x = winWidth - ball.radius
		ball.velX *= -1
	}
	if ball.y > winHeight-ball.radius {
		state.gameOver = true
	}
	if ball.y < ball.radius {
		ball.y = ball.radius
		ball.velY *= -1
	}

	for i := range blocks {
		for j := range blocks[i] {
			if rl.CheckCollisionCircleRec(rl.NewVector2(ball.x, ball.y), ball.radius, rl.NewRectangle(float32(blocks[i][j].x), float32(blocks[i][j].y), float32(blocks[i][j].width), float32(blocks[i][j].height))) {
				if ball.x > float32(blocks[i][j].x) {
					ball.velX *= -1
				}
				if ball.x < float32(blocks[i][j].x+blocks[i][j].width) {
					ball.velX *= -1
				}
				ball.velY *= -1
				state.score += int(blocks[i][j].pointValue)
				blocks[i][j].x = -100
				blocks[i][j].y = -100

				counter := 0
				if blocks[i][j].x > 0 {
					counter += 1
					if counter == 0 {
						nextStage(blocks, state, ball)
					} else {
						break
					}
				}
			}
		}
	}
}

func draw(paddle *paddle, blocks [][]block, ball *ball, state *state) {
	for i := range blocks {
		for j := range blocks[i] {
			rl.DrawRectangle(blocks[i][j].x, blocks[i][j].y, blocks[i][j].width, blocks[i][j].height, blocks[i][j].colour)
		}
	}
	rl.DrawRectangle(int32(paddle.x), paddle.y, paddle.width, paddle.height, rl.Blue)
	rl.DrawCircle(int32(ball.x), int32(ball.y), ball.radius, rl.White)
	rl.DrawText(fmt.Sprintf("score: %d", state.score), 10, winHeight-29, 20, rl.White)
}

func nextStage(blocks [][]block, state *state, ball *ball) {
	for i := range blocks {
		for j := range blocks[i] {
			blocks[i][j].x = 17 + int32(i)*32
			blocks[i][j].y = 10 + int32(j)*12
		}
	}
	state.score += 5000
	ball.velX += 40
	ball.velY += 40
}

func restart(blocks [][]block, state *state, ball *ball, paddle *paddle) {
	for i := range blocks {
		for j := range blocks[i] {
			blocks[i][j].x = 17 + int32(i)*32
			blocks[i][j].y = 10 + int32(j)*12
		}
	}

	state.gameOver = false
	state.score = 0

	ball.x = winWidth / 2
	ball.y = winHeight / 2
	ball.velX = 75
	ball.velY = 150

	paddle.x = float32(winWidth/2 - paddle.width/2)
}
