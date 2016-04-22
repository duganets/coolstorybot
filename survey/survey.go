package survey

import (
	"fmt"
)

type Survey struct {
	Id           uint
	Title        string
	Owner        string
	Question     *Question
	Users        []string
	answerByUser map[string]*Answer
}

func NewSurvey(title, owner, q string, users []string) *Survey {
	s := &Survey{
		0, title, owner, NewQuestion(q),
		make([]string, len(users)),
		make(map[string]*Answer, len(users)),
	}
	for ni, n := range users {
		s.Users[ni] = n
		s.answerByUser[n] = NewEmptyAnswer()
	}
	return s
}
func (s *Survey) IsSurveyUser(name string) bool {
	_, isExs := s.answerByUser[name]
	return isExs
}
func (s *Survey) IsDone() bool {
	ansCount := len(s.answerByUser)
	for _, a := range s.answerByUser {
		if a.Text() != "" {
			ansCount -= 1
		}
	}

	if ansCount == 0 {
		return true
	} else {
		return false
	}
}
func (s *Survey) AddAnswer(user, ans string) bool {
	if s.IsSurveyUser(user) {
		a, _ := s.answerByUser[user]
		a.Set(ans)
		return true
	} else {
		return false
	}
}
func (s *Survey) Info() string {
	str := fmt.Sprintf("*** %s ***\n%s", s.Title, s.Question.Text())
	for u, a := range s.answerByUser {
		str += "\n[" + u + "] " + a.Text()
	}
	return str
}
