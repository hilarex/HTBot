package htb


import (
    "../config"

    "io/ioutil"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

func StartRefreshUsers(ticker *time.Ticker){
    RefreshUsers()
    for {
        select {
            case <- ticker.C:
                RefreshUsers()
        }
    }
}

func RefreshUsers() {
/*
TODO:
- move embed to own function
- Add mutex to write !
*/

    var users []config.User
    byteValue, err := ioutil.ReadFile("users.json")
    if err != nil{
        fmt.Println("[!] Error RefreshUsers, cannot read users.json")
        return
    }
    json.Unmarshal(byteValue, &users)
    
    // Create progress slice
    progress := make([]config.Progress, len(users))

    // Parse users profil asynchronously
    var wg sync.WaitGroup

    for i, _ := range users {
        wg.Add(1)
        go ParseUserProfil(&wg, &users[i], &progress[i])
    }

    wg.Wait()

    // Create file with new data
    data, _ := json.Marshal(users)
    err = ioutil.WriteFile("users.json", data, 0644)
    if err != nil{
        fmt.Println("[!] error RefreshUsers : cannot create users file")
        return
    }
    // Create file with new data
    data, _ = json.Marshal(progress)
    err = ioutil.WriteFile("progress.json", data, 0644)
    if err != nil{
        fmt.Println("[!] error in RefreshUsers : cannot create progress file")
        return
    }

    return
}
