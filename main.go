package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// in millimeters

const (
	// plate
	PLATELENGTH = 285.0 //x
	PLATEHEIGHT = 1.2 //y
	PLATEWIDTH = 94.6 //z
	// keyswitch
	SWITCHWIDTH = 15 // square
	SWITCHLENGTH = 15 // square
	// keycap
	CAPLENGTH = 
)


func keeb() sdf.SDF3 {

	return nil
}

func main() {
	render.ToSTL(keeb(), 300, "pool1.stl", &render.MarchingCubesOctree{})
}

//-----------------------------------------------------------------------------
