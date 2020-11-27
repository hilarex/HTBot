package cmd

import (
	"../framework"
	"../config"

	"strings"
)

func RoleCommand(ctx framework.Context) {
    
    if len(ctx.Args) == 0{
        ctx.Reply("Choose a role to add : htb, paris, rennes, tours, lille, canada")
        return
    }
    
    roleInput := strings.ToLower(ctx.Args[0])
    remove := false

    if roleInput[0] == '!'{
        roleInput = roleInput[1:]
        remove = true
    }
    if ! framework.IsInSlice(roleInput, []string{"htb", "paris", "rennes", "tours", "lille", "canada"}){
    	ctx.Reply("I don't know this role..")
    	return
    }

    roles, _ := ctx.Discord.GuildRoles(config.Discord.GuildID)
    roleID := ""

    for _, role := range roles{
    	if roleInput == strings.ToLower(role.Name){
    		roleID = role.ID
            break
    	}

        if strings.ToLower(role.Name) == "htb player" && roleInput == "htb"{
            roleID = role.ID
            break
        }
    }

    if remove {
        ctx.Discord.GuildMemberRoleRemove(config.Discord.GuildID, ctx.User.ID, roleID)
        ctx.Reply("And it's over... 🔕")
    } else {
        member, _ := ctx.Discord.GuildMember(config.Discord.GuildID, ctx.User.ID)
        if !framework.IsInSlice(roleID, member.Roles){
            ctx.Discord.GuildMemberRoleAdd(config.Discord.GuildID, ctx.User.ID, roleID)
            ctx.Reply("You got promoted ! 🍻")
        }
    }

    return
}