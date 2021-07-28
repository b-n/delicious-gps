package location

import "github.com/stratoberry/go-gpsd"

type GPS_STATUS uint8

const (
	WAIT_SKY GPS_STATUS = iota
	WAIT_FIX
	FIX_WEAK
	FIX_GOOD
)

func CalculateState(sky *gpsd.SKYReport, tpv *gpsd.TPVReport) GPS_STATUS {
	haveSkyReport := sky != nil
	have3DFix := tpv.Mode == 3
	sats := 0
	if haveSkyReport {
		sats = len(sky.Satellites)
	}

	switch {
	case !haveSkyReport:
		return WAIT_SKY
	case haveSkyReport && !have3DFix:
		return WAIT_FIX
	case haveSkyReport && have3DFix && sats <= 6:
		return FIX_WEAK
	case haveSkyReport && have3DFix && sats > 6:
		return FIX_GOOD
	default:
		return 255
	}
}
