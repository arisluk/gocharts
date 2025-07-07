package piechart

import "fmt"

func Move(x, y int) string {
	return fmt.Sprintf("\033[%d;%dH", y, x)
}
