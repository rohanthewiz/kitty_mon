package main

import (
	"fmt"
	"kitty_mon/auth"
	"kitty_mon/config"
	"kitty_mon/km_db"
	"kitty_mon/kmclient"
	"kitty_mon/kmserver"
	"kitty_mon/kredis"
	"kitty_mon/loaders"
	"kitty_mon/node"
	"kitty_mon/reading"
	"kitty_mon/unloaders"
	"kitty_mon/util"
	"kitty_mon/web"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rohanthewiz/roredis"
	"github.com/rohanthewiz/serr"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	loaders.ConfigLoader()
	loaders.DBLoader()

	// ----- UTILITY OPTIONS -----

	// Return our db signature
	if config.Opts.WhoAmI {
		util.Fpl(auth.WhoAmI())
		return
	}

	// TODO let's modify this approach as it involves stopping the server
	// TODO Rather, let's use some kind of admin login approach
	// TODO - First let's put this in a function
	// Server - Generate an auth token for a client
	// The format of the generated token is: server_id-auth_token_for_the_client
	if config.Opts.GetTokenForNode != "" {
		pt, err := node.GetNodeToken(config.Opts.GetTokenForNode)
		if err != nil {
			util.Fpl("Error retrieving token")
			return
		}
		util.Fpf("Node token is: %s-%s\nYou will now need to run the client with \n'kitty_mon -save_node_token the_token'\n",
			auth.WhoAmI(), pt)
		return
	}

	// Client - Save a token generated for us by a server
	if config.Opts.SaveNodeToken != "" {
		node.SaveNodeToken(config.Opts.SaveNodeToken)
		return
	}

	// Server - Return the server's secret token
	// This is a master key and will allow any client to auth
	// We probably want to use the methods above instead
	if config.Opts.GetServerSecret {
		util.Fpl(auth.GetServerSecret())
		return
	}

	if config.Opts.SetupDb { // Migrate the DB
		km_db.Migrate(&reading.Reading{}, &node.Node{})
		auth.EnsureDBSig() // Initialize local with a SHA1 signature if it doesn't already have one
		return
	}

	// CORE PROCESSING

	// Disable image snapshots for now
	// snapshotsStopChan := make(chan bool)
	// snapshotsDoneChan := make(chan bool)

	if config.Opts.SynchClient != "" { // Become Client

		go reading.PollTemp() // save temp, whether real or bogus to local db

		// go snapshots.RunSnapshotLoop(snapshotsStopChan, snapshotsDoneChan)

		wait := config.ReadingsProdPollRate
		if config.Opts.Env == "dev" {
			wait = config.ReadingsDevPollRate
		}

		util.Lpl("I will periodically send data to server...")

		networkErrCount := 0

		for {
			if networkErrCount > 3 {
				_ = os.Setenv("KM_SHUTDOWN", "true") // let everyone know we are shutting down
				_ = unloaders.Reboot()
				break
			}

			// The app behavior can be dynamically changed via env vars
			if strings.ToLower(os.Getenv("KM_SHUTDOWN")) == "true" {
				break
			}

			if strRate := strings.ToLower(os.Getenv("KM_READINGS_POLLRATE")); strRate != "" {
				rate, err := strconv.Atoi(strRate)
				if err != nil {
					util.Lpl("Error converting readings pollrate from env var (KM_READINGS_POLLRATE) " + err.Error())
				} else {
					wait = time.Duration(rate) * time.Second
				}
			}

			time.Sleep(wait)

			err := kmclient.SynchAsClient(config.Opts.SynchClient, config.Opts.ServerSecret)
			if ser, ok := err.(serr.SErr); ok {
				mp := ser.FieldsMap()
				if str, ok := mp["msg"]; ok && strings.Contains(str, kmclient.NetworkConnErrorMsg) {
					networkErrCount++
				}
			} else {
				networkErrCount--
				if networkErrCount < 0 {
					networkErrCount = 0
				}
			}
		}

		// Disable image snapshots for now
		// close(snapshotsStopChan)
		// <-snapshotsDoneChan // give snapshots a chance to wrap up

	} else { // Become server
		// Server will store to Redis
		host := "localhost"
		if h := os.Getenv("REDIS_HOST"); h != "" {
			host = h
		}

		port := "6379"
		if p := os.Getenv("REDIS_PORT"); p != "" {
			port = p
		}

		db := 0
		if d := os.Getenv("REDIS_DB"); d != "" {
			dtmp, err := strconv.Atoi(d)
			if err == nil {
				db = dtmp
			} else {
				log.Fatal("Could not convert redis db ", d)
			}
		}

		fmt.Println("Host", host, "Port", port, "DB", db)
		kc := kredis.InitClient(host, port, db)

		if resp := roredis.Ping(kc); resp != "PONG" {
			fmt.Printf("resp %#v\n", resp)
			log.Fatal("Unable to ping redis server")
		}

		// Testing roredis
		// err = roredis.Set(kc, "key1", "val1", 15 * time.Second)
		// if err != nil {
		// 	log.Fatal("Redis set failed", err.Error())
		// }
		//
		// val, er := roredis.Get(kc, "key1")
		// if er != nil {
		// 	log.Fatal("Unable to get test val - ", er)
		// }
		// fmt.Println("Test value rcxd:", val)

		// Testing out sending a text
		// err := sms.NexmoSend("KittyMon web client starting " + fmt.Sprintf("%s", time.Now()))

		go web.Webserver(config.Opts.Port)
		kmserver.StartSynchServer()
	}
}
