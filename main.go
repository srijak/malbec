package main

import (
	"code.google.com/p/go-imap/go1/imap"
	"log"
	"os"
)

func main() {
	imap.DefaultLogger = log.New(os.Stdout, "", 0)
	//imap.DefaultLogMask = imap.LogConn | imap.LogRaw
	imap.DefaultLogMask = imap.LogConn

	acct := &IMAPAccount{Username: "vfct3st@gmail.com", Password: "GWxE6kBs436wa7tyedyU", Server: GmailServer()}
	ms := NewSqliteMetadata("metadata.sqlite")
  ep := &SqliteEmailProcessor{folder: "storage"}

	ic, err := NewIMAPConnection(acct)
	if err != nil {
		panic(err)
	}

	mboxes, err := ic.Mailboxes()

	for _, m := range mboxes {
		_, present := m.Attrs["\\Noselect"]
		if !present {
			fw := NewFolderWorker(m.Name, acct, ms, ep)
			fw.run()
			log.Printf("Synced mailbox: %v", m.Name)
		}
	}

	/* acct := &IMAPAccount{ Username: "vfct3st@gmail.com", Password: "GWxE6kBs436wa7tyedyU", Server: GmailServer()}
	  ic, err := NewIMAPConnection(acct)
	  if err != nil {
	    panic(err)
	  }

	  a := &AccountData{Name: "test"}
	  mboxes, err := ic.Mailboxes()

	  for _, m := range mboxes {
	    _, present := m.Attrs["\\Noselect"]
	    if ! present {
	      status, _:= ic.Examine(m.Name)
	      log.Printf("status: %v\n", status)
	      a.SetMbox(*status)
	    }
	  }
	  a.Save("account")
	  //ic.FetchUidsMostRecent("INBOX")
	  cc := make(chan uint32)
	  go func(){
	    ic.FetchAllUids("INBOX", cc)
	  }()

	//  fetcher, err := NewIMAPConnection(acct)
	  for {
	    uid := <- cc
	    log.Printf("%v", uid)
	  }
	*/
}
