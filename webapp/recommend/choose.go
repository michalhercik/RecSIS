package recommend

import "math/rand/v2"

func chooseRandom(items []string, n int) []string {
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	if len(items) <= n || n <= 0 {
		return items
	}
	return items[:n]
}
