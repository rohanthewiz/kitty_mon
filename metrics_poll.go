package main

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func pollTemp() {
	var reading Reading
	for {
		var wait time.Duration
		if opts.Env == "dev" {
			wait = 8 * time.Second
		} else {
			wait = 2 * time.Minute
		}
		time.Sleep(wait)

		// Temperature
		if opts.Bogus {
			reading = bogusReading()
		} else {
			reading = Reading{
				Guid:                 random_sha1(),
				SourceGuid:           whoAmI(),
				IPs:                  IPs(true),
				Temp:                 catTemp(),
				MeasurementTimestamp: time.Now(),
				Sent:                 0,
			}
		}
		db.Save(&reading)
		// Cleanup
		delOlderThanNDays(7) // TODO add config for this
	}
}

func catTemp() int {
	cmdArgs := []string{"/sys/class/thermal/thermal_zone0/temp"}
	byteTemps, err := exec.Command("cat", cmdArgs...).Output()
	if err != nil {
		lpl("Error acquiring temperature.")
		return -255
	}
	strTemp := strings.Trim(string(byteTemps), " \n\t") // clean up whitespace
	var temp int
	temp, err = strconv.Atoi(strTemp)
	if err != nil {
		lpl("Error converting temperature.")
		return -255
	}
	return temp
}
