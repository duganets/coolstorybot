package runner

import (
	//"coolstorybot/survey"
	"fmt"
	"log"
	"time"
)

const maxAttempts = 3

type conversation struct {
	srunner         *SurveysRunner
	username        string
	addCh           chan conversationItem
	doAskCh         chan bool
	replyCh         chan string
	timeoutCancelCh chan bool
	timeoutCh       chan bool
	itemsQ          []conversationItem
}
type conversationItem struct {
	id   uint
	user string
	out  string
	in   string
	tmpl string
}
type conversationMsgOut struct {
	username string
	message  string
}

func newConversation(srunner *SurveysRunner, username string) *conversation {
	log.Printf("conversation [%s] created", username)
	conv := &conversation{
		srunner,
		username,
		make(chan conversationItem, 32),
		make(chan bool, 32),
		make(chan string, 32),
		make(chan bool, 32),
		make(chan bool, 32),
		make([]conversationItem, 0),
	}
	go conv.run()
	return conv
}

func (conv *conversation) addItem(id uint, out string, tmpl string) {
	conv.addCh <- conversationItem{id, "", out, "", tmpl}
}
func (conv *conversation) reply(text string) {
	conv.replyCh <- text
}
func (conv *conversation) formatQ(curItem conversationItem) conversationMsgOut {
	msg := fmt.Sprintf(curItem.tmpl, curItem.out)
	return conversationMsgOut{conv.username, msg}
}
func (conv *conversation) run() {
	log.Printf("conversation [%s] started", conv.username)
	//ConvLoop:
	for {
		select {

		case item := <-conv.addCh:
			conv.itemsQ = append(conv.itemsQ, item)
			if len(conv.itemsQ) == 1 {
				conv.doAskCh <- true
			}
		case <-conv.doAskCh:
			if 0 < len(conv.itemsQ) {
				log.Printf("conversation [%s] sends question", conv.username)
				curItem := conv.itemsQ[0]
				conv.srunner.send(conv.formatQ(curItem))
				go func() {
					log.Printf("conversation [%s] timer started ", conv.username)
					select {
					case <-conv.timeoutCancelCh:
						log.Print("timer killed")
						return
					case <-time.After(time.Duration(10) * time.Minute):
						log.Printf("conversation [%s] timer triggered", conv.username)
						conv.timeoutCh <- true
					}
				}()
			}

		case ans := <-conv.replyCh:
			if 0 < len(conv.itemsQ) {
				log.Printf("conversation [%s] got reply [%s]", conv.username, ans)
				curItem := conv.itemsQ[0]
				conv.srunner.convDone(conversationItem{curItem.id, conv.username, "", ans, ""})
				conv.itemsQ = conv.itemsQ[1:]
				conv.timeoutCancelCh <- true
				conv.doAskCh <- true
			}

		case <-conv.timeoutCh:
			log.Printf("conversation [%s] timed out ", conv.username)
			if 0 < len(conv.itemsQ) {
				curItem := conv.itemsQ[0]
				conv.srunner.send(conversationMsgOut{conv.username, "истекло время ожидания ответа"})
				conv.srunner.convDone(conversationItem{curItem.id, conv.username, "", "-", ""})
				conv.itemsQ = conv.itemsQ[1:]
			}
			conv.doAskCh <- true
		}
	}
	log.Printf("conversation [%s] exits", conv.username)
}
