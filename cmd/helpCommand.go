package cmd

import (
	"../framework"
	"../config"
)

func HelpCommand(ctx framework.Context) {

	ctx.Reply( "```LaPiraterie's Bot"+`

Command     Options             Description                              
----------  ---------------     --------------------------------------------
ping                            Want to ping pong ?
echo        <sentence>          A simple echo command
help                            Shows this message

ippsec      <search_term>       Search through Ippsec videos

get_box     <box_name>          Get info on a box
get_chall   <chall_name>        Get info on a chall
get_user    <htb_user_name>     Stalk your competitors
last_box                        Get info on the newest box
list_boxes                      list all active boxs
list_challs <category>          list active challs by category

verify      <api_token>         Verify your HTB account (restricted)
me                              Get your HTB info (restricted)
leaderboard                     Get the leaderboard of the guild (restricted)
----------- -----------------   --------------------------------------------
Type `+config.Prefix+`help command for more info on a command.
You can also type `+config.Prefix+"help category for more info on a category."+"```" )
	
	return
}