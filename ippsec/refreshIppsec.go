package ippsec

import (
    "../config"
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "time"
)

func StartRefreshIppsec(ticker *time.Ticker){
    
    RefreshIppsec()
    for {
        select {
            case <- ticker.C:
                RefreshIppsec()
        }
    }
}

func RefreshIppsec() {
/*
function to create ippsec.json with all the video of ippsec
*/

	client := &http.Client{
  		Timeout: time.Second * 10,
	}
	
    // First request for boxes info
	req, _ := http.NewRequest("GET", "https://raw.githubusercontent.com/IppSec/ippsec.github.io/master/dataset.json", nil)
	req.Header.Add("User-Agent", config.USERAGENT)
	resp, err := client.Do(req)
    if err != nil {
        fmt.Println("[!] RefreshIppsec, error request")
        return
    }
    if resp.StatusCode != 200{
        fmt.Println("[!] RefreshIppsec, error code")
        return
    }
    defer resp.Body.Close()
    
    var videos []config.Video
    body, _ := ioutil.ReadAll(resp.Body)
    _ = json.Unmarshal(body, &videos)

    // Write to json file
    content, _ := json.Marshal(videos)
    err = ioutil.WriteFile("ippsec.json", content, 0644)
    if err != nil{
        fmt.Println("[!] RefreshIppsec, error writing file")
    }

    return
}