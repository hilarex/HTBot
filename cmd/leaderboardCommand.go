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

func LeaderboardCommand(ctx framework.Context) {
/*
TODO:
- change embed to own function in context
*/

   	if !IsMemberOfTeam(ctx.Discord, ctx.User.ID){
        ctx.Reply("Sorry, you're not in the team, you cannot see the leaderboard")
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

    // Order list by points
    sort.Slice(users, func(i, j int) bool {
	    a, _ := strconv.Atoi(users[i].Points)
	    b, _ := strconv.Atoi(users[j].Points)
	    return a > b
	})

	// Create board
	var board strings.Builder
	var displayed int
	maxDisplayed := 20
	if len(users) < maxDisplayed {
		displayed = len(users)
	} else{
		displayed = maxDisplayed
	}
	for i, user := range users[:displayed] {
		switch i{
			case 0:
				board.WriteString(fmt.Sprintf("👑 **%d. %v** (Points : %v, Ownership : %v)\n", i+1, user.Username, user.Points, user.Ownership))
			case 1:
				board.WriteString(fmt.Sprintf("💠 **%d. %v** (Points : %v, Ownership : %v)\n", i+1, user.Username, user.Points, user.Ownership))
			case 2:
				board.WriteString(fmt.Sprintf("🔶 **%d. %v** (Points : %v, Ownership : %v)\n", i+1, user.Username, user.Points, user.Ownership))
			default:
				board.WriteString(fmt.Sprintf("➡ **%d. %v** (Points : %v, Ownership : %v)\n", i+1, user.Username, user.Points, user.Ownership))
		}
	}	

	// Send leaderboard
    embed := &discordgo.MessageEmbed{
        Color: 0x69c0ce,
        Description: board.String(),
        Title: "🏆 Leaderboard 🏆 | "+ctx.Guild.Name,
    }
    ctx.ReplyEmbed(embed)

    return
}
