package location

func CalculateState(pd PositionData) uint8 {
	haveSkyReport := pd.SKYReport != nil
	have3DFix := (*pd.TPVReport).Mode == 3
	sats := 0
	if haveSkyReport {
		sats = len((*pd.SKYReport).Satellites)
	}

	switch {
	case !haveSkyReport:
		return 1
	case haveSkyReport && !have3DFix:
		return 2
	case haveSkyReport && have3DFix && sats <= 6:
		return 3
	case haveSkyReport && have3DFix && sats > 6:
		return 4
	default:
		return 255
	}
}
