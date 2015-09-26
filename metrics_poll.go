package main
import (
	"time"
	"strconv"
	"strings"
	"os/exec"
)

func poll_temp() {
	var reading Reading
	for {
		var wait time.Duration
		if opts.Env == "dev" {
			wait = 8 * time.Second
		} else {
			wait = 2 * time.Minute
		}
		time.Sleep(wait)

		//Temperature
		if opts.Bogus {
			reading = bogusReading()
		} else {
			reading = Reading{
				Guid: random_sha1(),
				SourceGuid: whoAmI(),
				IPs: IPs(true),
				Temp: cat_temp(),
				Sent: 0,
			}
		}
		db.Save(&reading)
	}
}

func cat_temp() int {
	cmdArgs := []string{"/sys/class/thermal/thermal_zone0/temp"}
	byte_temp, err := exec.Command("cat", cmdArgs...).Output()
	if err != nil {
		lpl("Error acquiring temperature.")
		return -255
	}
	str_temp := strings.Trim(string(byte_temp), " \n\t") // clean up whitespace
	//lpl("\"",str_temp,"\"")
	var temp int
	temp, err = strconv.Atoi(str_temp)
	if err != nil {
		lpl("Error converting temperature.")
		return -255
	}
	return temp
}
