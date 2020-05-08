package htb

import (
    "../config"
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "time"
)

func StartRefreshBoxes(ticker *time.Ticker){
    
    RefreshBoxes()
    for {
        select {
            case <- ticker.C:
                RefreshBoxes()
        }
    }
}

func RefreshBoxes() {
/*
function to create boxes.json with all the HTB box
TODO
- add verification to prevent rewrite
- add lock to write ?
*/

	client := &http.Client{
  		Timeout: time.Second * 10,
	}
	
    // First request for boxes info
	req, _ := http.NewRequest("GET", "https://www.hackthebox.eu/api/machines/get/all/?api_token="+ config.Htb.ApiToken, nil)
	req.Header.Add("User-Agent", config.USERAGENT)
	resp, err := client.Do(req)
    if err != nil {
        fmt.Println("[!] RefreshBoxes, error first")
        return
    }
    if resp.StatusCode != 200{
        fmt.Println("[!] RefreshBoxes, error all")
        return
    }
    defer resp.Body.Close()
    
    var boxes []config.Box
    // First response
    body, _ := ioutil.ReadAll(resp.Body)
    _ = json.Unmarshal(body, &boxes)

    // Second request for difficulty ratings
    req, _ = http.NewRequest("GET", "https://www.hackthebox.eu/api/machines/difficulty?api_token="+ config.Htb.ApiToken, nil)
    req.Header.Add("User-Agent", config.USERAGENT)
    resp, err = client.Do(req)
    if err != nil {
        fmt.Println("[!] RefreshBoxes, error second")
        return
    }
    if resp.StatusCode != 200{
        fmt.Println("[!] RefreshBoxes, error difficulty")
        return
    }
    defer resp.Body.Close()
    

    // Because there is not same the number of result in the two request
    // We create a temporary struct
    type Ratings struct{
        ID int  `json:"id"`
        Difficulty []int `json:"difficulty_ratings"`
    }
    var obj []Ratings
    body, _ = ioutil.ReadAll(resp.Body)
    _ = json.Unmarshal(body, &obj)

    // We create a map to access it by ID
    confMap := map[int][]int{}
    for _, o := range obj {
        confMap[o.ID] = o.Difficulty
    }

    // Then, we update the value of the difficulty ratings of the box if IDs match
    for i, b := range boxes{
        if val, ok := confMap[b.ID]; ok {
            boxes[i].Difficulty = val
        }
    }

    // Write to json file
    content, _ := json.Marshal(boxes)
    err = ioutil.WriteFile("boxes.json", content, 0644)
    if err != nil{
        fmt.Println("[!] RefreshBoxes, error writing file")
    }

    return
}