package piechart

import (
	"math"
	"time"
)

type Option func(*Model)

func WithShowLegend(showLegend bool) Option {
	return func(m *Model) {
		m.showLegend = showLegend
	}
}

func WithAnimation(withAnimation bool) Option {
	return func(m *Model) {
		m.withAnimation = withAnimation
	}
}

func WithAnimationDuration(animationDuration time.Duration) Option {
	return func(m *Model) {
		m.animationDuration = animationDuration
	}
}

func WithRadius(radius int) Option {
	return func(m *Model) {
		m.radius = radius
	}
}

func WithAspectRatio(aspectRatio float64) Option {
	return func(m *Model) {
		m.aspectRatio = aspectRatio
		m.centerX = int(math.Round(float64(m.radius) * aspectRatio))
	}
}

func WithValuePrefix(valuePrefix string) Option {
	return func(m *Model) {
		m.valuePrefix = valuePrefix
	}
}

func WithData(data PieData) Option {
	return func(m *Model) {
		m.PushAll(data.Values)
	}
}
