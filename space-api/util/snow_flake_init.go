package util

import (
	"log"
	"space-api/conf"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	nodeCreated, err := snowflake.NewNode(int64(conf.GetProjectViper().GetInt("app.nodeID")))
	if err != nil {
		log.Fatal("could not get the id generator: ", err)
	} else {
		node = nodeCreated
	}
}

func GetSnowFlakeNode() *snowflake.Node {
	return node
}
