package piechart

import (
	"math"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const POINT_SYMBOL = "•"
const LEGEND_PADDING = 3

// PieValue represents a value in the pie chart.
type PieValue struct {
	Name  string
	Color string
	Value float64
	Angle float64
}

// PieData contains the label and values for the pie chart.
type PieData struct {
	Label  string
	Values []*PieValue
}

// Model represents the model for the pie chart.
type Model struct {
	showLegend  bool
	radius      int
	centerX     int
	centerY     int
	aspectRatio int
	sum         float64
	data        *PieData
}

// New creates a new pie chart model.
func New(radius int, opts ...Option) Model {
	m := Model{
		showLegend:  true,
		radius:      radius,
		centerX:     radius * 2,
		centerY:     radius,
		aspectRatio: 2,
		sum:         0,
		data:        &PieData{},
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// Push a single value to the pie chart.
func (m *Model) Push(val *PieValue) {
	val.Value = math.Max(0, val.Value)
	m.data.Values = append(m.data.Values, val)
	m.sum += val.Value
}

// Push multiple values to the pie chart
func (m *Model) PushAll(data []*PieValue) {
	for _, d := range data {
		m.Push(d)
	}

	// Sort the values by value in descending order
	sort.Slice(m.data.Values, func(i, j int) bool {
		return m.data.Values[i].Value > m.data.Values[j].Value
	})
}

// Populate the angles for the pie chart
func (m *Model) PopulateAngles() {
	startAngle := 0.0
	for i, v := range m.data.Values {
		angle := v.Value / m.sum * 360
		v.Angle = startAngle + angle
		startAngle += angle
		m.data.Values[i] = v
	}
}

// Select the item from the pie chart based on the angle
func (m *Model) SelectItemFromAngle(angle float64) *PieValue {
	for _, v := range m.data.Values {
		if (360/2 - angle) <= v.Angle {
			return v
		}
	}
	return nil
}

// Render
func (m *Model) View() string {
	var sb strings.Builder

	m.PopulateAngles()

	labelIndex := 0
	legendPadding := math.Ceil(float64((m.radius*2 + 1 - len(m.data.Values))) / 2.0)
	legendBoundaryStart := -m.radius + int(legendPadding)
	legendBoundaryEnd := m.radius - int(legendPadding)

	for y := -m.radius; y < m.radius+1; y++ {
		width := int(math.Round(math.Sqrt(float64(m.radius*m.radius)-float64(y)*float64(y)) * float64(m.aspectRatio)))
		if width == 0 && m.aspectRatio != 1 {
			width = int(math.Round(float64(m.radius) / float64(m.aspectRatio)))
		}

		padding := m.centerX - width
		if padding < 0 {
			padding = -padding
		}
		sb.WriteString(strings.Repeat(" ", padding))

		for x := -width; x < width+1; x++ {
			angle := math.Atan2(float64(x), float64(y)) * (180 / math.Pi)

			item := m.SelectItemFromAngle(angle)
			if item != nil {
				sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(item.Color)).Render(POINT_SYMBOL))
			} else {
				sb.WriteString(" ")
			}
		}

		if m.showLegend {
			if y >= legendBoundaryStart && y <= legendBoundaryEnd && labelIndex < len(m.data.Values) {
				sb.WriteString(strings.Repeat(" ", m.centerX-width+LEGEND_PADDING))
				sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(m.data.Values[labelIndex].Color)).Render(POINT_SYMBOL))
				sb.WriteString(" " + m.data.Values[labelIndex].Name)
				sb.WriteString(strings.Repeat(" ", LEGEND_PADDING))
				labelIndex++
			}
		}

		if y != m.radius {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
