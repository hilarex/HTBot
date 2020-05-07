package htb

import (
    "../config"

    "fmt"
    "math"
)


func GetBoxInfo(box config.Box) map[string]string{
/*
Get info on box (difficulty, difficultyRating, status)
*/
    
    // status
    var status string
    if box.Retired == true{
        status = "Retired"
    }else{
        status = "Active"
    }

    var difficulty string
    switch box.Points{
        case 20:
            difficulty = "Easy"
        case 30:
            difficulty = "Medium"
        case 40:
            difficulty = "Hard"
        case 50:
            difficulty = "Insane"
        default:
            difficulty = "?"
    }

    var difficultyRating float64
    var count int
    for i, d := range box.Difficulty{
        difficultyRating +=  float64(i) * float64(d)
        count += d
    }
    difficultyRating /= float64(count)

    info := map[string]string{
        "status" : status,
        "difficulty" : difficulty,
        "difficultyRating" : fmt.Sprintf("%.1f", math.Round(difficultyRating*10)/10),
    }

    return info
}