package actor

import (
	"testing"

	"bytes"
	"fmt"
	"reflect"
)

type SayHello struct {
	Actor
}

func (s *SayHello) Recv() chan interface{} {
	s.receive = make(chan interface{})
	return s.receive
}

func (s *SayHello) Do(i int) {
	msg := <-s.receive
	fmt.Println(i + msg.(int))
	close(s.receive)
}

func Spawn(actor Actorable, startArgs []interface{}) chan interface{} {
	act := reflect.ValueOf(actor).MethodByName("Do")
	var buf bytes.Buffer
	for _, v := range startArgs {
		t := reflect.TypeOf(v)
		buf.WriteString(t.Name())
	}
	var buf2 bytes.Buffer
	for i := 0; i < act.Type().NumIn(); i++ {
		t := act.Type().In(i)
		buf2.WriteString(t.Name())
	}
	expected := buf2.String()
	input := buf.String()
	if expected != input {
		panic(fmt.Sprintf("expected: %s, but receive: %s", expected, input))
	}
	if act.Type().NumOut() > 0 {
		panic("expected no return!!!")
	}

	inputs := make([]reflect.Value, len(startArgs))
	for k, in := range startArgs {
		inputs[k] = reflect.ValueOf(in)
	}
	go act.Call(inputs)
	return actor.Recv()
}

func TestSpawn(t *testing.T) {
	sayHi := &SayHello{}
	pid := Spawn(sayHi, []interface{}{3})
	pid <- 30
}

func TestReflect(t *testing.T) {
	sayHi := &SayHello{}
	method := reflect.ValueOf(sayHi).MethodByName("Do")

	fmt.Printf("%v\n", method)
}