package cmd

import (
	"../framework"
	"../config"

	"encoding/json"
	"io/ioutil"
	"fmt"
)

func LastBoxCommand(ctx framework.Context) {
/*
Send info about the last HTB box
TODO :
	- move the reply function to context
	- use : https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds 
*/

 	// Read json data from file
    plan, err := ioutil.ReadFile("boxes.json")
    if err != nil {
		fmt.Println(err.Error())
	}
	
	var data []config.Box
	err = json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Println(err.Error())
	}
	
    // Get the last box
	lastBox := data[len(data)-1]

    // Function to print response (in getBoxCommand.go)
    ReplyBoxInfo(&ctx, &lastBox)
}