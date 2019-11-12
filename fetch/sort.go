package fetch

type SortByFromTime []Train

func (s SortByFromTime) Len() int {
	return len(s)
}
func (s SortByFromTime) Less(i, j int) bool {
	return s[i].FromTime.Before(s[j].FromTime)
}
func (s SortByFromTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
