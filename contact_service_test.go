package main

import (
	. "launchpad.net/gocheck"
)

type ContactServiceSuite struct{}

var _ = Suite(&ContactServiceSuite{})


func (s *ContactServiceSuite) TestContactService_Add(c *C) {

  cs := NewSqliteContactService("_test")
  cs.Add("Srijak Rijal", "srijak@gmail.com")
  name := cs.GetName("srijak@gmail.com")

  c.Assert(name, Equals, "Srijak Rijal")

}

