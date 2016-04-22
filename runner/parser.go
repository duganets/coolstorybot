package runner

import (
	"coolstorybot/survey"
	"errors"
	"fmt"
	"strings"
)

var (
	usersMinNumber    = 1
	usersMaxNumber    = 5
	questionMaxLength = 1024

	keywordNewSurvey = "спроси"

	perrParserError     = errors.New("Ошибка парсинга, ожидаю сообщение в формате [Спроси @user1 @user2 ... @userN что спросить]")
	perrNoKwFound       = errors.New(fmt.Sprintf("Сообщение должно начинаться с ключевого слова [%s]", keywordNewSurvey))
	perrNotEnoughUsers  = errors.New(fmt.Sprintf("Минимальное количество пользователей = %d", usersMinNumber))
	perrTooManyUsers    = errors.New(fmt.Sprintf("Максимальное количество пользователей = %d", usersMaxNumber))
	perrNoQuestion      = errors.New("Не найден вопрос")
	perrQuestionTooLong = errors.New(fmt.Sprintf("Максимальная длина вопроса = %d", questionMaxLength))
)

func ParseSurvey(surveyText string, owner string, userId2Name map[string]string) (resSrv *survey.Survey, perr error) {
	var userId string
	surveyText = strings.Trim(surveyText, " ")
	words := strings.Split(surveyText, " ")
	if len(words) < 3 {
		perr = perrParserError
		return
	}
	kw := strings.ToLower(words[0])
	if kw != keywordNewSurvey {
		perr = perrNoKwFound
		return
	}

	users := make([]string, 0)
	qwords := make([]string, 0)
	isScanUsersDone := false

	words = words[1:]
	for _, w := range words {
		w = strings.Trim(w, " ")
		if w == "" {
			continue
		}
		if !isScanUsersDone {
			if string([]byte{w[0]}) == "<" && 3 < len(w) {
				userId = w[2:(len(w) - 1)]
				if userId != "" {
					uname, uExs := userId2Name[userId]
					if uExs {
						users = append(users, uname)
						delete(userId2Name, userId)
					}
				} /* else {
					fmt.Printf("user [%s][%s] not found", w, userId)
					perr = perrParserError
					return
				}*/
				continue
			} else {
				isScanUsersDone = true
			}
		}

		if isScanUsersDone {
			qwords = append(qwords, w)
		}
	}

	if len(users) < usersMinNumber {
		perr = perrNotEnoughUsers
		return
	}
	if usersMaxNumber < len(users) {
		perr = perrTooManyUsers
		return
	}

	q := strings.Join(qwords, " ")

	if q == "" {
		perr = perrNoQuestion
		return
	}
	if questionMaxLength < len(q) {
		perr = perrQuestionTooLong
		return
	}

	if q[(len(q)-1):] != "?" {
		q = q + "?"
	}

	resSrv = survey.NewSurvey("Вопрос от ["+owner+"]",
		owner,
		q, users,
	)

	return
}
