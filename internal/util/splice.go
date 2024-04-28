package util

func SpliceContains(splice []any, sub any) bool {
	for _, v := range splice {
		if v == sub {
			return true
		}
	}
	return false
}
