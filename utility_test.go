package main

import (
	"testing"
)

var testSlice = []string{"one", "two", "three", "four"}

func TestShuffleStringSlice_Length(t *testing.T) {
	shuffledSlice := make([]string, len(testSlice))
	copy(shuffledSlice, testSlice)
	ShuffleStringSlice(shuffledSlice)

	if len(testSlice) != len(shuffledSlice) {
		t.Errorf("The length of slices are not equal. Got %d, expected %d", len(shuffledSlice), len(testSlice))
	}
}

func TestShuffleStringSlice_Items(t *testing.T) {
	shuffledSlice := make([]string, len(testSlice))
	copy(shuffledSlice, testSlice)
	ShuffleStringSlice(shuffledSlice)

	for _, item := range testSlice {
		if IsStringInSlice(item, shuffledSlice) == false {
			t.Errorf("Item \"%s\" not found in shuffledSlice: %+v", item, shuffledSlice)
		}
	}
}

func TestCrop(t *testing.T) {
	testSentence := "This is a test sentence for the unit test."
	textMock := []struct {
		Content     string
		Chars       int
		AfterString string
		Crop2Space  bool
		Result      string
	}{
		{testSentence, 0, "", false, testSentence},
		{testSentence, 99, "", false, testSentence},
		{testSentence, 13, "", false, "This is a te"},
		{testSentence, 13, "...", false, "This is a te..."},
		{testSentence, 13, "", true, "This is a"},
		{testSentence, 13, "...", true, "This is a..."},
		{testSentence, -99, "", false, testSentence},
		{testSentence, -13, "", false, "he unit test."},
		{testSentence, -13, "...", false, "...he unit test."},
		{testSentence, -13, "", true, "unit test."},
		{testSentence, -13, "...", true, "...unit test."},
	}

	for _, mock := range textMock {
		res := Crop(mock.Content, mock.Chars, mock.AfterString, mock.Crop2Space)
		if res != mock.Result {
			t.Errorf("Crop result is \"%s\", but expected \"%s\".", res, mock.Result)
		}
	}
}

func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
