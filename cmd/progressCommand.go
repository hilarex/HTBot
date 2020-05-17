package cmd

import (
    "../framework"
    "../config"

    "encoding/json"
    "io/ioutil"
    "github.com/bwmarrin/discordgo"
    "fmt"
    "strings"
)

func ProgressCommand(ctx framework.Context) {
/*
View progress of HTB Team player
*/

    if !IsMemberOfTeam(ctx.Discord, ctx.User.ID){
        ctx.Reply("Sorry, you're not in the team, you cannot see progress")
        return
    }

    if len(ctx.Args) == 0{
        ctx.Reply("Give me a box name")
        return
    }

    // Read json data from file
    data, err := ioutil.ReadFile("progress.json")
    if err != nil {
        fmt.Println(err.Error())
    }
    var progress []config.Progress
    err = json.Unmarshal(data, &progress)
    if err != nil {
        fmt.Println(err.Error())
    }
    if len(progress) == 0 {
        return
    }

    
    // Get players name that did the box
    box := strings.ToLower(ctx.Args[0])
    var player_user, player_root string

    for _, user := range progress{
        if framework.IsInSlice(box, user.Users){
            player_user += user.Username + "\n"
        }
        if framework.IsInSlice(box, user.Roots){
            player_root += user.Username + "\n"   
        }
    }

    // if no one did it
    if player_user == ""{
        player_user = "-"
    }
    if player_root == ""{
        player_root = "-"
    }


    // Send response
    embed := &discordgo.MessageEmbed{
    Color:       0x69c0ce, 
    Fields: []*discordgo.MessageEmbedField{
        &discordgo.MessageEmbedField{
            Name:   "Own User",
            Value:  player_user,
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Own Root",
            Value:  player_root,
            Inline: true,
        },
    },
    Title: "📊 Progress for "+ctx.Args[0],
    }

    ctx.ReplyEmbed( embed )

    return
}

