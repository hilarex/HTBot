package cmd

import(
	"../framework"
	"../config"
	"../htb"
	"../ippsec"

	"github.com/bwmarrin/discordgo"
	"fmt"
	"strings"
	"strconv"
	"time"
)


func Ready(session *discordgo.Session, event *discordgo.Ready) {

	htb.Login()

	tickerLogin 	:= time.NewTicker(30 * time.Minute)
	tickerIppsec 	:= time.NewTicker(30 * time.Minute)
	tickerUsers		:= time.NewTicker(10 * time.Minute)
	tickerBoxes		:= time.NewTicker(1  * time.Minute)
	tickerChallenges:= time.NewTicker(10 * time.Minute)
	tickerShoutbox 	:= time.NewTicker(3 * time.Second)

	go htb.StartLogin(tickerLogin)
	go htb.StartParseShoutbox(tickerShoutbox, session)
	go htb.StartRefreshUsers(tickerUsers)
	go htb.StartRefreshChallenges(tickerChallenges)
	go htb.StartRefreshBoxes(tickerBoxes)
	go ippsec.StartRefreshIppsec(tickerIppsec)
	
}

func CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	user := message.Author
	if user.ID == session.State.User.ID || user.Bot {
		return
	}
	content := message.Content
	if len(content) <= len(config.Prefix) {
		return
	}
	// Ignore all messages that don't have the prefix
	if content[:len(config.Prefix)] != config.Prefix {
		return
	}
	content = content[len(config.Prefix):]
	if len(content) < 1 {
		return
	}
	// Ignore citation message
	if string(content[0]) == " "{
		return
	}

	// Get user command
	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	
	channel, err := session.State.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel,", err)
		return
	}

	guild, err := session.State.Guild(config.Discord.GuildID)
	if err != nil {
		fmt.Println("Error getting guild,", err)
		return
	}

	// Create new context
	ctx := framework.NewContext(session, guild, channel, user, message)
	ctx.Args = args[1:]
	ctx.Shoutbox = config.Discord.Shoutbox

	// Call the corresponding handler
	switch name {
		case "echo":
			session.ChannelMessageSend(message.ChannelID, strings.Join(ctx.Args, " "))
		case "ping":
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("🏓 pong ! %v", session.HeartbeatLatency()))
		case "help":
			HelpCommand(*ctx)
		case "me":
			MeCommand(*ctx)		
		case "leaderboard":
			LeaderboardCommand(*ctx)
		case "prolab":
			ProlabCommand(*ctx)
		case "get_box":
			GetBoxCommand(*ctx)	
		case "get_chall":
			GetChallCommand(*ctx)
		case "last_box":
			LastBoxCommand(*ctx)	
		case "get_user":
			GetUserCommand(*ctx)
		case "verify":
			VerifyCommand(*ctx)
		case "list_boxes":
			ListBoxesCommand(*ctx)
		case "list_challs":
			ListChallsCommand(*ctx)
		case "ippsec":
			IppsecCommand(*ctx)
		default:
			session.ChannelMessageSend(message.ChannelID, "🤔 I don't know this command !\nFor a list of help topics, type `"+config.Prefix+"help`")
	}
}


// Handles all page related reactions for Ippsec videos
func ReactionsHandler(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	// Ignore all reactions created by the bot itself
	if reaction.UserID == session.State.User.ID {
		return
	}

	// Fetch some extra information about the message associated to the reaction
	m, err := session.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	// Ignore reactions on messages that have an error or that have not been sent by the bot
	if err != nil || m == nil || m.Author.ID != session.State.User.ID {
		return
	}

	if !isBotReaction(session, m.Reactions, &reaction.Emoji) {
		return
	}

	user, err := session.User(reaction.UserID)
	// Ignore when sender is invalid or is a bot
	if err != nil || user == nil || user.Bot {
		return
	}

	// Ignore message without embed content (like shoutbox)
	if len(m.Embeds) == 0{
		return
	}

	footer := strings.Split(m.Embeds[0].Footer.Text, "/")
	// Ensure valid footer command
	if len(footer) != 2 {
		return
	}

	currentPage, _ := strconv.Atoi(strings.Split(footer[0], " ")[2])
	totalPage, _ := strconv.Atoi(footer[1])

	// remove reaction so the user can react again
	defer session.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)

	var result string
	var newPage int
	switch reaction.Emoji.Name{
		case "⬅️":
			if (currentPage == 1){
				return
			}
			newPage = currentPage-1
			result, _ = ippsec.SearchIppsec(m.Embeds[0].Title, newPage)
		case "➡️":
			if (currentPage == totalPage){
				return
			}
			newPage = currentPage+1
			result, _ = ippsec.SearchIppsec(m.Embeds[0].Title, newPage)
	}

            
	embed := &discordgo.MessageEmbed{
    	Color:       0x69c0ce, 
 		Description: result,
	   	Title:  m.Embeds[0].Title,
	   	Footer: &discordgo.MessageEmbedFooter{
			Text:  "page : "+strconv.Itoa(newPage)+"/"+strconv.Itoa(totalPage),
		},
	}

	session.ChannelMessageEditEmbed(reaction.ChannelID, reaction.MessageID, embed)
}

// Check if users reaction is one preset by the bot
func isBotReaction(s *discordgo.Session, reactions []*discordgo.MessageReactions, emoji *discordgo.Emoji) bool {
	for _, r := range reactions {
		if r.Emoji.Name == emoji.Name && r.Me {
			return true
		}
	}

	return false
}

func IsMemberOfTeam(session *discordgo.Session, userID string) bool{
    member, err := session.State.Member(config.Discord.GuildID, userID)
    if err != nil {
        if member, err = session.GuildMember(config.Discord.GuildID, userID); err != nil {
            return false
        }
    }

    result := false
    // Iterate through the role IDs stored in member.Roles

    roles, _ := session.GuildRoles(config.Discord.GuildID)
    for _, role := range roles{
        if isInSlice(role.ID, member.Roles){
            if role.Name == "Pirates" || role.Name == "HideAndSec"{
                result = true
            }
        }
    }
    return result
}

func isInSlice(val string, slice []string) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}