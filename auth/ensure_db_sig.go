package auth

import (
	"kitty_mon/km_db"
	"kitty_mon/node"
	"kitty_mon/util"
)

func EnsureDBSig() {
	if LocalNode.Id > 0 && len(LocalNode.Token) == 40 {
		return /* all is good */
	}
	var nde node.Node
	if km_db.Db.Where("is_local = 1").First(&nde); nde.Id < 1 { // create the signature
		km_db.Db.Create(&node.Node{Guid: util.Random_sha1(), Token: util.Random_sha1(), IsLocal: 1})
		if km_db.Db.Where("is_local = 1").First(&nde); nde.Id > 0 && len(nde.Token) == 40 { // was it saved?
			util.Pl("Local signature created")
		}
	} else {
		util.Pl("Local db signature already exists")
	}
}
