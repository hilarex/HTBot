package cmd

import (
	"../framework"
	"../config"
	"../htb"

	"io/ioutil"
	"encoding/json"
	"strconv"
	"fmt"
)

func MeCommand(ctx framework.Context) {
/*
TODO:
- move embed
*/
    if !IsMemberOfTeam(ctx.Discord, ctx.User.ID){
        ctx.Reply("Sorry, you're not in the team, you cannot see the leaderboard")
        return
    }

	var users []config.User
    byteValue, err := ioutil.ReadFile("users.json")
    if err != nil{
    	fmt.Println("[!] Me command, cannot read users.json")
    }
    json.Unmarshal(byteValue, &users)
    var user config.User
    var match int
    for i := 0; i < len(users); i++ {
        id, _ := strconv.Atoi(ctx.User.ID)
        if id == users[i].DiscordID {
    		user = users[i]
    		match = 1
    		break
        }
    }

    if match == 0{
    	ctx.Reply("First you must verified your HTB account.\nPM me with the "+config.Prefix+"verify command")
    	return	
    }

	// Parse HTB profil to fill data
	htb.ParseUserProfil(nil, &user, nil)
    
    ReplyUserInfo(&ctx, &user)
    
    return
}