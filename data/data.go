package data

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

var ss *mgo.Session

func init() {
	var err error
	ss, err = mgo.Dial("localhost")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func": "init",
			"err":  err,
		}).Fatal("failed connetting to mongo.")
	}
	return
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func": "createUUID",
			"err":  err,
		}).Fatal("uuid gen failed.")
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
