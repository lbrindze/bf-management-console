package main

import "io/ioutil"

const (
	dataDir = "data/"
)

var client TextClient = &StubClient{}

func setClient(c TextClient) {
	client = c
}

type BoyfriendActor interface {
	GetName() string
	GetText() string
	NotifyBoyfriend() error
	Save() error
}

type GenericAction struct {
	Name        string
	TextMessage []byte
	Client      TextClient
}

func NewGenericAction(name string, textMsg string) *GenericAction {
	return &GenericAction{Name: name, TextMessage: []byte(textMsg), Client: client}
}

func (a *GenericAction) Save() error {
	filename := dataDir + a.Name + ".txt"
	return ioutil.WriteFile(filename, a.TextMessage, 0600)
}

func loadAction(name string) (*GenericAction, error) {
	filename := dataDir + name + ".txt"
	textMsg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewGenericAction(name, string(textMsg)), nil
}

func (a *GenericAction) GetName() string {
	return a.Name
}
func (a *GenericAction) GetText() string {
	return string(a.TextMessage)
}

func (a *GenericAction) NotifyBoyfriend() error {
	return a.Client.SendMessage(string(a.TextMessage))
}
