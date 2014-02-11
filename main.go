package main

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
)

func reverse(name string, things ...interface{}) string {
	//convert the things to strings
	strs := make([]string, len(things))
	for i, th := range things {
		strs[i] = fmt.Sprint(th)
	}
	//grab the route
	u, err := router.GetRoute(name).URL(strs...)
	if err != nil {
		panic(err)
	}
	return u.Path
}

func init() {
	gob.Register(bson.ObjectId(""))
}

var store sessions.Store
var session *mgo.Session
var database string
var router *pat.Router

func main() {
	var err error
	session, err = mgo.Dial(os.Getenv("MONGOLAB_URI"))
	if err != nil {
		panic(err)
	}
	database = session.DB("").Name

	//create an index for the username field on the users collection
	if err := session.DB("").C("users").EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	}); err != nil {
		panic(err)
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("KEY")))

	router = pat.New()
	router.Add("GET", "/login", handler(loginForm)).Name("login")
	router.Add("POST", "/login", handler(login))

	router.Add("GET", "/register", handler(registerForm)).Name("register")
	router.Add("POST", "/register", handler(register))

	router.Add("GET", "/logout", handler(logout)).Name("logout")

	router.Add("GET", "/", handler(hello)).Name("index")

	router.Add("POST", "/sign", handler(sign)).Name("sign")

	fmt.Println("About lo listen on", os.Getenv("PORT"))
	if err = http.ListenAndServe(":"+os.Getenv("PORT"), router); err != nil {
		panic(err)
	}
}
