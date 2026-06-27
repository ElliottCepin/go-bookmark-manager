package domain

import "time"

type Bookmark struct {
	Id int64
	URL string
	Title string
	Tags []string
	Time time.Time 
}



