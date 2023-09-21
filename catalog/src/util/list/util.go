package list

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

// FindMissing find items missing from second slice
func FindMissing[T comparable](s1 []T, s2 []T) []T {
	var missing []T

	for _, i := range s1 {
		if !Contains[T](s2, i) {
			missing = append(missing, i)
		}
	}

	return missing
}

// Contains checks if the slice contains the item
func Contains[T comparable](haystack []T, needle T) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}
