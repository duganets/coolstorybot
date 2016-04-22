package survey

type Answer struct {
	text string
}

func NewEmptyAnswer() *Answer {
	return NewAnswer("")
}
func NewAnswer(text string) *Answer {
	a := &Answer{text}
	return a
}
func (a *Answer) Set(ans string) {
	a.text = ans
}
func (a *Answer) Text() string {
	return a.text
}
