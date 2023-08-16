package poll

type PollQuestion struct {
	Id string `json:"id"`
  Question string `json:"question"`
  Options [2]string `json:"options"`
  submissions [2]int
}