package main

import (
	"github.com/rohanthewiz/serr"
	"kitty_mon/auth"
	"kitty_mon/config"
	"kitty_mon/km_db"
	"kitty_mon/kmclient"
	"kitty_mon/kmserver"
	"kitty_mon/loaders"
	"kitty_mon/node"
	"kitty_mon/reading"
	"kitty_mon/snapshots"
	"kitty_mon/unloaders"
	"kitty_mon/util"
	"kitty_mon/web"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var err error

func main() {
	loaders.ConfigLoader()
	loaders.DBLoader()

	// ----- UTILITY OPTIONS -----

	// Client - Return our db signature
	if config.Opts.WhoAmI {
		util.Fpl(auth.WhoAmI())
		return
	}

	// Server - Generate an auth token for a client
	// The format of the generated token is: server_id-auth_token_for_the_client
	if config.Opts.GetNodeToken != "" {
		pt, err := node.GetNodeToken(config.Opts.GetNodeToken)
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

	snapshotsStopChan := make(chan bool)
	snapshotsDoneChan := make(chan bool)

	if config.Opts.SynchClient != "" { // Become Client

		go reading.PollTemp() // save temp, whether real or bogus to local db

		go snapshots.RunSnapshotLoop(snapshotsStopChan, snapshotsDoneChan)

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

		close(snapshotsStopChan)
		<-snapshotsDoneChan // give snapshots a chance to wrap up

	} else { // Become server
		// Testing out sending a text
		// err := sms.NexmoSend("KittyMon web client starting " + fmt.Sprintf("%s", time.Now()))

		go web.Webserver(config.Opts.Port)
		kmserver.StartSynchServer()
	}
}
