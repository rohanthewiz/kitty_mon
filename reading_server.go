// GOB server
package main
import(
	//"reflect"
	"os"
	"net"
	"encoding/gob"
	"fmt"
	"strconv"
)

func synch_server() {
	ln, err := net.Listen("tcp", ":" + SYNCH_PORT) // counterpart of net.Dial
	if err != nil {
		fpl("Error setting up server listen on port", SYNCH_PORT)
		return
	}
	fmt.Println("Server listening on port: " + SYNCH_PORT + " - CTRL-C to quit")

	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil { continue }
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	msg := Message{}
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	var node_id string
	var node Node
	var err error
	var authorized bool = false // true to disable auth

	for {
		msg = Message{}
		rcxMsg(decoder, &msg)

		switch msg.Type {
		case "Hangup":
			printHangupMsg(conn); return

		case "Quit":
			pl("Quit message received. Exiting..."); os.Exit(1)

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
			//     (a) The client does ./go_notes -whoami  # output the client's guid
			//     (b) At the server a token is generated for the client with
			//         ./go_notes -get_node_token node_guid
			//     (c) The client saves this token with
			//         ./go_notes -save_node_token the_token (detail: the_token here is of the format "server_id-auth_token")
		case "WhoAreYou":
			node_id = msg.Param // receive the client's GUID (db signature)
			pd("Client id is:", short_sha(node_id))
			if msg.Param2 == get_server_secret() { // if client knows the server's token, give them a token of their own on this server
				pt, err := getNodeToken(node_id)
				if err != nil {
					msg.Param2 = ""
				} else {
					pl("Auth token generated:", pt)
					msg.Param2 = pt // include the auth token in next msg
				}
			} else {
				msg.Param2 = "" // no token for client returned
			}
			// Always return the server's id
			msg.Param = whoAmI() // reply with the server's db signature
			msg.Type = "WhoIAm"
			// Retrieve the actual node object which represents the client
			node, err = getNodeByGuid(node_id)
			if err != nil {
				fpl("Error retrieving node object for node:", short_sha(node_id));
				fpf("Could not find a record for %s on the server.\n", node_id)
				return
			}
			sendMsg(encoder, msg)

		case "AuthMe":
			if node.Token == msg.Param {
				authorized = true
				msg.Param = "Authorized"
			} else {
				msg.Param = "Declined"
			}
			sendMsg(encoder, msg)

			// Pickup Data from Client
		case "NumberOfReadings":
			if !authorized { pl(authFailMsg); return }
			numReadings, err := strconv.Atoi(msg.Param)
			if err != nil {
				pl("Could not decode the number of change messages"); return
			}
			if numReadings < 1 { pl("No remote changes."); return }

			readings := make([]Reading, numReadings) // prealloc
			sendMsg(encoder, Message{Type: "SendReadings"}) // Send the actual changes
			// Receive readings, extract the Readings, save into readings
			for i := 0; i < numReadings; i++ {
				msg = Message{}
				rcxMsg(decoder, &msg)
				readings[i] = msg.Reading
			}
			pf("\n%d readings received:\n", numReadings)
			processReadings(&readings)

//
//		case "NewSynchPoint": // New synch point at the end of synching
//			if !authorized { pl(authFailMsg); return }
//			synch_nc := msg.NoteChg
//			synch_nc.Id = 0 // so it will save
//			db.Save(&synch_nc)
//			node.SynchPos = synch_nc.Guid
//			db.Save(&node)

		default: // It is essential to do this here on GOB communication
			pl("Unknown message type received", msg.Type)
			printHangupMsg(conn); return
		}
	}
}

func sendMsg(encoder *gob.Encoder, msg Message) {
	encoder.Encode(msg); printMsg(msg, false)
	//time.Sleep(10)
}

func rcxMsg(decoder *gob.Decoder, msg *Message) {
	//time.Sleep(10)
	decoder.Decode(&msg); printMsg(*msg, true)
}

func printHangupMsg(conn net.Conn) {
	fmt.Printf("Closing connection: %+v\n----------------------------------------------\n", conn)
}

func printMsg(msg Message, rcx bool) {
	pl("\n----------------------------------------------")
	if rcx { print("Received: ")
	} else {
		print("Sent: ")
	}
	pl("Msg Type:", msg.Type, " Msg Param:", short_sha(msg.Param))
	msg.Reading.Print()
}


// CODE_SCRAP // Yes. A compiled language allows us to do this without any runtime penalty
//	fmt.Printf("encoder is a type of: %v\n", reflect.TypeOf(encoder))

//			// Send a Create Change
//			noteGuid := random_sha1() // we use the note guid in two places (a little denormalization)
//			note1Guid := noteGuid
//			msg.NoteChg = NoteChange{
//				Operation: 1,
//				Guid: random_sha1(),
//				NoteGuid: noteGuid,
//				Note: Note{
//					Guid: noteGuid, Title: "Synch Note 1",
//					Description: "Description for Synch Note 1", Body: "Body for Synch Note 1",
//					Tag: "tag_synch_1", CreatedAt: time.Now()},
//				NoteFragment: NoteFragment{},
//			}
//			sendMsg(enc, msg)
//
//			// Send another Create Change
//			noteGuid = random_sha1()
//			msg.NoteChg = NoteChange{
//				Operation: 1,
//				Guid: random_sha1(),
//				NoteGuid: noteGuid,
//				Note: Note{
//					Guid: noteGuid, Title: "Synch Note 2",
//					Description: "Description for Synch Note 2", Body: "Body for Synch Note 2",
//					Tag: "tag_synch_2", CreatedAt: time.Now().Add(time.Second)},
//				NoteFragment: NoteFragment{},
//			}
//			second_note_guid := msg.NoteChg.NoteGuid // save for use in update op
//			sendMsg(enc, msg)
//
//			// Send an update operation
//			msg.NoteChg = NoteChange{
//				Operation: 2,
//				Guid: random_sha1(),
//				NoteGuid: second_note_guid,
//				Note: Note{},
//				NoteFragment: NoteFragment{
//						Bitmask: 0xC, Title: "Synch Note 2 - Updated",
//						Description: "Updated!"},
//			}
//			sendMsg(enc, msg)
//
//			// Send a Delete Change
//			msg.NoteChg = NoteChange{
//				Operation: 3,
//				Guid: random_sha1(),
//				NoteGuid: note1Guid,
//				Note: Note{},
//				NoteFragment: NoteFragment{},
//			}
//			sendMsg(enc, msg)
