package ippsec

import(
	"../config"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strings"
)

func min(i1 int, i2 int) int{
	if (i1 < i2){
		return i1
	}else{
		return i2
	}
}

func SearchIppsec(search string, page int) (string, int){

	if page < 0{
		return "", 0
	}	

	// Read json data from file
    data, err := ioutil.ReadFile("ippsec.json")
    if err != nil {
		fmt.Println(err.Error())
	}
	
	var videos []config.Video
	err = json.Unmarshal(data, &videos)
	if err != nil {
		fmt.Println(err.Error())
	}
	
	var allMatch []config.Video
	for _, v := range videos{
		content := strings.ToLower(fmt.Sprintf(v.Machine, v.Line))
		match := true
		for _, term := range strings.Split(search," "){
			if !strings.Contains( content , strings.ToLower(term)){
				match = false
				break
			}
		}
		if match{
			allMatch = append(allMatch, v)
		}
	}

	if len(allMatch) == 0{
		return "", 0
	}

	lenPage := 6
	var result []config.Video
	result = allMatch[ min(lenPage*(page-1), len(allMatch)) : min(lenPage*(page), len(allMatch)) ]
	numPage := len(allMatch)/lenPage + 1

	var text string
	for _, r := range result{
		text += fmt.Sprintf("**%v** *(%02d:%02d:%02d)*\n%v\nhttps://youtube.com/watch?v=%v&t=%d\n\n", r.Machine, r.Timestamp.Minutes/60, r.Timestamp.Minutes - (r.Timestamp.Minutes/60)*60, r.Timestamp.Seconds, r.Line, r.VideoId, (r.Timestamp.Minutes*60+r.Timestamp.Seconds))
	}

	return text, numPage
}