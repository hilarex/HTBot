package cmd

import (
    "../config"
    "../framework"
    "../htb"

    "github.com/bwmarrin/discordgo"
    "net/http"
    "io/ioutil"
    "time"
    "encoding/json"
    "fmt"
    "strconv"
)

func VerifyCommand(ctx framework.Context) {
/*
TODO:
- add roles when we verify our account
- Add mutex to write !
- Add instructions in private message
*/
   
   if !IsMemberOfTeam(ctx.Discord, ctx.User.ID){
        ctx.Discord.ChannelMessageDelete(ctx.Channel.ID, ctx.Message.ID)
        ctx.Reply("Sorry, you're not in the team, you cannot get verify for now")
        return
    }

    if ctx.Channel.GuildID != ""{
        ctx.Discord.ChannelMessageDelete(ctx.Channel.ID, ctx.Message.ID)
        ctx.Reply("😱 Don't send it here !\nCome in private")
        return
    }

   if len(ctx.Args) == 0 {
        ctx.Reply("What's your Account Identifier ?\n`"+config.Prefix+"verify <account identifier>`")
        return
    }

    // Create users.json if doesn't exist 
    // and check if user already verified 
    var users []config.User
    byteValue, err := ioutil.ReadFile("users.json")
    if err != nil{
        ioutil.WriteFile("users.json", nil, 0644)
    } else {
        json.Unmarshal(byteValue, &users)
        for i := 0; i < len(users); i++ {
            id, _ := strconv.Atoi(ctx.User.ID) 
            if id == users[i].DiscordID {
                ctx.Reply("You already have verified your HTB account.")
                return
            }
        }   
    }

    // Request public info on user
    client := &http.Client{
        Timeout: time.Second * 10,
    }
    
    req, err := http.NewRequest("GET", "https://www.hackthebox.eu/api/users/identifier/"+ctx.Args[0], nil)
    req.Header.Add("User-Agent", config.USERAGENT)
    resp, err := client.Do(req)
    if err != nil {
        print(err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200{
        ctx.Reply("Oups, wrong Account Identifier..")
        return
    }

    // Read response
    body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
        print(err)
    }

    // Create user struct
    var user config.User
    err = json.Unmarshal(body, &user)
    if err != nil{
        fmt.Println(err)
    }

    // Fill user data from his HTB profil
    htb.ParseUserProfil(nil, &user, nil)
    
    // Add the discord ID
    user.DiscordID, _ = strconv.Atoi(ctx.User.ID)

    // Add the new user to the list
    users = append(users, user)

    // Create file with new data
    data, _ := json.Marshal(users)
    err = ioutil.WriteFile("users.json", data, 0644)
    if err != nil{
        fmt.Println("[!] error in verify : cannot create file")
    }

    embed := &discordgo.MessageEmbed{
        Color: 0x69c0ce,
        Description: "You are now verify 😊️",
        Title: "Congrats!",
    }
    ctx.ReplyEmbed(embed)

    // Update the users.json to prevent overwriting updated data (not optimal...)
    htb.RefreshUsers()
    return
}