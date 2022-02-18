package main

import (
	"log"

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

func plate() (sdf.SDF3, error) {
	return rect(PLATELENGTH, PLATEHEIGHT, PLATEWIDTH)
}

func main() {
	plate, err := plate()
	if err != nil {
		log.Fatal("error: %s", err)
	}
	render.ToSTL(plate, 300, "keeb.stl", &render.MarchingCubesOctree{})
}

func rect(x, y, z float64) (sdf.SDF3, error) {
	dimensions := []sdf.V2{
		{0.0, 0.0},
		{x, 0.0},
		{x, z},
		{0.0, z},
	}
	p := sdf.NewPolygon()
	p.AddV2Set(dimensions)
	profile, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(profile, y), err
}

//-----------------------------------------------------------------------------
