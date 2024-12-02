package util

import (
	"math/rand"
	"time"
)

func ShuffleValues[valueType any](values []valueType) []valueType {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	return values
}
