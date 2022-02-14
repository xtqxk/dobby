package dobby

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/nsqio/go-nsq"
)

type IWorker interface {
	Launch(topic string) error
}

type dobby struct {
	hdls map[string]func(message *nsq.Message) error
}

var (
	nsqRegister func(topic string, channel string, handler func(message *nsq.Message) error) error
	dLog        *log.Logger
)

func SetupNsqRegister(def func(topic string, channel string, handler func(message *nsq.Message) error) error, logger *log.Logger) error {
	if def == nil {
		return errors.New("nsq register is nil")
	}
	nsqRegister = def
	dLog = logger
	return nil
}
func (s *dobby) Launch(topic string) error {
	if nsqRegister == nil {
		return errors.New("pls setup NSQ register first")
	}
	for channel, handler := range s.hdls {
		err := nsqRegister(topic, channel, handler)
		if err != nil {
			dLog.Panicf("failed to register NSQ reader hander,err:%s", err.Error())
			return err
		}
	}
	return nil
}

func (s *dobby) registerWorkers(workers ...interface{}) error {
	s.hdls = make(map[string]func(message *nsq.Message) error)
	for _, worker := range workers {
		for _, fm := range getAllMethods(reflect.TypeOf(worker)) {
			if fm.Type.NumIn() < 2 || fm.Type.In(1).Kind() != reflect.Ptr || fm.Type.In(1).Elem().String() != "nsq.Message" || fm.Type.NumOut() != 1 || !fm.Type.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				continue
			}
			funcName := fm.Name
			lstInx := strings.LastIndex(funcName, "Handler")
			if lstInx <= 0 {
				continue
			}
			channelName := snakeString(funcName[0:lstInx])
			if _, ok := s.hdls[channelName]; ok {
				return fmt.Errorf("channel:%s is aready registed to %s, can't reqister it again", channelName, funcName)
			}
			dLog.Println(channelName, "registerd")
			s.hdls[channelName] = s.makeHandler(fm, worker)
		}
	}
	return nil
}

func (s *dobby) makeHandler(fm reflect.Method, worker interface{}) func(message *nsq.Message) error {
	return func(message *nsq.Message) error {
		// msg := new(NsqMessage)
		// err := json.Unmarshal(message.Body, msg)
		// if err != nil { //drop illegal msgs
		// 	dLog.Printf("NSQ message is not legal:" + err.Error())
		// 	return nil
		// }
		ret := fm.Func.Call([]reflect.Value{reflect.ValueOf(worker), reflect.ValueOf(message)})
		if ret[0].Interface() != nil {
			return ret[0].Interface().(error)
		}
		return nil
	}
}

func NewWorker(workers ...interface{}) (ret *dobby, err error) {
	s := new(dobby)
	err = s.registerWorkers(workers...)
	if err != nil {
		return
	}
	ret = s
	return
}
