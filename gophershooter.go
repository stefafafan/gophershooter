package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/ajhager/engi"
)

var (
	region        *engi.Region
	batch         *engi.Batch
	font          *engi.Font
	gopher        *Gopher
	speed         float32
	voices        []*Voice
	enemies       []*Gopher
	spawntime     int
	score         int
	scoreText     string
	tick          int
	lastShootTime int
)

type Gopher struct {
	*engi.Sprite
	Velocity   *engi.Point
	gopherType string
	goLeft     bool
}

type Voice struct {
	X, Y float32
}

type Game struct {
	*engi.Game
}

func (game *Game) Preload() {
	engi.Files.Add("gopher", "data/gopher.png")
	engi.Files.Add("gopher-red", "data/gopher-red.png")
	engi.Files.Add("gopher-blue", "data/gopher-blue.png")
	engi.Files.Add("gopher-green", "data/gopher-green.png")
	engi.Files.Add("font", "data/font.png")
}

func (game *Game) Setup() {
	engi.SetBg(0x2a2a2a)
	font = engi.NewGridFont(engi.Files.Image("font"), 20, 20)
	batch = engi.NewBatch(engi.Width(), engi.Height())
	gophertexture := engi.Files.Image("gopher")
	region = engi.NewRegion(gophertexture, 0, 0, int(gophertexture.Width()), int(gophertexture.Height()))
	gopher = &Gopher{
		engi.NewSprite(region, 0, 0),
		&engi.Point{0, 0},
		"player",
		false,
	}
	gopher.Position.X = (engi.Width() / 2) - 50
	gopher.Position.Y = engi.Height() - 100
	speed = 350
	voices = make([]*Voice, 0)
	enemies = make([]*Gopher, 0)

	game.Spawn()
	rand.Seed(time.Now().UnixNano())
	spawntime = rand.Intn(50) + 30

	score = 0
	scoreText = fmt.Sprintf("score: %d", score)

	tick = 0
	lastShootTime = 0
}

func (game *Game) Render() {
	batch.Begin()
	gopher.Render(batch)
	for _, enemy := range enemies {
		enemy.Render(batch)
	}
	for _, voice := range voices {
		font.Print(batch, "Go", voice.X, voice.Y, 0xffffff)
	}
	font.Print(batch, scoreText, 20, 60, 0xffffff)
	font.Print(batch, "Gopher Shooter", 20, 20, 0xffffff)
	batch.End()
}

func (g *Game) Key(key engi.Key, modifier engi.Modifier, action engi.Action) {
	if action == engi.PRESS {
		switch key {
		case engi.ArrowLeft:
			gopher.Velocity.X -= speed
		case engi.ArrowRight:
			gopher.Velocity.X += speed
		case engi.ArrowUp:
			gopher.Velocity.Y -= speed
		case engi.ArrowDown:
			gopher.Velocity.Y += speed
		case engi.Space:
			if tick-lastShootTime > 20 {
				g.Shoot()
				lastShootTime = tick
			}
		}
	} else if action == engi.RELEASE {
		switch key {
		case engi.ArrowLeft:
			gopher.Velocity.X += speed
		case engi.ArrowRight:
			gopher.Velocity.X -= speed
		case engi.ArrowUp:
			gopher.Velocity.Y += speed
		case engi.ArrowDown:
			gopher.Velocity.Y -= speed
		}
	}
}

func (g *Gopher) Intersects(v *Voice) bool {
	gLeft := g.Position.X - 10
	gRight := g.Position.X + 70
	gTop := g.Position.Y - 10
	gBottom := g.Position.Y + 70

	vLeft := v.X - 20
	vTop := v.Y + 20
	return (vLeft >= gLeft && vLeft <= gRight && vTop >= gTop && vTop <= gBottom)
}

func (g *Game) Shoot() {
	voice := &Voice{
		gopher.Position.X,
		gopher.Position.Y,
	}
	voices = append(voices, voice)
}

func (g *Game) Spawn() {
	var enemyTexture *engi.Texture
	var gType string
	var goLeft bool
	switch rand.Intn(3) {
	case 0:
		enemyTexture = engi.Files.Image("gopher-red")
		gType = "red"
		goLeft = false
	case 1:
		enemyTexture = engi.Files.Image("gopher-blue")
		gType = "blue"
		goLeft = false
	default:
		enemyTexture = engi.Files.Image("gopher-green")
		gType = "green"
		goLeft = rand.Intn(2) == 0
	}
	gopherRegion := engi.NewRegion(enemyTexture, 0, 0, int(enemyTexture.Width()), int(enemyTexture.Height()))
	enemy := &Gopher{
		engi.NewSprite(gopherRegion, 0, 0),
		&engi.Point{0, 0},
		gType,
		goLeft,
	}
	enemy.Position.X = float32(rand.Intn(int(engi.Width()) - 60))
	enemy.Position.Y = -50
	enemies = append(enemies, enemy)
}

func (g *Game) Update(dt float32) {
	tick += 1
	gopher.Position.X += gopher.Velocity.X * dt
	gopher.Position.Y += gopher.Velocity.Y * dt

	if gopher.Position.X > engi.Width()-60 {
		gopher.Position.X = engi.Width() - 60
	}
	if gopher.Position.X < 0 {
		gopher.Position.X = 0
	}
	if gopher.Position.Y > engi.Height()-70 {
		gopher.Position.Y = engi.Height() - 70
	}
	if gopher.Position.Y < 5 {
		gopher.Position.Y = 5
	}

	for _, voice := range voices {
		voice.Y -= 300 * dt
	}

	newEnemies := make([]*Gopher, 0)
	for _, enemy := range enemies {
		switch enemy.gopherType {
		case "red":
			enemy.Position.Y += 4
		case "blue":
			enemy.Position.X += float32(math.Sin(float64(tick/5))*5) * 5
			enemy.Position.Y += 6
		default:
			if enemy.goLeft {
				enemy.Position.X += 6
			} else {
				enemy.Position.X -= 6
			}

			enemy.Position.Y += 8
		}

		if enemy.Position.Y < engi.Height()+10 {
			newEnemies = append(newEnemies, enemy)
		}
		for _, voice := range voices {
			if enemy.Intersects(voice) {
				newEnemies = newEnemies[:len(newEnemies)-1]
				switch enemy.gopherType {
				case "red":
					score += 10
				case "blue":
					score += 20
				default:
					score += 30
				}
			}
		}
	}
	enemies = newEnemies
	scoreText = fmt.Sprintf("score: %d", score)

	spawntime--
	if spawntime < 0 {
		g.Spawn()
		spawntime = rand.Intn(50) + 30
	}
}

func (game *Game) Resize(w, h float32) {
	batch.SetProjection(w, h)
}

func main() {
	engi.Open("Gopher", 800, 450, false, &Game{})
}
