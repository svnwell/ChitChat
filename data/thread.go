package data

import (
	"sync/atomic"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Thread struct {
	Id        int    `bson:"id"`
	Uuid      string `bson:"uuid"`
	Topic     string `bson:"topic"`
	UserId    int    `bson:"userid"`
	CreatedAt string `bson:"createdat"`
}

type Post struct {
	Id        int    `bson:"id"`
	Uuid      string `bson:"uuid"`
	Body      string `bson:"body"`
	UserId    int    `bson:"userid"`
	ThreadId  int    `bson:"threadid"`
	CreatedAt string `bson:"createdat"`
}

var postid int32 = 0
var threadid int32 = 0

// format the CreatedAt date to display nicely on the screen
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt
}

// get the number of posts in a thread
func (thread *Thread) NumReplies() (count int) {
	conn := ss.Copy()
	defer conn.Close()
	coll := conn.DB("chat").C("posts")

	count, err := coll.Find(bson.M{"threadid": thread.Id}).Count()
	if err != nil {
		return 0
	}
	return count
}

// get posts to a thread
func (thread *Thread) Posts() (posts []Post, err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("posts").Find(bson.M{"threadid": thread.Id}).All(&posts)

	return
}

// Create a new thread
func (user *User) CreateThread(topic string) (conv Thread, err error) {
	conn := ss.Copy()
	defer conn.Close()

	id := atomic.AddInt32(&threadid, 1)
	uuid := createUUID()
	_, err = conn.DB("chat").C("threads").Upsert(bson.M{"id": id},
		bson.M{
			"id":        id,
			uuid:        uuid,
			"topic":     topic,
			"userid":    user.Id,
			"createdat": user.CreatedAt,
		})

	if err == nil {
		conv = Thread{
			int(id), uuid, topic, user.Id, user.CreatedAt,
		}
	}
	return
}

// Create a new post to a thread
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	conn := ss.Copy()
	defer conn.Close()

	id := atomic.AddInt32(&postid, 1)
	uuid := createUUID()
	ct := time.Now().Format("2006-01-02 15:04:05")
	_, err = conn.DB("chat").C("posts").Upsert(bson.M{"id": id},
		bson.M{
			"id":        id,
			"uuid":      uuid,
			"userid":    user.Id,
			"threadid":  conv.Id,
			"createdat": ct,
		})
	if err == nil {
		post = Post{
			Id:        int(id),
			Uuid:      uuid,
			UserId:    user.Id,
			ThreadId:  conv.Id,
			CreatedAt: ct,
		}
	}
	return
}

// Get all threads in the database and returns it
func Threads() (threads []Thread, err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("threads").Find(bson.M{}).Sort("-createat").All(&threads)

	return
}

// Get a thread by the UUID
func ThreadByUUID(uuid string) (conv Thread, err error) {
	conn := ss.Copy()
	defer conn.Close()

	err = conn.DB("chat").C("threads").Find(bson.M{"uuid": uuid}).One(&conv)

	return
}

// Get the user who started this thread
func (thread *Thread) User() (user User) {
	conn := ss.Copy()
	defer conn.Close()

	conn.DB("chat").C("users").Find(bson.M{"id": thread.UserId}).One(&user)
	return
}

// Get the user who wrote the post
func (post *Post) User() (user User) {
	conn := ss.Copy()
	defer conn.Close()

	conn.DB("chat").C("users").Find(bson.M{"id": post.UserId}).One(&user)

	return
}
