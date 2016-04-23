package runner

import (
	"coolstorybot/survey"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"strings"
)

const debugSlack = false

var surveysCounter = uint(0)

type SurveysRunner struct {
	slackToken  string
	rtm         *slack.RTM
	users       map[string]*slack.User
	usersByName map[string]*slack.User
	userId2Name map[string]string
	channels    map[string]*slack.Channel
	usersConvs  map[string]*conversation
	chMsg       chan *slack.MessageEvent
	chConvSend  chan conversationMsgOut
	chNew       chan *survey.Survey
	chConvDone  chan conversationItem
	surveys     map[string]*survey.Survey
	schS        *survey.Survey
	ChDone      chan bool
}

func New(slackToken string) *SurveysRunner {
	sr := &SurveysRunner{
		slackToken,
		nil, // rtm
		make(map[string]*slack.User, 0),
		make(map[string]*slack.User, 0),
		make(map[string]string, 0),
		make(map[string]*slack.Channel, 0),
		make(map[string]*conversation, 0),
		make(chan *slack.MessageEvent, 32),
		make(chan conversationMsgOut, 32),
		make(chan *survey.Survey, 32),
		make(chan conversationItem, 32),
		make(map[string]*survey.Survey, 0),
		nil,
		make(chan bool),
	}
	sr.initSlack()
	go sr.run()
	return sr
}

func (sr *SurveysRunner) ScheduleSurvey(sv *survey.Survey) {
	sr.schS = sv
}
func (sr *SurveysRunner) AddSurvey(sv *survey.Survey) {
	sr.chNew <- sv
}

