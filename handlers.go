package main

import (
	"errors"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	//set up the collection and query
	coll := ctx.C("entries")
	query := coll.Find(nil).Sort("-timestamp")

	//execute the query
	//TODO: add pagination :)
	var entries []Entry
	if err = query.All(&entries); err != nil {
		return
	}

	//execute the template
	return T("index.html").Execute(w, map[string]interface{}{
		"entries": entries,
		"ctx":     ctx,
	})
}

func sign(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	//we need a user to sign to
	if ctx.User == nil {
		err = errors.New("Can't sign without being logged in")
		return
	}

	entry := NewEntry()
	entry.Name = ctx.User.Username
	entry.Message = req.FormValue("message")

	if entry.Message == "" {
		entry.Message = "Some dummy who forgot a message."
	}

	coll := ctx.C("entries")
	if err = coll.Insert(entry); err != nil {
		return
	}

	//ignore errors: it's ok if the post count is wrong. we can always look at
	//the entries table to fix.
	ctx.C("users").Update(bson.M{"_id": ctx.User.ID}, bson.M{
		"$inc": bson.M{"posts": 1},
	})

	http.Redirect(w, req, reverse("index"), http.StatusSeeOther)
	return
}

func loginForm(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	t := T("login.html").Execute(w, map[string]interface{}{
		"ctx": ctx,
	})
	if err = ctx.Session.Save(req, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return t
}

func login(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	username, password := req.FormValue("username"), req.FormValue("password")

	user, e := Login(ctx, username, password)
	if e != nil {
		ctx.Session.AddFlash("Invalid Username/Password")
		return loginForm(w, req, ctx)
	}

	//store the user id in the values and redirect to index
	ctx.Session.Values["user"] = user.ID
	if err = ctx.Session.Save(req, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, reverse("index"), http.StatusSeeOther)
	return nil
}

func logout(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	delete(ctx.Session.Values, "user")
	if err = ctx.Session.Save(req, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, reverse("index"), http.StatusSeeOther)
	return nil
}

func registerForm(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	t := T("register.html").Execute(w, map[string]interface{}{
		"ctx": ctx,
	})
	if err = ctx.Session.Save(req, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return t
}

func register(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	username, password := req.FormValue("username"), req.FormValue("password")

	u := &User{
		Username: username,
		ID:       bson.NewObjectId(),
	}
	u.SetPassword(password)

	if err = ctx.C("users").Insert(u); err != nil {
		ctx.Session.AddFlash("Problem registering user.")
		return registerForm(w, req, ctx)
	}

	//store the user id in the values and redirect to index
	ctx.Session.Values["user"] = u.ID
	if err = ctx.Session.Save(req, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, reverse("index"), http.StatusSeeOther)
	return nil
}
