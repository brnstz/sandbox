package ohmy_test

import (
	"github.com/brnstz/sandbox/ohmy"
	"gopkg.in/mgo.v2"

	"testing"
)

func TestOhMy(t *testing.T) {
	ohmy.Doit()
	s, err := mgo.Dial("192.168.59.103")
	if err != nil {
		panic(err)
	}

	c := s.DB("ohmy").C("bands")
	band := ohmy.Band{
		Name: "The Rolling Stones",
	}
	err = c.Insert(band)
	if err != nil {
		panic(err)
	}
}