func (sr *SurveysRunner) send(outMsg conversationMsgOut) {
	sr.chConvSend <- outMsg
}
func (sr *SurveysRunner) convDone(ci conversationItem) {
	sr.chConvDone <- ci
}
func (sr *SurveysRunner) onMessage(msg *slack.MessageEvent) {
	sr.chMsg <- msg
}
func (sr *SurveysRunner) online() {
	log.Print("SurveysRunner.online")
	if len(sr.users) == 0 {
		sr.loadUsers()
	}
	if len(sr.channels) == 0 {
		sr.loadChannels()
	}
	if len(sr.usersConvs) == 0 {
		sr.initConversations()
	}

	if sr.schS != nil {
		sr.AddSurvey(sr.schS)
		sr.schS = nil
	}
}
func (sr *SurveysRunner) handleUserMessage(inMsg *slack.MessageEvent) {
	u := sr.users[inMsg.User]
	usrsMap := make(map[string]string, 0)
	for k, v := range sr.userId2Name {
		usrsMap[k] = v
	}
	newSurvey, nserr := ParseSurvey(inMsg.Text, u.Name, usrsMap)
	replyText := ""
	if nserr != nil {
		replyText = nserr.Error()
		sr.send(conversationMsgOut{u.Name, replyText})
	} else {
		sr.AddSurvey(newSurvey)
	}
}
func (sr *SurveysRunner) run() {
	log.Printf("SurveysRunner started run loop")
	for {
		select {
		case ci := <-sr.chConvDone:
			log.Print("SurveysRunner conversation done chan item process")
			for _, si := range sr.surveys {
				if si.Id == ci.id {
					si.AddAnswer(ci.user, ci.in)
					if si.IsDone() {
						log.Printf("SurveysRunner survey[%s] is done", si.Title)
						sr.send(conversationMsgOut{si.Owner, si.Info()})
						sr.surveys[si.Owner] = nil
						delete(sr.surveys, si.Owner)
					}
					break
				}
			}
		case sv := <-sr.chNew:
			_, isHaveRunning := sr.surveys[sv.Owner]
			if isHaveRunning {
				sr.send(conversationMsgOut{sv.Owner, "ваш предыдущий опрос еще не закончился"})
			} else {
				sr.initSurvey(sv)
				replyText := fmt.Sprintf("Спрошу у %s\n%s", strings.Join(sv.Users, ", "), sv.Question.Text())
				sr.send(conversationMsgOut{sv.Owner, replyText})
			}
		case inMsg := <-sr.chMsg:
			for uid, c := range sr.usersConvs {
				if uid == inMsg.User {
					if 0 < len(c.itemsQ) {
						log.Printf("SurveysRunner onMessage conversation found, uid=[%s]", uid)
						c.reply(inMsg.Text)
					} else {
						sr.handleUserMessage(inMsg)
					}
				}
			}
			log.Print(inMsg)

		case outMsg := <-sr.chConvSend:
			u := sr.usersByName[outMsg.username]
			_, _, chId, _ := sr.rtm.OpenIMChannel(u.ID)
			sr.rtm.SendMessage(sr.rtm.NewOutgoingMessage(outMsg.message, chId))

		}
	}
	log.Printf("SurveysRunner finished run loop")
}
func (sr *SurveysRunner) initSurvey(sv *survey.Survey) {
	surveysCounter += 1
	log.Print("initSurvey, surveysCounter=", surveysCounter)
	sv.Id = surveysCounter
	sr.surveys[sv.Owner] = sv
	for _, un := range sv.Users {
		log.Print("un=", un)
		u := sr.usersByName[un]
		c := sr.usersConvs[u.ID]
		c.addItem(sv.Id, sv.Question.Text(), "["+sv.Owner+" интересуется] %s")
	}
}
func (sr *SurveysRunner) loadUsers() {
	usersList, chErr := sr.rtm.GetUsers()
	if chErr != nil {
		log.Print("loadUsers error, GetUsers failed", chErr)
	} else {
		for ui, u := range usersList {
			if u.IsBot || u.ID == "USLACKBOT" {
				continue
			}
			sr.users[u.ID] = &usersList[ui]
			sr.usersByName[u.Name] = &usersList[ui]
			sr.userId2Name[u.ID] = u.Name
			//log.Printf("user: %s %q", u.Name, u)
		}
		log.Print("loadUsers ok len=", len(sr.users), sr.users)
	}
}
func (sr *SurveysRunner) loadChannels() {
	chList, chErr := sr.rtm.GetChannels(true)
	if chErr != nil {
		log.Print("SurveysRunner loadChannels error, GetChannels failed", chErr)
	} else {
		for ci, c := range chList {
			sr.channels[c.Name] = &chList[ci]
			log.Printf("SurveysRunner channel: %q", c.Name)
		}
		log.Print("SurveysRunner loadChannels ok len=", len(sr.channels), sr.channels)
	}
}
func (sr *SurveysRunner) initConversations() {
	for uid, u := range sr.users {
		if u.IsBot || u.ID == "USLACKBOT" {
			log.Print("skip conversations with bot=", u.Name)
			continue
		}
		conv := newConversation(sr, u.Name)
		sr.usersConvs[uid] = conv
	}
	log.Print("SurveysRunner initConversations ok len=", len(sr.usersConvs))
}
func (sr *SurveysRunner) initSlack() {
	log.Print("SurveysRunner.initSlack")

	api := slack.New(sr.slackToken)
	//api.SetDebug(true)

	sr.rtm = api.NewRTM()
	go sr.rtm.ManageConnection()

	go func() {
		log.Print("SurveysRunner starting routing for slack events handling")
		for {
			select {
			case msg := <-sr.rtm.IncomingEvents:
				if debugSlack {
					log.Print("SurveysRunner slack event received=", msg)
				}
				switch ev := msg.Data.(type) {
				case *slack.HelloEvent:
					// Ignore hello

				case *slack.ConnectedEvent:
					//fmt.Println("Infos:", ev.Info)
					if debugSlack {
						log.Print("SurveysRunner slack connected, connection counter=", ev.ConnectionCount)
					}
					sr.online()

				case *slack.MessageEvent:
					if debugSlack {
						log.Printf("Slack message %v\n", ev)
					}
					sr.onMessage(ev)
				case *slack.PresenceChangeEvent:
					if debugSlack {
						log.Printf("Slack presence Change: %v\n", ev)
					}

				case *slack.LatencyReport:
					if debugSlack {
						log.Printf("Current latency: %v\n", ev.Value)
					}

				case *slack.RTMError:
					log.Printf("Slack errror %s\n", ev.Error())

				case *slack.InvalidAuthEvent:
					log.Printf("Slask error, invalid credentials")
					//break Loop

				default:

					// Ignore other events..
					// fmt.Printf("Unexpected: %v\n", msg.Data)
				}
			}
		}
	}()
}

/*
func (sr *SurveysRunner) reportSurvey() {
	sInfo := sr.curSurvey.infoString(sr)
	log.Print(sInfo)
	ss := newSurveyStats(sr.curSurvey)
	ssInfo := ss.infoString()
	log.Print(ssInfo)
	log.Print(sr.channels["general"])
	report := sInfo + "\n" + ssInfo

	for _, u := range sr.users {
		if sr.curSurvey.sv.IsSurveyUser(u.Name) {
			log.Print("reporting to user=", u.Name)
			sr.send(u.ID, report)
		}
	}

		sr.rtm.SendMessage(
			sr.rtm.NewOutgoingMessage(report, sr.channels["general"].ID),
		)
}*/
