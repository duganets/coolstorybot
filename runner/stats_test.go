package runner

/*
import (
	"coolstorybot/survey"
	"testing"
)

func TestSurveyStats(t *testing.T) {
	qs := []string{"Question 1"}
	usersData := []struct {
		uid, ans string
	}{
		{"uid1", "answer1"},
		{"uid2", "answer1"},
		{"uid3", "answer2"},
	}
	fracs := map[string]float32{"answer1": 2. / 3., "answer2": 1. / 3.}
	convs := make([]*conversation, 0)
	sv := survey.NewWithQsList("Sample survey", qs)
	si := &surveyItem{sv, nil}
	for _, u := range usersData {
		c := newConversation(nil, u.uid)
		c.prog.add(survey.NewAnswer(u.ans))
		convs = append(convs, c)
	}
	si.convs = convs
	ss := newSurveyStats(si)
	if len(ss.qsFractions) != len(qs) {
		t.Logf("questions number and fractions not equal %d!=%d", len(qs), len(ss.qsFractions))
		t.Fail()
	}
	for a, f := range fracs {
		if f != ss.qsFractions[0][a] {
			t.Logf("answer [%s] fraction error exp=%f act=%f", a, f, ss.qsFractions[0][a])
			t.Fail()
		}
	}
	t.Log(ss.infoString())
	t.Fail()
}
*/
