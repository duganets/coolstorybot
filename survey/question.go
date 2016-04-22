package survey

type Question struct {
	text string
}

func NewQuestion(q string) *Question {
	qs := &Question{q}
	return qs
}
func (q *Question) Text() string {
	return q.text
}
