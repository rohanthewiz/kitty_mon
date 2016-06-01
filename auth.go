package main

import(
	"errors"
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
