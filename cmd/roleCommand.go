package cmd

import (
	"../framework"
	"../config"

	"strings"
)

func RoleCommand(ctx framework.Context) {
    
    if len(ctx.Args) == 0{
        ctx.Reply("Choose a role to add : paris, rennes, tours")
        return
    }
    
    roleWanted := strings.ToLower(ctx.Args[0])
    
    if ! framework.IsInSlice(roleWanted, []string{"paris", "rennes", "tours"}){
    	ctx.Reply("I don't know this role..")
    	return
    }

    roles, _ := ctx.Discord.GuildRoles(config.Discord.GuildID)
    newRole := ""

    for _, role := range roles{
    	if strings.ToLower(ctx.Args[0]) == strings.ToLower(role.Name){
    		newRole = role.ID
    	}
    }

    ctx.Discord.GuildMemberRoleAdd(config.Discord.GuildID, ctx.User.ID, newRole)
    ctx.Reply("You got promoted ! 🍻")

    return
}