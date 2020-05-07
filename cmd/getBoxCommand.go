package cmd

import (
    "../framework"
    "../config"
    "../htb"

    "encoding/json"
    "io/ioutil"
    "github.com/bwmarrin/discordgo"
    "fmt"
    "strconv"
    "strings"

)

func GetBoxCommand(ctx framework.Context) {
/*
Send info about the last HTB box
TODO :
    - move the reply function to context
    - use : https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds 
*/

    if len(ctx.Args) == 0{
        ctx.Reply("Which one ?")
        return
    }

    // Read json data from file
    plan, err := ioutil.ReadFile("boxes.json")
    if err != nil {
        fmt.Println(err.Error())
    }
    
    var boxes []config.Box
    err = json.Unmarshal(plan, &boxes)
    if err != nil {
        fmt.Println(err.Error())
    }

    var box config.Box
    var match int
    name := strings.ToLower(strings.Join(ctx.Args, " "))

    for _, b := range boxes{
        if strings.ToLower(b.Name) == name{
            box = b
            match = 1
            break
        }
    }

    if match == 0{
        ctx.Reply("This box doesn't exist..")
        return
    }

    ReplyBoxInfo(&ctx, &box)
    return
}


func ReplyBoxInfo(ctx *framework.Context, box *config.Box){
    var customEmoji string

    info := htb.GetBoxInfo(*box)

    for _, c := range ctx.Guild.Emojis{
        if c.Name == strings.ToLower(box.Os){
            customEmoji = c.MessageFormat()
        }
    }

    embed := &discordgo.MessageEmbed{
    Color:       0x69c0ce, 
    Fields: []*discordgo.MessageEmbedField{
        &discordgo.MessageEmbedField{
            Name:   "IP",
            Value:  box.IP,
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "OS",
            Value:  fmt.Sprintf("%v%v", customEmoji, box.Os),
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Difficulty",
            Value: info["difficulty"]+" ("+strconv.Itoa(box.Points)+" Points)",
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Rating",
            Value: "⭐ "+box.Rating,
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Real difficulty",
            Value: fmt.Sprintf("🛡️ %v/10", info["difficultyRating"]),
            Inline: true,
        },
         &discordgo.MessageEmbedField{
            Name:   "Owns",
            Value:  fmt.Sprintf("👤 %d #️⃣󠁲󠁯󠁯󠁴󠁿 %d", box.UserOwns, box.RootOwns ),
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Status",
            Value:  info["status"],
            Inline: true,
        },
        &discordgo.MessageEmbedField{
            Name:   "Release",
            Value:  box.Release,
            Inline: true,
        },
    },
    Footer: &discordgo.MessageEmbedFooter{
        IconURL:      box.AvatarThumb,
        Text:         "Maker: "+box.Maker.Name,
    },

    Thumbnail: &discordgo.MessageEmbedThumbnail{
        URL: box.AvatarThumb,
    },
    Title:   box.Name,
    }

    ctx.ReplyEmbed( embed )
}