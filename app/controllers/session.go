package controllers

import (
	web "babou/lib/web"
)

// Implements babou/app.Controller interface.
// Maps an action to results or returns 404 otherwise.

type SessionController struct {
	actionMap map[string]web.Action
}

// Registers actions for the HomeController and returns it.
// Note: State in the returned controller object is global to
// all requests the controller processes.
func NewSessionController() *SessionController {
	sc := &SessionController{}
	sc.actionMap = make(map[string]web.Action)

	//add your actions here.

	return sc
}

func (sc *SessionController) HandleRequest(action string,
	params map[string]string) *web.Result {

	if sc.actionMap[action] != nil {
		return sc.actionMap[action](params)
	} else {
		return &web.Result{Status: 404, Body: []byte("")}
	}
}
