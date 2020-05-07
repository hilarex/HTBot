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

func ListBoxesCommand(ctx framework.Context) {
/*
Send info about the last HTB box
TODO :
- move the reply function to context
- use : https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds 

- Add flag remaining for user
- Add flag to display only one difficulty
*/

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

    var lineEasy, lineMedium, lineHard, lineInsane string
    var numEasy, numMedium, numHard, numInsane int

    for _, box := range boxes{
        if !box.Retired{
            
            info := htb.GetBoxInfo(box)
            var line, customEmoji string

            for _, c := range ctx.Guild.Emojis{
                if c.Name == strings.ToLower(box.Os){
                    customEmoji = c.MessageFormat()
                }
            }

            line = customEmoji +" • **"+box.Name+"** (⭐ "+box.Rating+") (🛡️ "+info["difficultyRating"]+"/10)\n"
            switch box.Points{
                case 20:
                    numEasy += 1
                    lineEasy += strconv.Itoa(numEasy)+"."+line
                case 30:
                    numMedium += 1
                    lineMedium += strconv.Itoa(numMedium)+"."+line
                case 40:
                    numHard += 1
                    lineHard += strconv.Itoa(numHard)+"."+line
                case 50:
                    numInsane += 1
                    lineInsane += strconv.Itoa(numInsane)+"."+line
            }
        }
    }

	embed := &discordgo.MessageEmbed{
    Color:       0x69c0ce, 
    Fields: []*discordgo.MessageEmbedField{
        &discordgo.MessageEmbedField{
            Name:   "Easy",
            Value:  lineEasy,
        },
        &discordgo.MessageEmbedField{
            Name:   "Medium",
            Value:  lineMedium,
        },
        &discordgo.MessageEmbedField{
            Name:   "Hard",
            Value:  lineHard,
        },
        &discordgo.MessageEmbedField{
            Name:   "Insane",
            Value:  lineInsane,
        },
    },
    Title:   "Active boxes 💻",
	}

	ctx.ReplyEmbed( embed )

}