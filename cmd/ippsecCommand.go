package cmd

import (
	"../framework"
  	"../ippsec"

  	"github.com/bwmarrin/discordgo"
	"fmt"
    "strings"
    "strconv"
)

func IppsecCommand(ctx framework.Context) {
/*
Send info about the last HTB box
*/

	if len(ctx.Args) == 0{
		ctx.Reply("Enter a search term for an Ippsec Video")
		return
	}

	result, numPage := ippsec.SearchIppsec(fmt.Sprintf(strings.Join(ctx.Args," ")), 1)

	if result == "" {
		ctx.Reply("No content matches your request..")
		return
	}

	embed := &discordgo.MessageEmbed{
    	Color:       0x69c0ce, 
 		Description: result,
	   	Title:   fmt.Sprintf(strings.Join(ctx.Args," ")),
	   	Footer: &discordgo.MessageEmbedFooter{
			Text:  "page : 1/"+strconv.Itoa(numPage),
		},
	}

	msg := ctx.ReplyEmbed( embed )

	if numPage != 1{
		ctx.Discord.MessageReactionAdd(msg.ChannelID, msg.ID, "⬅️")
		ctx.Discord.MessageReactionAdd(msg.ChannelID, msg.ID, "➡️")
	}

	return
}