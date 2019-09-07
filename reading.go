package main

import (
	"fmt"
	"kitty_mon/sms"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

var criticalTemp int

func init() {
	var err error
	criticalTemp, err = strconv.Atoi(os.Getenv("CRITICAL_TEMP"))
	if err != nil || criticalTemp == 0 {
		panic("CRITICAL_TEMP env var is not a valid positive integer - e.g. 70")
	}
	criticalTemp *= 1000
}

// This represents the payload sent to the server
// This is the equivalent of a NoteChange in GoNotes
type Reading struct {
	Id                   int64
	Guid                 string    `sql: "size:40"` // random id for each message
	SourceGuid           string    `sql: "size:40"` // We will tag this with the client's db sig when reading sent
	IPs                  string    `sql: "size:254"`
	Sent                 int       `json:"-"` // has the reading been sent // bool 0 - false, 1 - true
	Temp                 int       // temperature
	MeasurementTimestamp time.Time // True CreatedAt for the reading,
	// since GORM also updates CreatedAt when saved on the server
	CreatedAt time.Time // GORM automatically updates this field on save
}

type ReadingEnriched struct {
	Reading
	Name   string
	Status string
}

type byCreatedAt []Reading

func (ncs byCreatedAt) Len() int {
	return len(ncs)
}
func (ncs byCreatedAt) Less(i, j int) bool {
	return ncs[i].CreatedAt.Before(ncs[j].CreatedAt)
}
func (ncs byCreatedAt) Swap(i, j int) {
	ncs[i], ncs[j] = ncs[j], ncs[i]
}

func processReadings(readings *[]Reading) {
	pl("Processing received changes...")
	sort.Sort(byCreatedAt(*readings)) // we will apply in created order
	var highReading Reading

	for _, reading := range *readings {

		// Skip if we already have this Reading locally
		local_reading := Reading{} // make sure local_change is inited here
		// otherwise GORM uses its id in the query - weird!
		db.Where("guid = ?", reading.Guid).First(&local_reading)
		if local_reading.Id > 1 {
			pf("We already have Reading: %s -- skipping\n", short_sha(local_reading.Guid))
			continue
		}

		// Track high temps
		if reading.Temp > highReading.Temp {
			highReading = reading
		}

		// Save the reading
		pl("__________SAVING READING_____________")
		reading.save()
	}

	// See if overtemp condition
	if highReading.Temp > criticalTemp {
		rdgsWNames := readingsPopulateNodeName([]Reading{highReading})
		name := "unknown"
		if len(rdgsWNames) > 0 {
			name = rdgsWNames[0].Name
		}
		msg := fmt.Sprintf("High temperature %d received from %s measured at %s",
			highReading.Temp, name, highReading.MeasurementTimestamp)
		fmt.Println(msg)
		err := sms.NexmoSend(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (r Reading) save() bool {
	r.Id = 0 // Make sure the reading has a zero id for db creation
	// A non-zero Id will not be created
	db.Create(&r)         // will auto create contained objects too and it's smart - 'nil' children will not be created :-)
	if !db.NewRecord(r) { // was it saved?
		pl("Reading saved:", short_sha(r.Guid))
		return true
	}
	fpl("Failed to save reading:", short_sha(r.Guid))
	return false
}

func doDelete(reading Reading) {
	if reading == (Reading{}) {
		pf("Internal error: cannot delete non-existent reading")
		return
	}
	db.Delete(&reading)
}

func (r *Reading) Print() {
	pf("%+v\n", r)
}

func bogusReading() Reading {
	return Reading{
		Guid:       random_sha1(),
		SourceGuid: whoAmI(),
		Temp:       rand.Intn(100),
	}
}

func find_reading_by_id(id int64) Reading {
	var reading Reading
	db.First(&reading, id)
	return reading
}
