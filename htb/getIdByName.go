package htb

import (
    "../config"

    "net/http"
    "net/url"
    "io/ioutil"
    "time"
    "encoding/json"
    "strings"
    "fmt"
)

type User struct{
    Id        int    `json:"id"`
    Username  string `json:"username"`
}

func GetIdByName(username string) int{
/*
Get ID of user by his username
*/
    client := &http.Client{
        Timeout: time.Second * 10,
    }
    
    // Get request for userID
    data := url.Values{}
    data.Set("username", username)
    req, err := http.NewRequest("POST", "https://www.hackthebox.eu/api/user/id?api_token="+config.Htb.ApiToken, strings.NewReader(data.Encode()))
    req.Header.Add("User-Agent", config.USERAGENT)
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    resp, err := client.Do(req)
    if err != nil {
        print(err)
        return 0
    }
    defer resp.Body.Close()
    
    // If valid
    if resp.StatusCode != 200{
        fmt.Println("[!] error getting user id of:", username)
        return 0
    }

    // Read response
    body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
        print(err)
        return 0
    }

    var result User
    json.Unmarshal(body, &result)

    return result.Id
}