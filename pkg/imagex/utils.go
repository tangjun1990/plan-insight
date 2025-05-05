package imagex

import (
	"image"
	"image/color"
)

// DrawLine 绘制直线 (使用Bresenham算法)
func DrawLine(img *image.RGBA, p1, p2 image.Point, c color.Color) {
	dx := Abs(p2.X - p1.X)
	dy := Abs(p2.Y - p1.Y)

	var sx, sy int
	if p1.X < p2.X {
		sx = 1
	} else {
		sx = -1
	}

	if p1.Y < p2.Y {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	x, y := p1.X, p1.Y
	for {
		img.Set(x, y, c)

		if x == p2.X && y == p2.Y {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}

		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// Abs 计算整数的绝对值
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Min 返回两个整数中的较小值
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max 返回两个整数中的较大值
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
