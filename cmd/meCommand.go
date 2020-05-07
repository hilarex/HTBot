package cmd

import (
	"../framework"
	"../config"
	"../htb"

	"github.com/bwmarrin/discordgo"
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
	htb.ParseUserProfil(&user)

    var team string = user.Team
 	if team != ""{
 		team = " | 🏡 " + user.Team
 	}
 	var vip string
 	if user.VIP == true{
 		vip = "  💠"
 	}

	embed := &discordgo.MessageEmbed{
    Color:       0x69c0ce, 
    Description: fmt.Sprintf("🎯 %v • 🏆 %v • 👤 %v • ⭐ %v", user.Points, user.Systems, user.Users, user.Respect),
    Fields: []*discordgo.MessageEmbedField{
        &discordgo.MessageEmbedField{
            Name:   "About",
            Value:  fmt.Sprintf("📍 %v | 🔰 %v%v\n\n**Ownership** : %v%% | **Rank** : %v | ⚙️ **Challenges** : %v", user.Country, user.Level, team, user.Ownership, user.Rank, user.Challs),
            Inline: true,
        },
    },
   	Thumbnail: &discordgo.MessageEmbedThumbnail{
        URL: user.Avatar,
    },
    Title:   user.Username + vip,
	}

	ctx.ReplyEmbed( embed )
}