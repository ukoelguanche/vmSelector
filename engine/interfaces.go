package engine

import (
	"time"

	"github.com/ukoelguanche/graphicsengine/core"
)

type Drawer interface {
	DrawSpriteRect(sprite *core.Sprite, rect core.Rect, position core.Point)
}

type Renderable interface {
	NextFrame()
	Draw(d Drawer)

	GetSprite() *core.Sprite
	GetPosition() core.Point
	GetTargetPosition() core.Point
	GetSpeed() core.Size
	SetPosition(core.Point)
	GetMovementFrameCount() float64
	GetMovementFrame() float64
	EndMovement()
	IsMoving() bool
	SetTargetPosition(core.Point)
	SetSpeed(float64)
	GetTotalDistance() float64
	SetEaseFunction(func(float64) float64)
	GetEaseFunction() func(float64) float64
	SetOnMovementComplete(func(Renderable))
	GetStartPosition() core.Point
	GetStartTime() time.Time
	GetDuration() time.Duration
}

func UpdatePosition(r Renderable) {
	if !r.IsMoving() {
		return
	}

	// 1. Calculamos el progreso temporal (t) de 0.0 a 1.0
	// Esto NUNCA se queda en cero porque el tiempo siempre pasa.
	elapsed := time.Since(r.GetStartTime())
	duration := r.GetDuration()

	t := elapsed.Seconds() / duration.Seconds()
	if t > 1.0 {
		t = 1.0
	}

	// 2. Pasamos ese 't' por la función de Ease
	// Aquí es donde sucede la magia: t avanza lineal,
	// pero easedT avanza con curvas (lento-rápido-lento)
	easedT := t
	if easeFunc := r.GetEaseFunction(); easeFunc != nil {
		easedT = easeFunc(t)
	}

	// 3. Interpolación Lineal (LERP) usando el easedT
	start := r.GetStartPosition()
	target := r.GetTargetPosition()

	// Fórmula: inicio + (destino - inicio) * progreso_suave
	nextX := start.X + (target.X-start.X)*easedT
	nextY := start.Y + (target.Y-start.Y)*easedT

	r.SetPosition(core.Point{X: nextX, Y: nextY})

	// 4. Condición de parada: Si el tiempo se agotó
	if t >= 1.0 {
		r.SetPosition(target) // Aseguramos el píxel perfecto al final
		r.EndMovement()
	}
}
