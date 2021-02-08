package auth

import (
	"encoding/json"
	"kitty_mon/config"
	"kitty_mon/km_db"
	"kitty_mon/kredis"
	"kitty_mon/node"
	"kitty_mon/util"
	"os"

	"github.com/rohanthewiz/roredis"
)

const serverKey = "km.server"

// Ensure this node has an identity (DB Signature)
// (This is a node entry representing this node (servers act as one node btw) with is_local = 1)
// Servers will share the DBSig -- however server will only use Redis for storage
func EnsureDBSig() { // TODO - let's return an error if something goes wrong here
	// Local node info is cached
	if LocalNode.Id > 0 && len(LocalNode.Token) == 40 {
		return /* all is good */
	}

	var nde node.Node

	if config.Opts.SynchClient != "" { // we are a client
		if km_db.Db.Where("is_local = 1").First(&nde); nde.Id < 1 { // create the signature
			km_db.Db.Create(&node.Node{Guid: util.Random_sha1(), Token: util.Random_sha1(), IsLocal: 1})
			if km_db.Db.Where("is_local = 1").First(&nde); nde.Id > 0 && len(nde.Token) == 40 { // was it saved?
				util.Pl("Local signature created")
			}
		} /*else {
			util.Pl("Local db signature already exists")
		}*/
	} else { // server
		cl := kredis.GetClient()
		if cl == nil {
			// TODO handle this error condition
			return
		}

		byts, err := roredis.GetBytes(cl, serverKey)
		if err != nil {
			util.Pl("Error obtaining server node info", err)
			return
		}

		err = json.Unmarshal(byts, &nde)
		if err != nil {
			util.Pl("Error unmarshalling server node info", err)
			return
		}

		if nde.Guid == "" {
			guid := ""
			if v := os.Getenv("KM_SERVER_GUID"); v != "" {
				guid = v
			} else {
				guid = util.Random_sha1()
			}
			nde.Guid = guid

			if nde.Token == "" {
				token := ""
				if v := os.Getenv("KM_SERVER_TOKEN"); v != "" {
					token = v
				} else {
					token = util.Random_sha1()
				}
				nde.Token = token
			}

			nde.IsLocal = 1

			// Save it back to Redis
			byts, err = json.Marshal(nde)
			if err != nil {
				util.Pl("Error marshalling server node")
				return
			}

			err = roredis.Set(cl, serverKey, string(byts), 0)
			if err != nil {
				util.Pl("Error writing server node to redis")
				return
			}
		}
	}

	// Cache local copy
	LocalNode = nde
}
