package id

import (
	"log"
	"space-api/conf"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	appConf := conf.ProjectConf.GetAppConf()
	nodeCreated, err := snowflake.NewNode(appConf.NodeID)
	if err != nil {
		log.Fatal("could not get the id generator: ", err)
	} else {
		node = nodeCreated
	}
}

func GetSnowFlakeNode() *snowflake.Node {
	return node
}
