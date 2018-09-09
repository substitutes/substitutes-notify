package mail

import "time"

type Update struct {
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Class     string `json:"class"`
	Name      string `json:"name"`
	Date      string `json:"date"`
}

func NewUpdate(URL, Class, Name, Date string) *Update {
	return &Update{
		URL:       URL,
		Class:     Class,
		Name:      Name,
		Date:      Date,
		Timestamp: time.Now().String(),
	}
}
