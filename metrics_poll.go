package main
import (
	"time"
	"strconv"
	"strings"
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
	cmdArgs := []string{"/sys/class/thermal/thermal_zone0/temp"}
	byte_temp, err := exec.Command("cat", cmdArgs...).Output()
	if err != nil {
		log.Println("Error acquiring temperature.")
		return -255
	}
        str_temp := strings.Trim(string(byte_temp), " \n\t") // clean up whitespace
	lpl("\"",str_temp,"\"")
	var temp int
	temp, err = strconv.Atoi(str_temp)
        lpl("The converted temperature is:", temp)
	if err != nil {
		log.Println("Error converting temperature.")
		return -255
	}
	return temp
}
