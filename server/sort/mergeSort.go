package sort

func merge(a []string, b []string) []string {
    final := []string{}
    i := 0
    j := 0
    for i < len(a) && j < len(b) {
        if a[i] < b[j] {
            final = append(final, a[i])
            i++
        } else {
            final = append(final, b[j])
            j++
        }
    }
    for ; i < len(a); i++ {
        final = append(final, a[i])
    }
    for ; j < len(b); j++ {
        final = append(final, b[j])
    }
    return final
}

func Sort(items []string) []string {
    if len(items) < 2 {
        return items
    }
    first := Sort(items[:len(items)/2])
    second := Sort(items[len(items)/2:])
	return merge(first, second)
}