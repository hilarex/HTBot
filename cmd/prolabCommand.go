package cmd

import (
	"../framework"
	"../config"

    "github.com/bwmarrin/discordgo"
	"fmt"
	"encoding/json"
	"strconv"
	"io/ioutil"
	"sort"
	"strings"
)

func ProlabCommand(ctx framework.Context) {
/*
Get completion of prolab for the team
*/

   	if !IsMemberOfTeam(ctx.Discord, ctx.User.ID){
        ctx.Reply("Sorry, you're not in the team, you cannot see the prolab board")
        return
    }

    if len(ctx.Args) == 0{
    	ctx.Reply("Which prolab scores do you want to see ?\nRastaLabs, Offshore, Cybernetics")
    	return
    }

    // Get users list
    var users []config.User
    byteValue, err := ioutil.ReadFile("users.json")
    if err != nil{
        ctx.Reply("No HTB account register.\nType "+config.Prefix+"verify to do it !")
        return
    } 
    json.Unmarshal(byteValue, &users)

    key := strings.ToLower(ctx.Args[0])
    if users[0].Prolabs[key] == ""{
    	ctx.Reply("This lab doesn't exist")
    	return
    }


    // Remove user with 0% completion on a lab
    var usersWithLab []config.User
    for _, user := range users{
    	if user.Prolabs[key] != "0"{
    		usersWithLab = append(usersWithLab, user)
    	}
    }

    if len(usersWithLab) == 0{
    	ctx.Reply("No one has done this lab yet")
    	return
    }

    // Order list by points
    json.Unmarshal(byteValue, &users)
    sort.Slice(usersWithLab, func(i, j int) bool {
	    a, _ := strconv.ParseFloat(usersWithLab[i].Prolabs[key], 64)
	    b, _ := strconv.ParseFloat(usersWithLab[j].Prolabs[key], 64)
	    return a > b
	})
	
	// Create board
	var completed strings.Builder
	var progress strings.Builder
	var displayed int
	maxDisplayed := 20
	if len(usersWithLab) < maxDisplayed {
		displayed = len(usersWithLab)
	} else{
		displayed = maxDisplayed
	}
	for _, user := range usersWithLab[:displayed] {
		if user.Prolabs[key] == "100"{
			completed.WriteString(fmt.Sprintf("100%% : **%v**\n", user.Username))
		}else{
			progress.WriteString(fmt.Sprintf(" %v%% : **%v**\n", user.Prolabs[key], user.Username))
		}
	}	

	if completed.Len() == 0{
		completed.WriteString("-")
	}
	if progress.Len() == 0{
		progress.WriteString("-")
	}

	// Send leaderboard
    embed := &discordgo.MessageEmbed{
        Color: 0x69c0ce,
        Title: "🏆 Prolab : "+key,
        Fields: []*discordgo.MessageEmbedField{
	        &discordgo.MessageEmbedField{
	            Name:   "Completed",
	            Value:  completed.String(),
	        },
	        &discordgo.MessageEmbedField{
	            Name:   "In Progress",
	            Value:  progress.String(),
	        },
	     },
    }
    ctx.ReplyEmbed(embed)

    return
}
