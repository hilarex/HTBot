package htb

import (
    "../config"
    
    "fmt"
    "io/ioutil"
    "net/http"
    "net/http/cookiejar"

    "regexp"
    "net/url"
    "time"
    "strings"
)

func StartLogin(ticker *time.Ticker){
    for {
        select {
            case <- ticker.C:
                Login()
        }
    }
}

func Login() {

	// Create cookie jar
	jar, _ := cookiejar.New(nil)
/*	proxyUrl, err := url.Parse("http://127.0.0.1:8080")
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        Proxy: http.ProxyURL(proxyUrl),
    }
*/
	client := &http.Client{
  		Timeout: time.Second * 10,
  		Jar: jar,
// 		Transport: tr,
	}
	
	// Request for csrf token
	req, err := http.NewRequest("GET", "https://www.hackthebox.eu/login", nil)
	req.Header.Add("User-Agent", config.USERAGENT)
	resp, err := client.Do(req)
    if err != nil {
        print(err)
        return
    }
    defer resp.Body.Close()
    
 	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        print(err)
        return
    }

    // find token 
    r, _ := regexp.Compile("type=\"hidden\" name=\"_token\" value=\"(.+?)\"")
    token := r.FindStringSubmatch(string(body))
    if len(token) == 0{
        return
    }
	crsf_token := token[1]

    // Post request to login
    params := url.Values{}
	params.Set("_token", crsf_token)
	params.Set("email",  config.Htb.Email)
	params.Set("password", config.Htb.Password)
	postData := strings.NewReader(params.Encode())
    req, err = http.NewRequest("POST", "https://www.hackthebox.eu/login", postData)
	req.Header.Add("User-Agent", config.USERAGENT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") 
	
	resp, err = client.Do(req)
	defer resp.Body.Close()
    if err != nil {
        print(err)
        return
    }
    if resp.StatusCode != 200{
    	fmt.Println("Error connecting to HTB")
        return
    }

    config.Htbcookies = jar

    return
}