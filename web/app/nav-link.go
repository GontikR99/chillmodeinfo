package main

import (
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/vugu/vugu"
)

type NavLink struct {
	Icon        string
	DisplayName string
	Place string
}

func (c *NavLink) classText() string {
	if place.GetPlace()==c.Place {
		return "nav-link active"
	} else {
		return "nav-link"
	}
}

func (c *NavLink) follow(event vugu.DOMEvent) {
	event.PreventDefault()
	place.NavigateTo(event.EventEnv(), c.Place)
}