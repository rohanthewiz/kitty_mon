package main

import (
	"encoding/gob"
	"net"
	"strconv"
)

func synch_client(host string, server_secret string) {
	conn, err := net.Dial("tcp", host+":"+opts.SynchPort)
	if err != nil {
		lpl("Error connecting to server ", err)
		return
	}
	defer func() {
		conn.Close()
		if r := recover(); r != nil {
			lpl("Recovered in synch_client", r)
		}
	}()
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
		guid := msg.Param // retrieve the server's guid
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
		// Get the server's node info from our DB
		node, err := getNodeByGuid(guid)
		if err != nil {
			fpl("Error retrieving node object")
			return
		}
		msg.Param2 = "" // clear for next msg

		// Auth
		msg.Type = "AuthMe"
		msg.Param = node.Token // This is our token for communication with this node (server). It is set by one of two access granting mechanisms

		if opts.NodeName != "" {
			node.Name = opts.NodeName // The server will know this node as this name
			db.Save(&node)            // Save it locally
			msg.Param2 = opts.NodeName
		}
		sendMsg(enc, msg)
		rcxMsg(dec, &msg)
		if msg.Param != "Authorized" {
			fpl("The server declined the authorization request")
			return
		}

		// The Client will send one or more messages to the server
		readings := retrieveUnsentReadings()
		pf("%d unsent readings found\n", len(readings))
		if len(readings) > 0 {
			sendMsg(enc, Message{Type: "NumberOfReadings",
				Param: strconv.Itoa(len(readings))})
			rcxMsg(dec, &msg)
			if msg.Type == "SendReadings" {
				msg.Type = "Reading"
				msg.Param = ""

				for _, reading := range readings {
					reading.SourceGuid = whoAmI()
					msg.Reading = reading
					sendMsg(enc, msg)
					// Let's go ahead and delete here
					//reading.Sent = 1
					db.Delete(&reading) //db.Save(&reading)
				}
			}
		}

	} else {
		fpl("Node does not respond to request for database id")
		fpl("Make sure both server and client databases have been properly setup(migrated) with the -setup_db option")
		fpl("or make sure kitty_mon version is >= 0.9")
		return
	}

	lpl("Synch Operation complete")
}

func retrieveUnsentReadings() []Reading {
	var readings []Reading
	if opts.Bogus == false {
		db.Where("sent = ?", 0).Order("created_at desc").Limit(opts.L).Find(&readings)
	} else {
		// Send some bogus readings for development
		for i := 0; i < 3; i++ {
			readings = append(readings, bogusReading())
		}
	}
	return readings
}
