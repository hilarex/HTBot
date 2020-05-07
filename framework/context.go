package framework

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Discord      *discordgo.Session
	Guild        *discordgo.Guild
	Channel  	 *discordgo.Channel
	User         *discordgo.User
	Message      *discordgo.MessageCreate
	Args         []string
	Shoutbox     string
}

func NewContext(discord *discordgo.Session, guild *discordgo.Guild, channel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.Channel = channel
	ctx.User = user
	ctx.Message = message
	return ctx
}

func (ctx Context) Reply(content string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.Channel.ID, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}

func (ctx Context) ReplyEmbed(content *discordgo.MessageEmbed) *discordgo.Message {
	/*embed := &discordgo.MessageEmbed{
    	Color: color,
    	Description: content,
	    Title: title,
	}*/

	msg, err := ctx.Discord.ChannelMessageSendEmbed(ctx.Channel.ID, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}

func (ctx Context) ReplyShoutbox(content string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.Shoutbox, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}