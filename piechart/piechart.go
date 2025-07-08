package piechart

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const POINT_SYMBOL = "•"
const LEGEND_PADDING = 3
const PERCENT_WIDTH = 4

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
	showLegend    bool
	withAnimation bool
	radius        int
	centerX       int
	centerY       int
	aspectRatio   float64
	valuePrefix   string
	sum           float64
	data          *PieData

	// Animation state
	sweepAngle         float64
	animationStartTime time.Time
	animationDuration  time.Duration
}

// New creates a new pie chart model.
func New(radius int, opts ...Option) Model {
	m := Model{
		showLegend:    true,
		withAnimation: false,
		radius:        radius,
		centerX:       radius * 2,
		centerY:       radius,
		aspectRatio:   2.0,
		valuePrefix:   "",
		sum:           0,
		data:          &PieData{},

		// Initialize animation state
		sweepAngle:         0,
		animationStartTime: time.Now(),
		animationDuration:  500 * time.Millisecond,
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
func (m *Model) populateAngles() {
	startAngle := 0.0
	for i, v := range m.data.Values {
		angle := v.Value / m.sum * 360
		v.Angle = startAngle + angle
		startAngle += angle
		m.data.Values[i] = v
	}
}

// Select the item from the pie chart based on the angle, respecting sweep animation
func (m *Model) selectItemFromAngle(angle float64) *PieValue {
	for _, v := range m.getVisibleSegments() {
		if (360/2 - angle) <= v.Angle {
			return v
		}
	}
	return nil
}

// getVisibleSegments returns only the segments that should be visible based on animation
func (m *Model) getVisibleSegments() []*PieValue {
	if !m.withAnimation {
		return m.data.Values
	}

	visibleSegments := make([]*PieValue, 0)

	for _, v := range m.data.Values {
		if v.Angle <= m.sweepAngle {
			visibleSegments = append(visibleSegments, v)
		} else {
			if len(visibleSegments) == 0 || visibleSegments[len(visibleSegments)-1].Angle < m.sweepAngle {
				prevAngle := 0.0
				if len(visibleSegments) > 0 {
					prevAngle = visibleSegments[len(visibleSegments)-1].Angle
				}

				if m.sweepAngle > prevAngle {
					partialSegment := &PieValue{
						Name:  v.Name,
						Color: v.Color,
						Value: v.Value * (m.sweepAngle - prevAngle) / (v.Angle - prevAngle),
						Angle: m.sweepAngle,
					}
					visibleSegments = append(visibleSegments, partialSegment)
				}
			}
			break
		}
	}

	return visibleSegments
}

// Render
func (m *Model) View() string {
	var sb strings.Builder

	m.populateAngles()

	labelIndex := 0
	legendPadding := math.Ceil(float64((m.radius*2 + 1 - len(m.data.Values))) / 2.0)
	legendBoundaryStart := -m.radius + int(legendPadding)
	legendBoundaryEnd := m.radius - int(legendPadding)

	for y := -m.radius; y < m.radius+1; y++ {
		width := int(math.Round(math.Sqrt(float64(m.radius*m.radius)-float64(y)*float64(y)) * m.aspectRatio))
		if width == 0 && m.aspectRatio != 1.0 {
			width = int(math.Round(float64(m.radius) / m.aspectRatio))
		}

		padding := m.centerX - width
		if padding < 0 {
			padding = -padding
		}
		sb.WriteString(strings.Repeat(" ", padding))

		for x := -width; x < width+1; x++ {
			angle := math.Atan2(float64(x), float64(y)) * (180 / math.Pi)

			item := m.selectItemFromAngle(angle)
			if item != nil {
				sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(item.Color)).Render(POINT_SYMBOL))
			} else {
				sb.WriteString(" ")
			}
		}

		if m.showLegend {
			if y >= legendBoundaryStart && y <= legendBoundaryEnd && labelIndex < len(m.data.Values) {
				maxNameLen := 0
				for _, v := range m.data.Values {
					if len(v.Name) > maxNameLen {
						maxNameLen = len(v.Name)
					}
				}
				valueWidth := 0
				for _, v := range m.data.Values {
					if len(fmt.Sprintf("%s%.2f", m.valuePrefix, v.Value)) > valueWidth {
						valueWidth = len(fmt.Sprintf("%s%.2f", m.valuePrefix, v.Value))
					}
				}

				name := fmt.Sprintf("%-*s", maxNameLen, m.data.Values[labelIndex].Name)
				percent := fmt.Sprintf("%*.0f%%", PERCENT_WIDTH-1, m.data.Values[labelIndex].Value/m.sum*100)
				value := fmt.Sprintf("[%s%*s]", m.valuePrefix, valueWidth, fmt.Sprintf("%.2f", m.data.Values[labelIndex].Value))

				sb.WriteString(strings.Repeat(" ", m.centerX-width+LEGEND_PADDING))
				sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(m.data.Values[labelIndex].Color)).Render(POINT_SYMBOL))
				sb.WriteString(" " + name + " " + percent + " " + value)
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

// UpdateAnimation updates the animation state using linear progression
func (m *Model) UpdateAnimation() {
	if !m.withAnimation {
		m.sweepAngle = 360
		return
	}

	elapsed := time.Since(m.animationStartTime)
	progress := float64(elapsed) / float64(m.animationDuration)

	if progress >= 1.0 {
		progress = 1.0
		m.sweepAngle = 360.0
		return
	}

	uniformAngleRad := progress * 2.0 * math.Pi
	x := 1 / m.aspectRatio * math.Cos(uniformAngleRad)
	y := math.Sin(uniformAngleRad)
	correctedAngleRad := math.Atan2(y, x)
	correctedAngleDeg := correctedAngleRad * 180.0 / math.Pi

	if correctedAngleDeg < 0 {
		correctedAngleDeg += 360
	}

	m.sweepAngle = correctedAngleDeg
}

// IsAnimationComplete returns true if the animation has finished
func (m *Model) IsAnimationComplete() bool {
	if !m.withAnimation {
		return true
	}
	return time.Since(m.animationStartTime) >= m.animationDuration
}

// RestartAnimation resets the animation to the beginning
func (m *Model) RestartAnimation() {
	if m.withAnimation {
		m.sweepAngle = 0
		m.animationStartTime = time.Now()
	}
}
