package main

import (
	"math/rand"
)

// ShuffleStringSlice will randomize a string slice.
// I know that is a really bad shuffle logic (i won`t call this an algorithm,
// why? because i wrote and understand it :D)
// But this is YOUR chance to contribute to an open source project.
// Replace this by a cool one!
func ShuffleStringSlice(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
