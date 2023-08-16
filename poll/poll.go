package poll

type PollQuestion struct {
	Id string `json:"id"`
  Question string `json:"question"`
  Options [2]string `json:"options"`
  submissions [2]int
}

func (p *PollQuestion) AddSubmissions(option1Count int, option2Count int) {
  p.submissions = [2]int {option1Count, option2Count}
}