package config

// TODO implement persistent storage
func GetTitleBlacklist(channel string) []string {
	return []string{}
}

func GetDescriptionBlacklist(channel string) []string {
	return []string{}
}

func GetRatingThreshold(channel string) float64 {
	return 0
}

func GetViewThreshold(channel string) float64 {
	return 0
}

func AddTitleBlacklist(channel, word string) {

}
func AddDescriptionBlacklist(channel, word string) {

}
func SetRatingThreshold(channel string, threshold float64) {

}
func SetViewThreshold(channel string, threshold uint64) {

}
