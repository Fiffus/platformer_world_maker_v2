package attributes

import "image/color"

type Vector struct {
	X float64
	Y float64
}

type Color struct {
	Current   color.RGBA
	Normal    color.RGBA
	Highlight color.RGBA
}
