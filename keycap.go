package main

import "github.com/deadsy/sdfx/sdf"

const (
	CAP1LENGTH  = 17.5 //x
	CAP15LENGTH = 17.5 * 1.5
	CAP2LENGTH  = 17.5 * 2
	CAPHEIGHT   = 1    //y for reference only
	CAPWIDTH    = 16.5 //z
)

type keycap struct {
	x1, x2, z1, z2 float64
}

func (kc keycap) Dimensions() []sdf.V2 {
	return []sdf.V2{
		{kc.x1, kc.z1},
		{kc.x2, kc.z1},
		{kc.x2, kc.z2},
		{kc.x1, kc.z2},
	}
}
