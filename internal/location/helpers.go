package location

import "github.com/stratoberry/go-gpsd"

type GPS_STATUS uint8

const (
	WAIT_SKY GPS_STATUS = iota
	WAIT_FIX
	FIX
)

func CalculateState(sky *gpsd.SKYReport, tpv *gpsd.TPVReport) GPS_STATUS {
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
