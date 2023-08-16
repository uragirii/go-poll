package poll

import (
	"math/rand"
	"strings"
	"time"
)

const USER_ID_LENGTH = 20

type PollQuestion struct {
	Id          string    `json:"id"`
	Question    string    `json:"question"`
	Options     [2]string `json:"options"`
	submissions [2]int
}

func (p *PollQuestion) AddSubmissions(option1Count int, option2Count int) {
	p.submissions = [2]int{option1Count, option2Count}
}

func (p *PollQuestion) GetSubmissions() [2]int {
	return p.submissions
}

type User struct {
	Id             string
	SubmittedPolls []string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// taken from Stack overflow answer https://stackoverflow.com/a/22892986/8077711
func RandUserId() string {
	b := make([]rune, USER_ID_LENGTH)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ParseIds(options string) []string {
	split := strings.Split(options, ",")

	return split
}

func (user *User) HasSubmitted(pollId string) bool {
	for _, id := range user.SubmittedPolls {
		if id == pollId {
			return true
		}
	}

	return false
}
