package controllers

import (
	"github.com/drbawb/babou/app/filters"
	"github.com/drbawb/babou/app/models"

	"errors"
	"fmt"
	"github.com/drbawb/babou/lib/web"
	"strconv"
)

var renderer web.Renderer = web.NewMustacheRenderer("app/admin/views")

type UsersController struct {
	*App
	Auth *filters.AuthContext
}

func (au *UsersController) Dispatch(action, accept string) (web.Controller, web.Action) {
	newAu := &UsersController{}
	newAu.App = &App{}

	switch action {
	case "index":
		return newAu, newAu.Index
	case "delete":
		return newAu, newAu.Delete
	}

	panic("unreachable")
}

func (au *UsersController) Index() *web.Result {
	res := &web.Result{Status: 200}

	usersList, err := models.AllUsers()
	if err != nil {
		//todo flash error
		res.Body = []byte(err.Error())
		return res
	}

	context := &struct {
		Users []*models.User
	}{
		Users: usersList,
	}

	///res.Body = []byte(fmt.Sprintf("len context.Users: %d", len(context.Users)))
	res.Body = []byte(au.Out.RenderWith("bootstrap", "user", "index", context))
	return res
}

func (au *UsersController) Delete() *web.Result {
	res := &web.Result{Status: 200}

	userToDestroy := &models.User{}

	userId, err := strconv.Atoi(au.Dev.Params.All["id"])
	if err != nil {
		res.Body = []byte(err.Error())
		return res
	}

	err = userToDestroy.SelectId(userId)
	if err != nil {
		res.Body = []byte(err.Error())
		return res
	}

	err = userToDestroy.Delete()
	if err != nil {
		res.Body = []byte(err.Error())
		return res
	}

	res.Body = []byte(fmt.Sprintf(
		"user [%s] has been judged.",
		userToDestroy.Username))

	return res
}

func (uc *UsersController) SetAuthContext(context *filters.AuthContext) error {
	if context == nil {
		return errors.New("No AuthContext was supplied to this controller!")
	}

	uc.Auth = context
	uc.Auth.Required = false

	return nil
}
