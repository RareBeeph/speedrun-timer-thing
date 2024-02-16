package timer

func (t *timer) previousSegment() *Split {
	if t.segment == 0 {
		return t.run.Segments[t.segment]
	}

	// TODO: Actually check that this a valid segment
	return t.run.Segments[t.segment-1]
}
