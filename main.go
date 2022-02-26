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
	return rect(0, 0, PLATELENGTH, PLATEHEIGHT, PLATEWIDTH)
}

func keyCaps() []sdf.V2 {
	// this is about to get a lil silly
	sideWidthGaps := gapLength(PLATEWIDTH, (SWITCHWIDTH * 5), 5+1)                              // 5 rows, so 6 gaps including exterior, so we add 1
	row1Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP15LENGTH)+(12*CAP1LENGTH), 14-1) // we want consistent gaps on the outside, so remove those, we don't include exterior at all so we subtract 1.
	row2Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP15LENGTH)+(12*CAP1LENGTH), 14-1) // split for unneccessary futureproofing
	row3Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP2LENGTH)+(11*CAP1LENGTH), 13-1)
	row4Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP2LENGTH)+(11*CAP1LENGTH), 13-1)
	row5Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (3*CAP2LENGTH)+(1*CAP15LENGTH)+(7*CAP1LENGTH)+(0.5*CAP1LENGTH), 11) // this one has an additional gapLength to split the arrow keys

	row1 := createRow(sideWidthGaps, sideWidthGaps, CAPHEIGHT, row1Gaps, CAP15LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH)
	row2 := createRow(sideWidthGaps, (sideWidthGaps*2)+(CAPWIDTH*1), CAPHEIGHT, row2Gaps, CAP15LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH)
	row3 := createRow(sideWidthGaps, (sideWidthGaps*3)+(CAPWIDTH*2), CAPHEIGHT, row3Gaps, CAP2LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP2LENGTH)
	row4 := createRow(sideWidthGaps, (sideWidthGaps*4)+(CAPWIDTH*3), CAPHEIGHT, row4Gaps, CAP2LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP2LENGTH)
	row5 := createRow(sideWidthGaps, (sideWidthGaps*5)+(CAPWIDTH*4), CAPHEIGHT, row5Gaps, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH, CAP1LENGTH, CAP2LENGTH, CAP2LENGTH, CAP2LENGTH, CAP1LENGTH/2, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH)
	dimensions := append(row1, row2...)
	dimensions = append(dimensions, row3...)
	dimensions = append(dimensions, row4...)
	dimensions = append(dimensions, row5...)
	return dimensions
}

func drawKeyCaps() (sdf.SDF3, error) {
	p := sdf.NewPolygon()
	p.AddV2Set(keyCaps())
	profile, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(profile, CAPHEIGHT), err
}

func main() {
	plate, err := plate()
	if err != nil {
		log.Fatal("error: %s\n", err)
	}
	keycaps, err := drawKeyCaps()
	if err != nil {
		log.Fatal("error: %s\n", err)
	}
	render.ToSTL(plate, 300, "keeb.stl", &render.MarchingCubesUniform{})
	render.ToSTL(keycaps, 300, "keycaps.stl", &render.MarchingCubesOctree{})
}

func rect(x1, z1, x2, y, z2 float64) (sdf.SDF3, error) {
	dimensions := []sdf.V2{
		{x1, z1},
		{x2, z1},
		{x2, z2},
		{x1, z2},
	}
	p := sdf.NewPolygon()
	p.AddV2Set(dimensions)
	profile, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(profile, y), err
}

func gapLength(plateLength float64, keyLength float64, gaps float64) float64 {
	return (plateLength - keyLength) / gaps
}

func createRow(inX, inZ, y, gap float64, keys ...float64) []sdf.V2 {
	dimensions := make([]sdf.V2, 0, 0)
	x, z := inX, inZ
	for _, ele := range keys {
		dimensions = append(dimensions, sdf.V2{x, z})
		dimensions = append(dimensions, sdf.V2{x + ele, z})
		dimensions = append(dimensions, sdf.V2{x + ele, z + CAPWIDTH})
		dimensions = append(dimensions, sdf.V2{x, z + CAPWIDTH})
		x = x + gap + ele
	}
	return dimensions
}

//-----------------------------------------------------------------------------
