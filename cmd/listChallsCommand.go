package cmd

import (
	"../framework"
	"../config"

	"encoding/json"
	"io/ioutil"
	"github.com/bwmarrin/discordgo"
	"fmt"
    "strconv"
    "strings"
)

func ListChallsCommand(ctx framework.Context) {
/*
Send info about the last HTB box
TODO :
- move the reply function to context
- use : https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds 

- Add flag remaining for user
- Add flag to display only one difficulty
*/

    if len(ctx.Args) == 0{
        ctx.Reply("Which category ?\n(reversing, crypto, stego, pwn, web, misc, forensics, mobile, osint)")
        return
    }

    // Verify that the argument is in the list
    categoryList := []string{"reversing", "crypto", "stego", "pwn", "web", "misc", "forensics", "mobile","osint"}
    category := strings.ToLower(ctx.Args[0])
    var isIn int
    for _, c := range categoryList{
        if category == c{
            isIn = 1
            break
        }
    }
    if isIn != 1{
        ctx.Reply("This category doesn't exist")
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

    var line, lineEasy, lineMedium, lineHard string
    var numEasy, numMedium, numHard int

    for _, chall := range challs{
        if strings.ToLower(chall.Category) == category {
            line = "["+chall.Points+" Points] **"+chall.Name+"** \n"
            switch chall.Difficulty{
                case "Easy":
                    numEasy += 1
                    lineEasy += line
                case "Medium":
                    numMedium += 1
                    lineMedium += line
                case "Hard":
                    numHard += 1
                    lineHard += line
            }
        }
    }
    if len(lineEasy) == 0{lineEasy = "_"}
    if len(lineMedium) == 0{ lineMedium = "_"}
    if len(lineHard) == 0{ lineHard = "_"}

	embed := &discordgo.MessageEmbed{
    Color:       0x69c0ce, 
    Fields: []*discordgo.MessageEmbedField{
        &discordgo.MessageEmbedField{
            Name:   "Easy ("+strconv.Itoa(numEasy)+")",
            Value:  lineEasy,
        },
        &discordgo.MessageEmbedField{
            Name:   "Medium ("+strconv.Itoa(numMedium)+")",
            Value:  lineMedium,
        },
        &discordgo.MessageEmbedField{
            Name:   "Hard ("+strconv.Itoa(numHard)+")",
            Value:  lineHard,
        },
    },
    Title:   "Active Challenges 🧩",
	}

	ctx.ReplyEmbed( embed )

    return
}