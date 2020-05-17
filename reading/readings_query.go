package reading

import (
	"kitty_mon/config"
	"kitty_mon/km_db"
	"time"
)

const ReadingsLimit = 120

func GetRecentReadings() ([]Reading, error) {
	var readings []Reading
	q := km_db.Db.Order("created_at desc").Limit(config.Opts.L).Find(&readings)
	return readings, q.Error
}

func QueryReadings(limit int) ([]Reading, error) {
	var readings []Reading

	if limit == -1 || limit == 0 || limit > 120 {
		limit = ReadingsLimit
	}

	q := km_db.Db.Order("created_at desc").Limit(limit).Find(&readings)
	return readings, q.Error
}

func GetAllReadings() ([]Reading, error) {
	var readings []Reading
	q := km_db.Db.Order("created_at desc").Find(&readings)
	return readings, q.Error
}

func DelOlderThanNDays(nDays int) {
	threshold := time.Now().Add(-time.Duration(24*nDays) * time.Hour)
	km_db.Db.Where("created_at < ?", threshold).Delete(Reading{})
}

func DeleteGt2weeks() {
	now := time.Now()
	lastWeek := now.Add(-time.Duration(24*7*2) * time.Hour)
	km_db.Db.Where("created_at < ?", lastWeek).Delete(Reading{})

	//if err := db.Where("name = ?", "jinzhu").First(&user).Error; err != nil {
	//// error handling...
	//}
}
