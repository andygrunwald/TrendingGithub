package main

import (
	"testing"
)

func TestTrending_GetTimeFrames_Length(t *testing.T) {
	trend := Trend{}
	timeFrames := trend.GetTimeFrames()

	if len(timeFrames) == 0 {
		t.Errorf("Expected more than %d timeframes", len(timeFrames))
	}
}
