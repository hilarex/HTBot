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

func GetChallCommand(ctx framework.Context) {
/*
Send info about the last HTB box
TODO :
    - Add Real Difficulty !
*/

    if len(ctx.Args) == 0{
        ctx.Reply("Which one ?")
        return
    }

    // Read json data from file
    data, err := ioutil.ReadFile("challs.json")
    if err != nil {
        fmt.Println(err.Error())
    }
    
    var challs []config.Challenge
    err = json.Unmarshal(data, &challs)
    if err != nil {
        fmt.Println(err.Error())
    }

    var chall config.Challenge
    var match int
    name := strings.ToLower(strings.Join(ctx.Args, " "))

    for _, c := range challs{
        if strings.ToLower(c.Name) == strings.ToLower(name){
            chall = c
            match = 1
            break
        }
    }

    if match == 0{
        ctx.Reply("This challenge doesn't exist..")
        return
    }


    embed := &discordgo.MessageEmbed{
    Color:       0x69c0ce, 
    Fields: []*discordgo.MessageEmbedField{
        &discordgo.MessageEmbedField{
            Name:   "Rating",
            Value:  fmt.Sprintf("👍 %s 👎 %s", chall.Rates.Pro, chall.Rates.Sucks),
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Solvers",
            Value:  fmt.Sprintf("#️⃣ %s", chall.Owns),
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Difficulty",
            Value:  fmt.Sprintf("%s (%s points)", chall.Difficulty, chall.Points),
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Release",
            Value:  chall.Release,
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Status",
            Value:  "Active",
            Inline: true,
        },
    /*    &discordgo.MessageEmbedField{
            Name:   "Real difficulty",
            Value:  ...
            Inline: true,
        },*/
         &discordgo.MessageEmbedField{
            Name:   "Description",
            Value:  chall.Description,
        },
    },
    Footer: &discordgo.MessageEmbedFooter{
        Text:         "Maker: "+chall.Maker,
    },
    Title:  fmt.Sprintf("%s (%s)", chall.Name, chall.Category),
    }

    ctx.ReplyEmbed( embed )

    return
}