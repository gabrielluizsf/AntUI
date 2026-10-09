package css

// backMiddle is the default gradient centre: the middle of each axis, with no
// offset.
func backMiddle() BackPos {
	return BackPos{
		{Edge: BackPosMiddle, Off: Zero()},
		{Edge: BackPosMiddle, Off: Zero()},
	}
}
