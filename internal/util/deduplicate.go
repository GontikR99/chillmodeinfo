package util

func Deduplicate(values []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, v := range values {
		if _, ok := seen[v]; !ok {
			result = append(result, v)
			seen[v]=struct{}{}
		}
	}
	return result
}