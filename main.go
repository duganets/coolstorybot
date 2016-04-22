package main

import (
	"coolstorybot/runner"
	"flag"
	"log"
)

var flagSlackBotIntegrationKey string

func main() {
	log.Println("starting..")
	flag.StringVar(&flagSlackBotIntegrationKey, "slack_key", "", "slack integration key")
	flag.Parse()

	srunner := runner.New(flagSlackBotIntegrationKey)

	//sr := makeTest(3)

	/*sr := survey.NewSurvey("Some test survey",
		"mike",
		"How are you?", []string{"mike", "kenruska"},
	)
	srunner.ScheduleSurvey(sr)*/

	<-srunner.ChDone
}
