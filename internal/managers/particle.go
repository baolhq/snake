package managers

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Particle handle all particles related logic
var Particle = NewParticleSystem()

// particleData represents a single particle
type particleData struct {
	X, Y     float64
	VX, VY   float64
	Life     float64
	MaxLife  float64
	Color    color.RGBA
	Size     float64
	Image    *ebiten.Image
	Rotation float64
}

// ParticleSystem manages multiple particles
type ParticleSystem struct {
	Particles []*particleData
}

// NewParticleSystem creates an empty particle system
func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		Particles: []*particleData{},
	}
}

// Explode spawns particles at x,y using existing image
func (ps *ParticleSystem) Explode(img []byte, x, y int, count int) {
	decodedImg, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		log.Fatal(err)
	}
	particleImg := ebiten.NewImageFromImage(decodedImg)
	w, h := particleImg.Bounds().Dx(), particleImg.Bounds().Dy()

	for range count {
		px := float64(rand.Intn(w)) - float64(w)/2
		py := float64(rand.Intn(h)) - float64(h)/2

		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*1.5 + 0.5

		size := 12.0 + rand.Float64()*5
		life := rand.Float64()*1.0 + 0.5

		ps.Particles = append(ps.Particles, &particleData{
			X:        float64(x) + px*0.3,
			Y:        float64(y) + py*0.3,
			VX:       math.Cos(angle) * speed,
			VY:       math.Sin(angle) * speed,
			Life:     life,
			MaxLife:  life,
			Color:    color.RGBA{255, 255, 255, 255},
			Size:     size,
			Image:    particleImg,
			Rotation: rand.Float64() * 2 * math.Pi,
		})
	}
}

// Update moves particles, shrinks them, and fades alpha
func (ps *ParticleSystem) Update(dt float64) {
	newParticles := ps.Particles[:0]
	for _, p := range ps.Particles {
		p.X += p.VX
		p.Y += p.VY
		p.Life -= dt

		if p.Life > 0 {
			p.Size *= 0.98

			frac := p.Life / p.MaxLife
			if frac > 1 {
				frac = 1
			}
			if frac < 0 {
				frac = 0
			}
			p.Color.A = uint8(float64(255) * frac)

			newParticles = append(newParticles, p)
		}
	}
	ps.Particles = newParticles
}

// Draw renders particles using your existing image with fading alpha
func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for _, p := range ps.Particles {
		op := &ebiten.DrawImageOptions{}

		w := p.Image.Bounds().Dx()
		h := p.Image.Bounds().Dy()

		// center, rotate, scale, and translate
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Rotate(p.Rotation)
		op.GeoM.Scale(p.Size/float64(w), p.Size/float64(h))
		op.GeoM.Translate(p.X, p.Y)

		alphaF := float32(p.Color.A) / 255
		op.ColorScale.Scale(
			float32(p.Color.R)/255*alphaF,
			float32(p.Color.G)/255*alphaF,
			float32(p.Color.B)/255*alphaF,
			alphaF,
		)

		screen.DrawImage(p.Image, op)
	}
}
