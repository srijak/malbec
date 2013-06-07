package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"database/sql"
//	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
  "os"
  "path"
  "io/ioutil"
)

type EmailProcessor interface {
	Add(acct *IMAPAccount, mbox_name string, uid uint32, flags Flags, msg *mail.Message) (err error)
}

type SqliteEmailProcessor struct {
  folder string
  conns map[string]*sql.DB
}

func NewSqliteEmailProcessor(folder string) (s *SqliteEmailProcessor) {
  // if folder doesn't exist, create it
  // then open up the metadata file in it, or create if it doesn't exist 
  // the metadata folder will store 
  os.MkdirAll(folder, 0700)
  s = &SqliteEmailProcessor{folder: folder}
  return
}

func (s *SqliteEmailProcessor) getIndexFor(mbox_name string) (db *sql.DB){
  folder := path.Join(s.folder, "index")
  os.MkdirAll(folder, 0700)
  filename := path.Join(folder, "index")

  if s.conns  == nil {
    s.conns = make(map[string]*sql.DB, 10)
  }
  db, present := s.conns[filename]
	if !present {
    db, _ = sql.Open("sqlite3", filename)
    _, err := db.Exec("CREATE VIRTUAL TABLE IF NOT EXISTS uids USING fts4" +
      " (pk INTEGER PRIMARY KEY, uid INTEGER, deleted INTEGER, flags TEXT, " +
      "  subject TEXT, from TEXT, cc TEXT, to TEXT, account varchar(256), mbox TEXT" +
      " ); ")
    _, err = db.Exec("create index idx_uid_key on uids(uid)")
    if err != nil {
      log.Printf("\n\nError creating table: %v \n\n", err)
    }
  }
  if db != nil {
    s.conns[filename] = db
  }

  return
}

func (s *SqliteEmailProcessor) addToIndex(acct *IMAPAccount, mbox_name , subject, from, cc, to string, flags Flags, uid uint32) (err error){
  idx := s.getIndexFor(mbox_name)

  log.Printf("\n\tInserting uid: %v  flags: %v", uid, flags, subject, from, cc, to, mbox_name, acct.Username)
	cmd := "insert into uids(uid, deleted, flags, subject, from, cc, to, mbox, account) " +
		" VALUES (?,?,?,?,?,?,?,?,?) "
    flag_json, _ := json.Marshal(flags)
	_, err = idx.Exec(cmd, uid, 0, string(flag_json))
  if err != nil {
    log.Printf("Error inserting uid %v for mailbox [%v] ", uid, mbox_name)
  }
  return
}

func (s *SqliteEmailProcessor) Add(acct *IMAPAccount, mbox_name string, uid uint32,flags Flags, msg *mail.Message) (err error){
  folder := path.Join(s.folder, "emails", mbox_name)
  os.MkdirAll(folder, 0700)

  file := path.Join(folder, fmt.Sprintf("%d", uid))

  log.Printf("Writing %v to file: %v", uid, file)



	var msgdata = map[string]string{}

	for headerkey := range msg.Header {
	  val := msg.Header.Get(headerkey)
	  msgdata[headerkey] = val
	}

	msgdata["imap_uid"] = fmt.Sprintf("%d", uid)
	   if b, err := TextBody(msg); err == nil {
	     msgdata["text_body"] = b
	   }
	   if b, err := HTMLBody(msg); err == nil {
	     msgdata["html_body"] = b
	   }
	o, err := json.Marshal(msgdata)
	if err != nil {
		log.Println("error marshaling message as JSON: ", err.Error()[:100])
    ioutil.WriteFile(file+".ERR", []byte(fmt.Sprintf("Error decoding message body: %v\nout: %v", err.Error(), o)), 0700)
	} else {
    ioutil.WriteFile(file, o, 0700)
	}

  //func (s *SqliteEmailProcessor) addToIndex(acct *IMAPAccount, mbox_name string, subject, from, cc, to, flags Flags, uid uint32) (err error){
  subject := "hi"
  from := "hello"
  cc := "cced to"
  to := "to"
  s.addToIndex(acct, mbox_name, subject, from, cc, to, flags, uid)

  return nil
}

type PrintingEmailProcessor struct {
	MetadataService MetadataService
}

func (p *PrintingEmailProcessor) Add(account *IMAPAccount, mbox_name string, uid uint32, flags string, msg *mail.Message) (err error) {
	var msgdata = map[string]string{}

  log.Printf("FLAGS: %v", flags)
	/*for headerkey := range msg.Header {
	  val := msg.Header.Get(headerkey)
	  msgdata[headerkey] = val
	}*/

	msgdata["imap_uid"] = fmt.Sprintf("%d", uid)
	/*
	   if b, err := TextBody(msg); err == nil {
	     msgdata["text_body"] = b
	   }
	   if b, err := HTMLBody(msg); err == nil {
	     msgdata["html_body"] = b
	   }
	*/
	o, err := json.Marshal(msgdata)
	if err != nil {
		log.Println("error marshaling message as JSON: ", err.Error()[:100])
	} else {
		fmt.Println(string(o))
	}
	return nil
}
