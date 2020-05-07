package htb

import (
    "../config"
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "time"
    "regexp"
    "strconv"
    "sync"
)

var(
    mu = &sync.Mutex{}
)

func StartRefreshChallenges(ticker *time.Ticker){
    RefreshChallenges()
    for{
        select {
            case <- ticker.C:
                RefreshChallenges()
        }
    }
}

func RefreshChallenges() {
/*
function to create boxes.json with all the HTB box
TODO
- add verification to prevent rewrite
- add lock to write ?
- Add Real difficulty / Rate
*/
    categories := []string{"Reversing", "Crypto", "Stego", "Pwn", "Web", "Misc", "Forensics", "Mobile", "OSINT"}
    

    var challs []config.Challenge

    // Parse users profil asynchronously
    var wg sync.WaitGroup

    for _, category := range categories{
        wg.Add(1)
        go parse(&wg, &challs, category)
    }

    wg.Wait()   

    // Create file with new data
    data, _ := json.Marshal(challs)
    err := ioutil.WriteFile("challs.json", data, 0644)
    if err != nil{
        fmt.Println("[!] error in challenges : cannot create file")
    }

    return
}


func parse(wg *sync.WaitGroup, challs *[]config.Challenge, category string) {
    defer wg.Done() 
    
    client := &http.Client{
        Timeout: time.Second * 15,
        Jar: config.Htbcookies,
    }

    req, _ := http.NewRequest("GET", "https://www.hackthebox.eu/home/challenges/"+ category, nil)
    req.Header.Add("User-Agent", config.USERAGENT)
    resp, err := client.Do(req)
    if err != nil{
        fmt.Println(err)
        return
    }
    if resp.StatusCode != 200{
        fmt.Println("[!] RefreshChallenges, error all")
        return
    }
    defer resp.Body.Close()

    // Read response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        print(err)
    }

    // Regex doesn't work for retired challs, use this for retired as well :
    // panel-tools\"> (\d*\/\d*\/\d*) (?:.*?)\"text-(success|warning|danger)\">(?:.*?)(?:\[(\d*?) Points\]) <\/span> (.*?) \[by <(?:.*?)>(.*?)<\/a>\](?:.*?)\[(\d*?) solvers\](?:.*?)challenge=\"(.*?)\" data-toggle=(?:.*?)Rate Pro\">(\d*?) <(?:.*?)Rate Sucks\">(\d*?) <(?:.*?)> First Blood: <(?:.*?)>(.*?)<(?:.*?)><\/span><br><br>([\s\S]*?)<br> <br> (?:<p|<\/div)
    r := regexp.MustCompile(`panel-tools\"> (\d*\/\d*\/\d*) (?:.*?)\"text-(success|warning|danger)\">(?:.*?)(?:\[(\d*?) Points\]) <\/span> (.*?) \[by <(?:.*?)>(.*?)<\/a>\](?:.*?)\[(\d*?) solvers\](?:.*?)challenge=\"(.*?)\" data-toggle=(?:.*?)Rate Pro\">(\d*?) <(?:.*?)Rate Sucks\">(\d*?) <(?:.*?)> First Blood: <(?:.*?)>(.*?)<(?:.*?)><\/span><br><br>([\s\S]*?)<br> <br> (?:<p|<\/div)`)

    html := string(body)
    match := r.FindAllStringSubmatch(html, -1)

    for _, m := range(match){
        var chall config.Challenge
        chall.Release = m[1]
        chall.Difficulty = getDifficulty(m[2])
        chall.Points = m[3]
        chall.Name = m[4]
        chall.Maker = m[5]
        chall.Owns = m[6]
        chall.ID, _ = strconv.Atoi(m[7])
        chall.Rates.Pro = m[8]
        chall.Rates.Sucks = m[9]
        chall.Blood = m[10]
        chall.Description = m[11]
        chall.Category = category

        mu.Lock()
        *challs = append(*challs, chall)
        mu.Unlock()
    }

    return
}

func getDifficulty(text string) string{
    switch text{
        case "success":
            return "Easy"
        case "warning":
            return "Medium"
        case "danger":
            return "Hard"
        default:
            return "Unknown"
    }
}