package main

import (
	"code.google.com/p/go-imap/go1/imap"
	"fmt"
	"log"
	"time"
)

type IMAPServer struct {
	Host string
	Port uint16
}

type IMAPAccount struct {
	Username string
	Password string
	Server   *IMAPServer
}

type IMAPConnection struct {
	conn    *imap.Client
	Account *IMAPAccount
}

//  def __init__(self, server="imap.gmail.com", username="vfct3st@gmail.com", passwd="GWxE6kBs436wa7tyedyU"):
func GmailServer() (server *IMAPServer) {
	return &IMAPServer{Host: "imap.gmail.com", Port: 993}
}
func DovecotServer() (server *IMAPServer) {
	return &IMAPServer{Host: "dovecot", Port: 143}
}


func NewIMAPConnection(acct *IMAPAccount) (ic *IMAPConnection, err error) {
	conn, err := connect(acct.Server)
	if err != nil {
		goto RetError
	}

	ic = &IMAPConnection{conn: conn}
	err = ic.login(acct)
	if err != nil {
		goto RetError
	}
	return ic, nil

RetError:
	return nil, err
}

func (ic *IMAPConnection) Close() {
	if ic.conn != nil {
		ic.conn.Logout(time.Duration(2))
	}
}

func connect(server *IMAPServer) (c *imap.Client, err error) {
	addr := fmt.Sprintf("%s:%d", server.Host, server.Port)
	c, err = imap.DialTLS(addr, nil)
	return
}

func (ic *IMAPConnection) login(acct *IMAPAccount) (err error) {
	_, err = ic.conn.Login(acct.Username, acct.Password)
	if err != nil {
		return err
	}
	ic.Account = acct
	return nil
}

func (ic *IMAPConnection) Examine(mbox string) (status *MboxStatus, err error) {
	log.Printf("Selecting mbox %v", mbox)
	_, err = ic.conn.Select(mbox, true)
	status = NewFromMailboxStatus(ic.conn.Mailbox)

	return status, err
}

func (ic *IMAPConnection) FetchAllUids(mbox string, chunk_chan chan uint32) (err error) {
	tm := NewTimer()
	tm.Add("Select")

	ic.conn.Select(mbox, true)
	tm.Add("Select Done")

	set, _ := imap.NewSeqSet("1:*")

	tm.Add("FetchAll")
	cmd, err := ic.conn.UIDFetch(set, "")
	for cmd.InProgress() {
		ic.conn.Recv(-1)
		for _, rsp := range cmd.Data {
			uid := imap.AsNumber(rsp.MessageInfo().Attrs["UID"])
			chunk_chan <- uid
		}
		cmd.Data = nil
	}

	tm.Add("FetchAll Done")
	tm.Report()
	return nil
}

func (ic *IMAPConnection) FetchUidsMostRecent(mbox string) (uids []uint32, err error) {
	ic.conn.Select(mbox, true)
	uid_next := ic.conn.Mailbox.UIDNext

	recent_size := uint32(50)
	uids = make([]uint32, recent_size)

	var set *imap.SeqSet
	if uid_next <= recent_size {
		set, _ = imap.NewSeqSet("1:*")
	} else {
		set, _ = imap.NewSeqSet(fmt.Sprintf("%v:%v", uid_next-recent_size, uid_next))
	}
	cmd, err := ic.conn.UIDFetch(set, "RFC822.SIZE")
	total := 0
	for cmd.InProgress() {
		ic.conn.Recv(-1)
		for _, rsp := range cmd.Data {
			uid := imap.AsNumber(rsp.MessageInfo().Attrs["UID"])
			uids[total] = uid
		}
		cmd.Data = nil
		total++
	}

	uids = uids[:total-1]
	return
}

func fetch(c *IMAPConnection, mbox string, uid_chan chan uint32) {
	// set up where to store the emails etc
	// expects the uid chan to send uids in the mbox mailbox.

}

func (ic *IMAPConnection) VerifyConnected() (err error) {
	if ic.conn.State() == imap.Closed {
		ic.conn, err = connect(ic.Account.Server)
	}
	if err != nil {
		return
	}
	if ic.conn.State() == imap.Logout {
		ic.login(ic.Account)
	}
	return
}

func (ic *IMAPConnection) Mailboxes() (mboxes []*MboxInfo, err error) {
	cmd, err := imap.Wait(ic.conn.List("", "*"))
	mboxes = make([]*MboxInfo, len(cmd.Data))

	for i, d := range cmd.Data {
		mboxes[i] = NewMboxInfoFromMailboxInfo(d.MailboxInfo())
	}
	return mboxes, nil
}
