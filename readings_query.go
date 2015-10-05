package main
import (
	"errors"
	"time"
)

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
	q := db.Order("created_at desc").Limit(opts.L).Find(&readings)
	return readings, q.Error
}

func getAllReadings() ([]Reading, error) {
	var readings []Reading
	q := db.Order("created_at desc").Find(&readings)
	return readings, q.Error
}

func delete_gt_2weeks() {
	now := time.Now()
	last_week := now.Add(-time.Duration(24*7*2)*time.Hour)
	db.Where("created_at < ?", last_week).Delete(Reading{})

	//if err := db.Where("name = ?", "jinzhu").First(&user).Error; err != nil {
	//// error handling...
	//}
}

