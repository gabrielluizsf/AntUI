package css

// okLabToSRGB maps OKLab into linear sRGB, the closed-form matrix from the
// CSS Color 4 spec.
func okLabToSRGB(l, a, c float64) (r, g, b float64) {
	l_, m_, s_ := l+0.3963377774*a+0.2158037573*c,
		l-0.1055613458*a-0.0638541728*c,
		l-0.0894841775*a-1.2914855480*c
	l3, m3, s3 := l_*l_*l_, m_*m_*m_, s_*s_*s_
	return 4.0767416621*l3 - 3.3077115913*m3 + 0.2309699292*s3,
		-1.2684380046*l3 + 2.6097574011*m3 - 0.3413193965*s3,
		-0.0041960863*l3 - 0.7034186147*m3 + 1.7076147010*s3
}
