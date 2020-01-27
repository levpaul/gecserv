package eventbus

import "github.com/levpaul/idolscape-backend/internal/fb"

var pipeErr chan<- error

type player struct {
}

func Start(pErr chan<- error) error {
	pipeErr = pErr
	go start()
	return nil
}

func start() {
	for {
		select {}
	}
}

type InputEvent struct {
	Topic uint64
	Data  fb.PlayerInputT
}

func AddEventTopic(topicID uint64) {

}
func RemoveEventTopic(topicID uint64) {

}

//func GetEventTopics() []

//func SubmitEvent
