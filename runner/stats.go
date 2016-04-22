package runner

/*
import (
	"coolstorybot/survey"
	"fmt"
)

type surveyStats struct {
	s           *survey.Survey
	qsFractions []map[string]float32
}

// TODO add question stat entity later
func newSurveyStats(s *surveyItem) *surveyStats {
	ss := &surveyStats{
		s.sv, make([]map[string]float32, len(s.sv.Questions)),
	}

	total := float32(len(s.convs))
	for qi, _ := range s.sv.Questions {
		textAnswersCount := make(map[string]uint)
		for _, c := range s.convs {
			a := c.prog.answers[qi]
			_, isTxtAExists := textAnswersCount[a.Text]
			if !isTxtAExists {
				textAnswersCount[a.Text] = 0
			}
			textAnswersCount[a.Text] += 1
		}
		fs := make(map[string]float32)
		for t, tc := range textAnswersCount {
			fs[t] = float32(tc) / total
		}
		ss.qsFractions[qi] = fs
	}

	return ss
}

func (ss *surveyStats) infoString() string {
	s := ""
	for qi, q := range ss.s.Questions {
		s += fmt.Sprintf("\n%s", q.Text)
		for t, tc := range ss.qsFractions[qi] {
			s += fmt.Sprintf("\n%.1f%s ответили [%s]", tc*100., "%", t)
		}
	}
	return s
}
*/
