package timefmt

import "time"

const Layout = "2006-01-02 15:04:05.000"

var localLocation = loadLocalLocation()

func loadLocalLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("CST", 8*60*60)
	}
	return location
}

// Format formats time with millisecond precision using a shared project layout.
func Format(t time.Time) string {
	return t.In(localLocation).Format(Layout)
}
