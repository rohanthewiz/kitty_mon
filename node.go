package main

import (
	"time"
	"strings"
	"errors"
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
	IsLocal int  // (bool) 0 - false, 1 - true
	CreatedAt 	time.Time
	UpdatedAt	time.Time
}

// We no longer create Node here
// since node needs to have been created to have an auth token
func getNodeByGuid(node_id string) (Node, error) {
	var node Node
	db.Where("guid = ?", node_id).First(&node)
	if node.Id < 1 {
		return node, errors.New("Could not find node")
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
