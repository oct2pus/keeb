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
	PLATEHEIGHT = 1.2   //y
	PLATEWIDTH  = 94.6  //z
	// keyswitch
	SWITCHWIDTH  = 15 // square
	SWITCHLENGTH = 15 // square
	// keycap
	CAP1LENGTH  = 17.5 //x
	CAP15LENGTH = 17.5 * 1.5
	CAP2LENGTH  = 17.5 * 2
	CAPHEIGHT   = 1    //y for reference only
	CAPWIDTH    = 16.5 //z
)

func keeb() sdf.SDF3 {

	return nil
}

func main() {
	render.ToSTL(keeb(), 300, "pool1.stl", &render.MarchingCubesOctree{})
}

//-----------------------------------------------------------------------------
