package main

import (
  . "launchpad.net/gocheck"
//  "testing"
)

// hook up gocheck into go test runner

type MetadataSuite struct{}
var _ = Suite(&MetadataSuite{})

func (s *MetadataSuite) TestAccoutData_Serialization(c *C){
  m := NewMboxStatus()
  a := &AccountData{Name: "test"}
  a.SetMbox(*m)
  a.Save("account")

  b, err := LoadAccountData("account")
  c.Check(err, Equals, nil)
  c.Check(b.Equals(a), Equals, true)
}
 

