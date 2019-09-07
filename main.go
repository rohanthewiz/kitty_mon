package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"kitty_mon/sms"
	"os"
	"time"
)

const app_name = "Kitty Monitor"
const version string = "0.1.9"

var opts *Opts //Cmdline options and flags

// Init db
var db gorm.DB
var err error

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
	if LocalNode.Id > 0 && len(LocalNode.Token) == 40 {
		return /* all is good */
	}
	var node Node
	if db.Where("is_local = 1").First(&node); node.Id < 1 { // create the signature
		db.Create(&Node{Guid: random_sha1(), Token: random_sha1(), IsLocal: 1})
		if db.Where("is_local = 1").First(&node); node.Id > 0 && len(node.Token) == 40 { // was it saved?
			pl("Local signature created")
		}
	} else {
		pl("Local db signature already exists")
	}
}

func main() {
	opts = NewOpts()
	db, err = gorm.Open("sqlite3", opts.DbPath)
	if err != nil {
		lf("There was an error connecting to the DB.\nDBPath: " + opts.DbPath)
		os.Exit(2)
	}

	//Do we need to migrate?
	if !db.HasTable(&Node{}) || !db.HasTable(&Reading{}) {
		migrate()
	}

	if opts.V {
		fpl(app_name, version)
		fpl(catTemp())
		fpl("Local IPs:", IPs(false))
		return
	}

	db.LogMode(opts.Debug) // Set debug mode for Gorm db

	if opts.Admin == "delete_tables" {
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
	if opts.WhoAmI {
		fpl(whoAmI())
		return
	}

	// Server - Generate an auth token for a client
	// The format of the generated token is: server_id-auth_token_for_the_client
	if opts.GetNodeToken != "" {
		pt, err := getNodeToken(opts.GetNodeToken)
		if err != nil {
			fpl("Error retrieving token")
			return
		}
		fpf("Node token is: %s-%s\nYou will now need to run the client with \n'kitty_mon -save_node_token the_token'\n",
			whoAmI(), pt)
		return
	}

	// Client - Save a token generated for us by a server
	if opts.SaveNodeToken != "" {
		saveNodeToken(opts.SaveNodeToken)
		return
	}

	// Server - Return the server's secret token
	// This is a master key and will allow any client to auth
	// We probably want to use the methods above instead
	if opts.GetServerSecret {
		fpl(get_server_secret())
		return
	}
	if opts.SetupDb { // Migrate the DB
		migrate()
		return
	}

	// CORE PROCESSING

	if opts.SynchClient != "" {
		go pollTemp() // save temp, whether real or bogus to the db

		lpl("I will periodically send data to server...")
		for {
			var wait time.Duration
			if opts.Env == "dev" {
				wait = 12 * time.Second
			} else {
				wait = 4 * time.Minute
			}
			time.Sleep(wait)

			synch_client(opts.SynchClient, opts.ServerSecret)
		}

	} else { // Become server
		// Testing out sending a text
		err := sms.NexmoSend("KittyMon web client starting " + fmt.Sprintf("%s", time.Now()))
		if err != nil {
			fmt.Println(err)
		}

		go webserver(opts.Port)
		synch_server()
	}
}
