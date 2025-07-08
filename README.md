# gocharts

[BubbleTea](https://github.com/charmbracelet/bubbletea) compatible chart package that renders charts in the terminal.

## Features

| Feature     | Status      |
| ------------- | ------------- |
| `piechart` | Complete ✅ |

## API

| Function     | Usage      |
| ------------- | ------------- |
| `piechart.New(radius int, ...opts)`| Init a new pie chart |
| `piechart.View()` | BubbleTea compatible view function |
| `piechart.Push(*PieValue)` | Push a single `PieValue` to the chart instance |
| `piechart.PushAll([]*PieValue)` | Push multiple `PieValue`s to the chart instance |
| `piechart.UpdateAnimation()` | Update animation to next frame |
| `piechart.IsAnimationComplete()` | Check if animation is completed |
| `piechart.RestartAnimation()` | Restart animation |

### Init Options
| Option     | Usage      |
| ------------- | ------------- |
| `piechart.WithData(piechart.PieData)` | to populate the piechart with data |
| `piechart.WithShowLegend(bool)` | to show or hide legend |
| `piechart.WithRadius(int)` | to modify radius |
| `piechart.WithAspectRatio(float64)` | to modify the aspect ratio, where the aspect ratio is the ratio width/height |
| `piechart.WithValuePrefix(string)` | to add a prefix to the displayed values in the legend |
| `piechart.WithAnimation(bool)` | to turn on or off the sweeping animation when rendering |
| `piechart.WithAnimationDuration(time.Duration)` | to modify the duration of the sweeping animation |

## References

- Chart interface inspired by [ntcharts](https://github.com/NimbleMarkets/ntcharts)
- Pie chart logic and rendering logic from [term-piechart](https://github.com/va-h/term-piechart)
