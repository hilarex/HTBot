package framework

/*
import (
)
*/

func IsInSlice(val string, slice []string) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}