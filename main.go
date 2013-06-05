package main

import (
	"code.google.com/p/go-imap/go1/imap"
	"log"
	"os"
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
  "code.google.com/p/go.net/websocket"
  l4g "code.google.com/p/log4go"
)

func VersionHandler(w http.ResponseWriter, req *http.Request){
  fmt.Fprintf(w, "HI")
}

func main() {
  r := mux.NewRouter()
  l4g.Info("Starting.")

  r.Handle("/ws", websocket.Handler(wsHandler))
  r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/")))
  http.Handle("/", r)
  l4g.Info("Serving...")
  err := http.ListenAndServe(":8081", nil)

  if err != nil {
    l4g.Info("Error serving: %v", err)
  }

}

func haha(){
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
			l4g.Info("Synced mailbox: %v", m.Name)
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
	      l4g.Info("status: %v\n", status)
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
	    l4g.Info("%v", uid)
	  }
	*/
}
