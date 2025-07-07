package piechart

type Option func(*Model)

func WithShowLegend(showLegend bool) Option {
	return func(m *Model) {
		m.showLegend = showLegend
	}
}

func WithRadius(radius int) Option {
	return func(m *Model) {
		m.radius = radius
	}
}

func WithAspectRatio(aspectRatio int) Option {
	return func(m *Model) {
		m.aspectRatio = aspectRatio
		m.centerX = m.radius * int(aspectRatio)
	}
}

func WithData(data PieData) Option {
	return func(m *Model) {
		m.PushAll(data.Values)
	}
}
