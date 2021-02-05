package electron

import "syscall/js"

type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (rectangle *Rectangle) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"x":      rectangle.X,
		"y":      rectangle.Y,
		"width":  rectangle.Width,
		"height": rectangle.Height,
	})
}

func JSValueToRectangle(value js.Value) *Rectangle {
	return &Rectangle{
		X:      value.Get("x").Int(),
		Y:      value.Get("y").Int(),
		Width:  value.Get("width").Int(),
		Height: value.Get("height").Int(),
	}
}
