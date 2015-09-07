package main

import (
	"fmt"
	"os"
	//"strings"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

const app_name = "Kitty Monitor"
const version string = "0.1.5"

// Get Commandline Options and Flags
var opts_str, opts_intf = getOpts()

// Init db
var db gorm.DB; var err error

func migrate() {
	// Create or update the table structure as needed
	pl("Migrating the DB...")
	db.AutoMigrate(&Node{})
	db.AutoMigrate(&Reading{})
	//According to GORM: Feel free to change your struct, AutoMigrate will keep your database up-to-date.
	// Fyi, AutoMigrate will only *add new columns*, it won't update column's type or delete unused columns for safety
	// If the table is not existing, AutoMigrate will create the table automatically.

	db.Model(&Reading{}).AddUniqueIndex("idx_reading_guid", "guid")
	db.Model(&Reading{}).AddIndex("idx_reading_source_guid", "source_guid")
	db.Model(&Reading{}).AddIndex("idx_reading_created_at", "created_at")

	db.Model(&Node{}).AddUniqueIndex("idx_node_guid", "guid")
	// This would disallow blanks //db.Model(&Node{}).AddUniqueIndex("idx_node_name", "name")

	pl("Migration complete")
	ensureDBSig() // Initialize local with a SHA1 signature if it doesn't already have one
}

func ensureDBSig() {
	if LocalNode.Id > 0 && len(LocalNode.Token) == 40 { return /* all is good */ }
	var node Node
	if db.Where("is_local = 1").First(&node); node.Id < 1 { // create the signature
		db.Create(&Node{Guid: random_sha1(), Token: random_sha1(), IsLocal: 1})
		if db.Where("is_local = 1").First(&node); node.Id > 0 && len(node.Token) == 40 { // was it saved?
			pl("Local signature created")
		}
	} else {
			lf("Something is whacked with the local db signature. Manually check the DB.")
	}
}

func main() {
	var err error
	db, err = gorm.Open("sqlite3", opts_str["db_path"])
	if err != nil {
		fpl("There was an error connecting to the DB")
		fpl("DBPath: " + opts_str["db_path"])
		os.Exit(2)
	}

	//Do we need to migrate?
	if ! db.HasTable(&Node{}) || ! db.HasTable(&Reading{}) { migrate() }

	if opts_intf["v"].(bool) {
		fpl(app_name, version)
		return
	}

	db.LogMode(opts_intf["debug"].(bool)) // Set debug mode for Gorm db

	if opts_str["admin"] == "delete_tables" {
		fmt.Println("Are you sure you want to delete all data? (N/y)")
		var input string
		fmt.Scanln(&input) // Get keyboard input
		pd("input", input)
		if input == "y" || input == "Y" {
			db.DropTableIfExists(&Reading{})
			db.DropTableIfExists(&Node{})
			pl("Readings tables deleted")
		}
		return
	}

	// Client - Return our db signature
	if opts_intf["whoami"].(bool) {
		fpl(whoAmI())
		return
	}

	// Server - Generate an auth token for a client
	// The format of the generated token is: server_id-auth_token_for_the_client
	if opts_str["get_node_token"] != "" {
		pt, err := getNodeToken(opts_str["get_node_token"])
		if err != nil {fpl("Error retrieving token"); return}
		fpf("Node token is: %s-%s\nYou will now need to run the client with \n'go_notes -save_node_token the_token'\n",
			whoAmI(), pt)
		return
	}

	// Client - Save a token generated for us by a server
	if opts_str["save_node_token"] != "" {
		saveNodeToken(opts_str["save_node_token"])
		return
	}

	// Server - Return the server's secret token
	// This is a master key and will allow any client to auth
	// We probably want to use the methods above instead
	if opts_intf["get_server_secret"].(bool) {
		fpl(get_server_secret())
		return
	}
	if opts_intf["setup_db"].(bool) { // Migrate the DB
		migrate()
		return
	}

	// CORE PROCESSING

	if opts_str["synch_client"] != "" {
		go poll_temp() // save temp, whether real or bogus to the db

		lpl("I will periodically send data to server...")
		for {
			var wait time.Duration
			if opts_str["env"] == "dev" {
				wait = 12 * time.Second
			} else {
				wait = 2 * time.Minute
			}
			time.Sleep(wait)

			synch_client(opts_str["synch_client"], opts_str["server_secret"])
		}

	} else {  // Become server
		go webserver(opts_str["port"])
		synch_server()
		// opts_intf["synch_server"].(bool)
	}


// CODE SCRAP

//	} else if opts_str["q"] != "" || opts_intf["qi"].(int64) != 0 ||
//				opts_str["qg"] != "" || opts_str["qt"] != "" ||
//				opts_str["qb"] != "" || opts_str["qd"] != "" {
//		// QUERY
//		readings := queryReadings()
//
//		// List Readings found
//		fpl("")  // for UI sake
//		listReadings(readings, true)
//
//		// Options that can go with Query
//		// export
//		if opts_str["exp"] != "" {
//			arr := strings.Split(opts_str["exp"], ".")
//			arr_item_last := len(arr) -1
//			if arr[arr_item_last] == "csv" {
//				exportCsv(readings, opts_str["exp"])
//			}
//			if arr[arr_item_last] == "gob" {
//				exportGob(readings, opts_str["exp"])
//			}
//		} else if opts_intf["upd"].(bool) { // update
//			updateReadings(readings)
//
//			// See if we want to delete
//		} else if opts_intf["del"].(bool) {
//			deleteReadings(readings)
//		}
//		// Other options
//	} else if opts_str["imp"] != "" { // import
//			arr := strings.Split(opts_str["imp"], ".")
//			arr_item_last := len(arr) -1
//			if arr[arr_item_last] == "csv" {
//				importCsv(opts_str["imp"])
//			}
//			if arr[arr_item_last] == "gob" {
//				importGob(opts_str["imp"])
//			}
		// Create
//	} else if opts_str["t"] != "" { // No query options, we must be trying to CREATE
//		createReading(opts_str["t"], opts_str["d"], opts_str["b"], opts_str["g"])
}
