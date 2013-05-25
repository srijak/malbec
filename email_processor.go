package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
)

type EmailProcessor interface {
	Add(acct *IMAPAccount, mbox_name string, uid uint32, msg *mail.Message) (err error)
}

type PrintingEmailProcessor struct {
	MetadataService MetadataService
}

func (p *PrintingEmailProcessor) Add(account *IMAPAccount, mbox_name string, uid uint32, msg *mail.Message) (err error) {
	var msgdata = map[string]string{}

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
