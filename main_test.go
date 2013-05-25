package main

import (
	. "launchpad.net/gocheck"
	"testing"
)

// hook up gocheck into go test runner
func Test(t *testing.T) { TestingT(t) }

type ImapSuite struct{}

var _ = Suite(&ImapSuite{})

func (s *ImapSuite) TestDefaultGmailServer(c *C) {
	d := GmailServer()
	c.Check(d.Host, Equals, "imap.gmail.com")
	c.Check(d.Port, Equals, uint16(993))
}
