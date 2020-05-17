package cmd

import (
	"../framework"
	"../config"
	"../htb"

	"github.com/bwmarrin/discordgo"
	"fmt"
)

func GetUserCommand(ctx framework.Context) {

	if len(ctx.Args) == 0 {
		ctx.Reply("Give me a name\n")
		return
	}

	id := htb.GetIdByName(ctx.Args[0])
	if(id == 0){
		ctx.Reply("User doesn't exist..")
		return
	}


	var user config.User
	user.UserID = id

	// Parse HTB profil to fill data
	htb.ParseUserProfil(nil, &user, nil)

	ReplyUserInfo(&ctx, &user)

	return
}

func ReplyUserInfo(ctx *framework.Context, user *config.User){
    
    var team string = user.Team
 	if team != ""{
 		team = " | 🏡 " + user.Team
 	}
 	var vip string
 	if user.VIP == true{
 		vip = "  💠"
 	}

    var labs string
    labs = "|"
    for v := range user.Prolabs{
    	labs += fmt.Sprintf(" **%v**: %v%% |", v, user.Prolabs[v])
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
        &discordgo.MessageEmbedField{
            Name:   "Prolabs",
            Value:  labs,
        },
    },
   	Thumbnail: &discordgo.MessageEmbedThumbnail{
        URL: user.Avatar,
    },
    Title:   user.Username + vip,
	}

	ctx.ReplyEmbed( embed )
}