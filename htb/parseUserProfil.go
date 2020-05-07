package htb

import (
    "../config"

    "net/http"
    "io/ioutil"
    "time"
    "regexp"
    "strconv"
    "fmt"
)

func ParseUserProfil(user *config.User){
/*
Get information about an user by scrapping his HTB profil
*/

    client := &http.Client{
        Timeout: time.Second * 10,
        Jar: config.Htbcookies,
    }

    // Get request for userId
    req, err := http.NewRequest("GET", "https://www.hackthebox.eu/home/users/profile/"+strconv.Itoa(user.UserID), nil)
    req.Header.Add("User-Agent", config.USERAGENT)
    resp, err := client.Do(req)
    if err != nil {
        print(err)
    }
    defer resp.Body.Close()
    
    // If valid
    if resp.StatusCode != 200{
        fmt.Println("[!] error getting info on: "+strconv.Itoa(user.UserID))
        return
    }

    // Read response
    body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
        print(err)
    }
    html := string(body)

    var username, avatar, points, systems, users, respect, country, teamName, level, rank, challenges, ownership string

    // Scrapping with regex
    r := regexp.MustCompile(`<small> Member (.+?) is at position (\d+) of the Hall of Fame.`)
    match := r.FindStringSubmatch(html)
    if match != nil{
        rank = match[2]
    }
    r = regexp.MustCompile(`www.hackthebox.eu/home/teams/profile/(\d+)\">(.+?)<`)
    match = r.FindStringSubmatch(html)
    if(match != nil){
        teamName = match[2]
        //teamID = match[1]
    }
    r = regexp.MustCompile(`([a-zA-Z]+)"><span class="flag flag-`)
    match = r.FindStringSubmatch(html)
    if(match != nil){
        country = match[1]
        //teamID = match[1]
    }
    r = regexp.MustCompile("Ownership: (.+?)%")
    ownership = r.FindStringSubmatch(html)[1]
    r  = regexp.MustCompile(`>Rank: (.+?)<`)
    level = r.FindStringSubmatch(html)[1]
    r = regexp.MustCompile(`title="Points">(.+?)(\d+)<`)
    points = r.FindStringSubmatch(html)[2]
    r = regexp.MustCompile(`title="Owned Systems">(.+?)(\d+)<`)
    systems = r.FindStringSubmatch(html)[2]
    r = regexp.MustCompile(`title="Owned Users">(.+?)(\d+)<`)
    users = r.FindStringSubmatch(html)[2]
    r = regexp.MustCompile(`title="Respect">(.+?)(\d+)<`)
    respect = r.FindStringSubmatch(html)[2]
    r = regexp.MustCompile(`align="left" src="https://www.hackthebox.eu/storage/avatars/([0-9a-z]+)\.png`)
    avatar = "https://www.hackthebox.eu/storage/avatars/"+r.FindStringSubmatch(html)[1]+".png"
    r = regexp.MustCompile(`has solved (\d+) challenges`)
    challenges = r.FindStringSubmatch(html)[1]
    r = regexp.MustCompile(`class="m-n">(.+?) <`)
    username = r.FindStringSubmatch(html)[1]


    user.Username = username 
    user.Ownership = ownership
    user.Avatar = avatar
    user.Points = points
    user.Systems = systems
    user.Users = users
    user.Respect = respect
    user.Country = country
    user.Team = teamName
    user.Level = level
    user.Rank = rank
    user.Challs = challenges
    user.Ownership = ownership

    return 
}