package main
import(
	"errors"
	"strings"
)

const authFailMsg = "Authentication failure. Generate authorization token with -synch_auth\nThen store in node entry on client with -store_synch_auth"
var LocalNode Node // cache the local node

// The low-down on auth.
// Each node will have a Node table
// Each record will store the GUID of a node, and an associated token
//  which is the token required to authenticate with the node,
// or the server's secret token if the node is the local machine
//  depending on the setting of the is_local field
//

// Get local DB signature
func whoAmI() string {
	var node Node
	var err error
	if node, err = get_local_node(); err != nil {
		return ""
	}
	return node.Guid
}

func get_server_secret() string {
	var node Node
	var err error
	if node, err = get_local_node(); err != nil {
		return ""
	}
	return node.Token
}

func get_local_node() (Node, error) {
	if LocalNode.Id > 0 { // it has been inited
		return LocalNode, nil
	}
	var node Node
	db.Where("is_local = 1").First(&node)
	if node.Id < 1 { // no such node
		ensureDBSig()
		db.Where("is_local = 1").First(&node)
		if node.Id < 1 {
			str := (`Could not locate or create local database signature.
			Delete the local database, and try again`)
			fpl(str)
			return node, errors.New(str)
		}
	}
	LocalNode = node // cache
	return node, nil
}

// We no longer create Node here
// since node needs to have been created to have an auth token
func getNodeByGuid(node_id string) (Node, error) {
	var node Node
	db.Where("guid = ?", node_id).First(&node)
	if node.Id < 1 {
		return node, errors.New("Could not create node")
	}
	return node, nil
}

// Get node token or create node entry and token on the server and return node token
// (If there is already a valid one for this node, use that)
func getNodeToken(guid string) (string, error) {
	var node Node
	db.Where("guid = ?", guid).First(&node)
	if node.Id < 1 {
		token := random_sha1()
		db.Create(&Node{Guid: guid, Token: token})
		pl("Creating new node entry for:", short_sha(guid))
		db.Where("guid = ?", guid).First(&node) // read it back
		if node.Id < 1 {
			return "", errors.New("Could not create node entry")
		} else {
			return token, nil
		}
	  // Node already exists - make sure it has an auth token
	} else if len(node.Token) == 0 {
		token := random_sha1()
		node.Token = token
		db.Save(&node)
		return token, nil
	} else {
		return node.Token, nil
	}
}

// Use this method for manual generation of a token for the client
// The client will save the token for later access to the server node
// The client saves the Node Server's Guid along with the required auth token for that node.
func saveNodeToken(compound string) {
	arr := strings.Split(strings.TrimSpace(compound), "-")
	guid, token := arr[0], arr[1]
	pf("Node: %s, Auth Token: %s\n", guid, token)
	err := setNodeToken(guid, token)
	if err != nil { pl(err) }
}

// The client saves the token required to dial the node whose GUID is guid
func setNodeToken(guid string, token string) (error) {
	var node Node
	db.Where("guid = ?", guid).First(&node)
	if node.Id < 1 {
		pl("Creating new node entry for:", short_sha(guid))
		db.Create(&Node{Guid: guid, Token: token})
		// Verify
		db.Where("guid = ?", guid).First(&node)
		if node.Id < 1 {
			return errors.New("Could not create node entry")
		}
	} else { // Node already exists - make sure it has an auth token
		node.Token = token // always update
		db.Save(&node)
		pf("Updated token for node entry: %s", short_sha(guid))
	}
	return nil
}
