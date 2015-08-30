package main
import(
	"log"
	"net"
	"encoding/gob"
	"time"
	"strconv"
)

// A Node represents a single network entity.
// IsLocal determines if the entry is for self
// Token is the local db signature if IsLocal
// or the mutual auth key for access between local and the specified node
type Node struct {
	Id			int64
	Guid		string `sql: "size:255"`
	Name		string `sql: "size:255"`
	Token		string `sql: "size:255"`
	IP string `sql: "size:255"`
	Status string `sql: "size:255"`
	IsLocal int  // bool 0 - false, 1 - true
	CreatedAt 	time.Time
	UpdatedAt	time.Time
}

const SYNCH_PORT  string = "8090"

func synch_client(host string, server_secret string) {
	conn, err := net.Dial("tcp", host + ":" + SYNCH_PORT)
	if err != nil {log.Fatal("Error connecting to server ", err)}
	defer conn.Close()
	msg := Message{} // init to empty struct
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	defer sendMsg(enc, Message{Type: "Hangup"})

	// Send handshake - Client initiates
	sendMsg(enc, Message{
		Type: "WhoAreYou", Param: whoAmI(), Param2: server_secret,
	})
	rcxMsg(dec, &msg) // Decode the response
	if msg.Type == "WhoIAm" {
		guid := msg.Param  // retrieve the server's guid
		pl("The server's guid is", short_sha(guid))
		if len(guid) != 40 {
			fpl("The server's id is invalid. Run the server once with the -setup_db option")
			return
		}
		// Is there an auth token for us?
		if len(msg.Param2) == 40 {
			setNodeToken(guid, msg.Param2) // make sure to save new auth
									// i.e. Given a node (server) with id guid, our auth token on that server is msg.Param2
		}
		// Obtain the node object which represents the server
		node, err := getNodeByGuid(guid)
		if err != nil { fpl("Error retrieving node object"); return }
		msg.Param2 = ""  // hide the evidence

		// Auth
		msg.Type = "AuthMe"
		msg.Param = node.Token // This is our auth token for this node (server). It is set by one of two access granting mechanisms
		sendMsg(enc, msg)
		rcxMsg(dec, &msg)
		if msg.Param != "Authorized" {
			fpl("The server declined the authorization request")
			return
		}

		// The Client will send one or more messages to the server

    // Get Local Changes
		readings := retrieveUnsentReadings()
		pf("%d unsent readings found\n", len(readings))
		// Push local changes to server
		if len(readings) > 0 {
			sendMsg(enc, Message{Type: "NumberOfReadings",
				Param: strconv.Itoa(len(readings))})
			rcxMsg(dec, &msg)
			if msg.Type == "SendReadings" {
				msg.Type = "Reading"
				msg.Param = ""

				for _, reading := range (readings) {
						reading.SourceGuid = whoAmI()
					msg.Reading = reading
					sendMsg(enc, msg)
					// Mark as sent
					reading.Sent = 1
					db.Save(&reading)
				}
			}
		}

	} else {
			fpl("Node does not respond to request for database id")
			fpl("Make sure both server and client databases have been properly setup(migrated) with the -setup_db option")
			fpl("or make sure node version is >= 0.9")
			return
    }

	defer fpl("Synch Operation complete")
}

func retrieveUnsentReadings() []Reading {
	var readings []Reading
	if opts_str["env"] == "prod" {
		db.Where("sent = ?", 0).Find(&readings)
	} else {
		// Sent some bogus readings for development
		for i := 0; i < 3; i++ {
			readings = append(readings, bogusReading())
		}
	}
	return readings
}
