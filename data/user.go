package data

import (
	"errors"
	"sync/atomic"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id        int    `bson:"id"`
	Uuid      string `bson:"uuid"`
	Name      string `bson:"name"`
	Email     string `bson:"email"`
	Password  string `bson:"password"`
	CreatedAt string `bson:"createdat"`
}

type Session struct {
	Id        int    `bson:"id"`
	Uuid      string `bson:"uuid"`
	Email     string `bson:"email"`
	UserId    int    `bson:"userid"`
	CreatedAt string `bson:"createdat"`
}

var sid int32 = 0
var userid int32 = 0

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	conn := ss.Copy()
	defer conn.Close()

	id := atomic.AddInt32(&sid, 1)
	uuid := createUUID()
	ct := time.Now().Format("2006-01-02 15:04:05")
	_, err = conn.DB("chat").C("ssessions").Upsert(bson.M{"id": id},
		bson.M{
			"id":        id,
			"uuid":      uuid,
			"email":     user.Email,
			"userid":    user.Id,
			"createdat": ct,
		})

	if err == nil {
		session = Session{
			Id:        int(id),
			Uuid:      uuid,
			Email:     user.Email,
			UserId:    user.Id,
			CreatedAt: ct,
		}
	}

	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("sessions").Find(bson.M{"userid": user.Id}).One(&session)

	return
}

// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	conn := ss.Copy()
	defer conn.Close()

	valid = true
	se := Session{}
	err = conn.DB("chat").C("sessions").Find(bson.M{"uuid": session.Uuid}).One(&se)
	if err != nil {
		valid = false
	}
	if se.Id == 0 {
		valid = false
		err = errors.New("no session found.")
	}

	return
}

// Delete session from database
func (session *Session) DeleteByUUID() (err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("sessions").Remove(bson.M{"uuid": session.Uuid})

	return
}

// Get the user from the session
func (session *Session) User() (user User, err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("users").Find(bson.M{"id": session.UserId}).One(&user)

	return
}

// Delete all sessions from database
func SessionDeleteAll() (err error) {
	conn := ss.Copy()
	defer conn.Close()

	_, err = conn.DB("chat").C("sessions").RemoveAll(bson.M{})

	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {
	conn := ss.Copy()
	defer conn.Close()

	id := atomic.AddInt32(&userid, 1)
	uuid := createUUID()
	ct := time.Now().Format("2006-01-02 15:04:05")
	_, err = conn.DB("chat").C("users").Upsert(bson.M{"id": id},
		bson.M{
			"id":        id,
			"uuid":      uuid,
			"name":      user.Name,
			"email":     user.Email,
			"password":  user.Password,
			"createdat": ct,
		})
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	conn := ss.Copy()
	defer conn.Close()

	_, err = conn.DB("chat").C("users").RemoveAll(bson.M{"id": user.Id})

	return
}

// Update user information in the database
func (user *User) Update() (err error) {
	conn := ss.Copy()
	defer conn.Close()

	_, err = conn.DB("chat").C("users").UpdateAll(bson.M{"id": user.Id},
		bson.M{
			"$set": bson.M{
				"name":  user.Name,
				"email": user.Email,
			},
		})
	return
}

// Delete all users from database
func UserDeleteAll() (err error) {
	conn := ss.Copy()
	defer conn.Close()

	_, err = conn.DB("chat").C("users").RemoveAll(bson.M{})
	return
}

// Get all users in the database and returns it
func Users() (users []User, err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("users").Find(bson.M{}).All(&users)
	return
}

// Get a single user given the email
func UserByEmail(email string) (user User, err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("users").Find(bson.M{"email": email}).One(&user)
	return
}

// Get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("users").Find(bson.M{"email": uuid}).One(&user)
	return
}
