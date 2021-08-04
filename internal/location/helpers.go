package location

import "github.com/stratoberry/go-gpsd"

type GPSState uint8

const (
	WAIT_SKY GPSState = iota
	WAIT_FIX
	FIX
)

func CalculateState(sky *gpsd.SKYReport, tpv *gpsd.TPVReport) GPSState {
	haveSkyReport := sky != nil
	have3DFix := tpv.Mode == 3

	switch {
	case !haveSkyReport:
		return WAIT_SKY
	case haveSkyReport && !have3DFix:
		return WAIT_FIX
	case haveSkyReport && have3DFix:
		return FIX
	default:
		return 255
	}
}
