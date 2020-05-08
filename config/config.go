package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http/cookiejar"
)

// Struct for config.json

type Config struct {
	Prefix      string `json:"Prefix"`
	Htb    		ConfigHtb `json:"HTB"`
	Discord     ConfigDiscord `json:"Discord"`
}
type ConfigHtb struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
	ApiToken	string `json:"api_token"`
}
type ConfigDiscord struct {
	Guild 		string `json:"guild_name"`
	Token 	    string `json:"bot_token"`
	GuildID		string `json:"guild_id"`
	Shoutbox    string `json:"shoutbox_id"`
}

// Struct for HTB json files

type User struct {
	DiscordID  int 	`json:"discord_id"`
	UserID     int    `json:"user_id"`
    VIP         bool   `json:"vip"`
    Username string `json:"user_name"`
    Avatar string `json:"avatar"`
    Points string `json:"points"`
    Systems string `json:"systems"`
    Users string `json:"users"`
    Respect string `json:"respect"`
    Country string `json:"country"`
    Team string `json:"team"`
    Level string `json:"level"`
    Rank string `json:"rank"`
    Challs string `json:"challs"`
    Ownership string `json:"ownership"`
    Prolabs map[string]string `json:"prolabs"`
}


type Maker struct{
	ID		int  `json:"id"`
	Name 	string `json:"name"`
}
type Box struct {
	ID      	int 	`json:"id"`
	Name    	string 	`json:"name"`
	Os      	string 	`json:"os"`
	IP 			string	`json:"ip"`	 
	AvatarThumb string	`json:"avatar_thumb"`	
	Points 		int		`json:"points"`
	Release 	string	`json:"release"`
	RetiredDate	string	`json:"retired_date"`
	Maker 		Maker	`json:"maker"`
	Maker2		Maker	`json:"maker2"`
	Rating 		string	`json:"rating"`
	UserOwns 	int		`json:"user_owns"`
	RootOwns	int 	`json:"root_owns"`
	Retired 	bool	`json:"retired"`
	Free		bool	`json:"free"`
	Difficulty 	[]int   `json:"difficulty_ratings"`
}

type Rate struct{
	Pro 		string `json:"pro"`
	Sucks 		string `json:"sucks"`
	Difficulty 	string `json:"difficulty"`
}
type Challenge struct{
    ID      	int 	`json:"id"`
    Name 		string  `json:"name"`
    Category 	string  `json:"category"`
    Difficulty 	string  `json:"difficulty"`
    Points  	string  `json:"points"`
    Owns    	string 	`json:"owns"`
    Rates    	Rate  	`json:"rates"`
    Release 	string 	`json:"release"`
  //  Status   	string 	`json:"status"`
    Maker    	string 	`json:"maker"`
    Blood    	string 	`json:"blood"`
    Description string 	`json:"description"`
}

type Timestamp struct{
	Minutes 	int 	`json:"minutes"`
	Seconds 	int 	`json:"seconds"`
}
type Video struct{
	Machine 	string 	`json:"machine"`
	VideoId 	string 	`json:"videoId"`
	Timestamp 	Timestamp `json:"timestamp"`
	Line 		string 	`json:"line"`
}

type Notifs struct{
	Success string 		`json:"success"`
	Html 	[]string 	`json:"html"`
}

// Global variables
var Prefix string
var Htb ConfigHtb
var Discord ConfigDiscord


// Variable for http.Client
const USERAGENT = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.85 Safari/537.36"
var Htbcookies *cookiejar.Jar

func init(){
	
	var conf Config
	values, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(values, &conf)
	
	Prefix = conf.Prefix
	Htb = conf.Htb
	Discord = conf.Discord
}