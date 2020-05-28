package htb

import (
    "../config"

    "net/http"
    "io/ioutil"
    "time"
    "regexp"
    "strconv"
    "fmt"
    "strings"
    "sync"
)

func ParseUserProfil(wg *sync.WaitGroup, user *config.User, progress *config.Progress){
/*
Get information about an user by scrapping his HTB profil
*/
    if wg != nil{
        defer wg.Done()    
    }
    

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
        return
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
        return
    }
    html := string(body)

    var VIP bool
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
        //teamID = match[1]
        teamName = match[2]
    }
    r = regexp.MustCompile(`([a-zA-Z]+)"><span class="flag flag-`)
    match = r.FindStringSubmatch(html)
    if(match != nil){
        country = match[1]
    }
    r = regexp.MustCompile(`V.I.P Member`)
    match = r.FindStringSubmatch(html)
    if(match != nil){
        VIP = true
    }

    r = regexp.MustCompile(`<td class="col-md-3 text-right"> (RastaLabs|Offshore|Cybernetics) </td> <td> <div> <div data-toggle="tooltip" title="(([0-9]*[.])?[0-9]+)\%"`)
    prolabs := r.FindAllStringSubmatch(html, -1)
    tmp := make(map[string]string)
    for i,_ := range prolabs{
        tmp[ strings.ToLower(prolabs[i][1]) ] = prolabs[i][2]
    }
    user.Prolabs = tmp
    
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

    user.VIP = VIP
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

    if progress != nil{
        r = regexp.MustCompile(username+` owned (root|user) (?:.*?)(?:\d)">(.*?)<\/a>`)
        matches := r.FindAllStringSubmatch(html, -1)
        progress.Username = username
        progress.Users = nil
        progress.Roots = nil
        progress.Challs = nil
        for _, match := range matches{
            switch string(match[1]){
                case "root": 
                    progress.Roots = append(progress.Roots, strings.ToLower(match[2]))
                case "user":
                    progress.Users = append(progress.Users, strings.ToLower(match[2]))
            }
        }
    }

    return 
}