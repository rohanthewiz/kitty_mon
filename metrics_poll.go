package main
import (
	"time"
	"strconv"
	"os/exec"
	"log"
)

func poll_temp() {
	var reading Reading
	for {
		var wait time.Duration
		if opts_str["env"] == "dev" {
			wait = 8 * time.Second
		} else {
			wait = time.Minute
		}
		time.Sleep(wait)

		//Temperature
		if opts_intf["bogus"].(bool) {
			reading = bogusReading()
		} else {
			reading = Reading{
				Guid: random_sha1(),
				SourceGuid: whoAmI(),
				Temp: cat_temp(),
				Sent: 0,
			}
		}
		db.Save(&reading)
	}
}

func cat_temp() int {
	str_temp, err := exec.Command(`cat /sys/class/thermal/thermal_zone0/temp`).Output()
	if err != nil {
		log.Println("Error acquiring temperature.")
		return -255
	}
	var temp int
	temp, err = strconv.Atoi(string(str_temp))
	if err != nil {
		log.Println("Error acquiring temperature.")
		return -255
	}
	return temp
}