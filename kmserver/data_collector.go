package kmserver

import (
	"encoding/gob"
	"kitty_mon/auth"
	"kitty_mon/config"
	"kitty_mon/km_db"
	"kitty_mon/message"
	node2 "kitty_mon/node"
	"kitty_mon/reading"
	"kitty_mon/util"
	"net"
	"os"
	"sort"
	"strconv"
)

func StartSynchServer() {
	ln, err := net.Listen("tcp", ":"+config.Opts.SynchPort) // counterpart of net.Dial
	if err != nil {
		util.Fpl("Error setting up server listen on port", config.Opts.SynchPort)
		return
	}
	util.Fpl("GOB Server listening on port: "+config.Opts.SynchPort+" - CTRL-C to quit",
		"\nLocal IPs:", util.IPs(false))

	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			continue
		}
		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	msg := message.Message{}
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	var node_id string
	var node node2.Node
	var err error
	var authorized bool // true to disable auth

	for {
		msg = message.Message{}
		message.RcxMsg(decoder, &msg)

		switch msg.Type {
		case "Hangup":
			message.PrintHangupMsg(conn)
			return

		case "Quit":
			util.Pl("Quit message received. Exiting...")
			os.Exit(1)

			// This is the point of id exchange between the server and client
			// Normal auth process is that the client provides its db signature(node_id)
			// which the server already has a record of (together with its auth_token)
			// The client provides the required auth_token in the following auth. request message
			// The server must know of the client and auth_token before actual authentication.
			// This is achieved one of two ways:
			// (1) The client knows the server's secret token and provides it in the WhoAreYou message
			//     A token is then automatically generated for the client and stored on both ends.
			//     The server stores this as a new node entry.
			// (2) In a manual operation (At the cmd line for now), the client provides its db signature (guid)
			//     to the server in a request for an auth_token
			//     (a) The client does ./kitty_mon -whoami  # output the client's guid
			//     (b) At the server a token is generated for the client with
			//         ./kitty_mon -get_node_token node_guid
			//     (c) The client saves this token with
			//         ./kitty_mon -save_node_token the_token (detail: the_token here is of the format "server_id-auth_token")
		case "WhoAreYou":
			node_id = msg.Param // receive the client's GUID (db signature)
			util.Pd("Client id is:", util.Short_sha(node_id))
			if msg.Param2 == auth.GetServerSecret() { // TODO - Do this differently // if client knows the server's token, give them a token of their own on this server
				pt, err := node2.GetNodeToken(node_id)
				if err != nil {
					msg.Param2 = ""
				} else {
					util.Pl("Auth token generated:", pt)
					msg.Param2 = pt // include the auth token in next msg
				}
			} else {
				msg.Param2 = "" // no token for client returned
			}
			// Always return the server's id
			msg.Param = auth.WhoAmI() // reply with the server's db signature
			msg.Type = "WhoIAm"
			// Retrieve the actual node object which represents the client
			node, err = node2.GetNodeByGuid(node_id)
			if err != nil {
				util.Fpl("Error retrieving node object for node:", util.Short_sha(node_id))
				util.Fpf("Could not find a record for %s on the server.\n", node_id)
				return
			}
			message.SendMsg(encoder, msg)

		case "AuthMe":
			if node.Token == msg.Param {
				authorized = true
				msg.Param = "Authorized"
				util.Fpf("Message on auth: %v", msg)
				if msg.Param2 != "" { // While we are here, client wants to set/update it's name
					node.Name = msg.Param2
					km_db.Db.Save(&node)
				}
			} else {
				msg.Param = "Declined"
			}
			message.SendMsg(encoder, msg)

			// Pickup Data from Client
		case "NumberOfReadings":
			if !authorized {
				util.Pl(auth.AuthFailMsg)
				return
			}
			numReadings, err := strconv.Atoi(msg.Param)
			if err != nil {
				util.Pl("Could not decode the number of readings")
				return
			}
			if numReadings < 1 {
				util.Pl("No remote changes.")
				return
			}

			readings := make([]reading.Reading, numReadings)                // prealloc
			message.SendMsg(encoder, message.Message{Type: "SendReadings"}) // Request to send the actual changes
			// Receive readings, extract the Readings, save into readings
			for i := 0; i < numReadings; i++ {
				msg = message.Message{}
				message.RcxMsg(decoder, &msg)
				readings[i] = msg.Reading
			}
			util.Pf("\n%d readings received:\n", numReadings)
			ProcessReadings(&readings)

		default: // It is essential to do this here on GOB communication
			util.Pl("Unknown message type received", msg.Type)
			message.PrintHangupMsg(conn)
			return
		}
	}
}

func ProcessReadings(readings *[]reading.Reading) {
	util.Pl("Processing received changes...")
	sort.Sort(reading.ByCreatedAt(*readings)) // we will apply in created order
	var highReading reading.Reading

	for _, rdg := range *readings {

		// Skip if we already have this Reading locally
		localReading := reading.Reading{} // make sure local_change is inited here
		// otherwise GORM uses its id in the query - weird!
		km_db.Db.Where("guid = ?", rdg.Guid).First(&localReading)
		if localReading.Id > 1 {
			util.Pf("We already have Reading: %s -- skipping\n", util.Short_sha(localReading.Guid))
			continue
		}

		// Track high temps
		if rdg.Temp > highReading.Temp {
			highReading = rdg
		}

		// Save the reading
		util.Pl("__________SAVING READING_____________")
		rdg.Save()
	}

	// Turn off for now due to use of Odroid UX4 // See if overtemp condition
	// if highReading.Temp > reading.CriticalTemp {
	// 	rdgsWNames := reading.ReadingsPopulateNodeName([]reading.Reading{highReading})
	// 	name := "unknown"
	// 	if len(rdgsWNames) > 0 {
	// 		name = rdgsWNames[0].Name
	// 	}
	// 	msg := fmt.Sprintf("High temperature %d received from %s measured at %s",
	// 		highReading.Temp, name, highReading.MeasurementTimestamp)
	// 	fmt.Println(msg)
	// 	err := sms.NexmoSend(msg)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
}
