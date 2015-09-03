package main
import "errors"

func getReading(guid string) (Reading, error) {
	var reading Reading
	db.Where("guid = ?", guid).First(&reading)
	if reading.Id != 0 {
		return reading, nil
	} else {
		return reading, errors.New("Note not found")
	}
}

func getRecentReadings() ([]Reading, error) {
	var readings []Reading
	q := db.Order("created_at desc").Limit(60).Find(&readings)
	return readings, q.Error
}
