package main

import (
	"bytes"
	"code.google.com/p/go-imap/go1/imap"
	"log"
	"net/mail"
)

type FolderWorker struct {
	Folder          string
	Account         *IMAPAccount
	MetadataService MetadataService
	ep              EmailProcessor
	conn            *IMAPConnection
}

// Things a sync needs:
//  lastseenuuid
//
//  With no more info the client can start fetching the new messages
//  UID FETCH <lastseenuuid+1>:* <descriptors>
//   You need to do this only if the UIDNEXT has changed since the last time.
//

func NewFolderWorker(folder string, acct *IMAPAccount, md MetadataService, ep EmailProcessor) *FolderWorker {
	return &FolderWorker{Folder: folder, Account: acct, MetadataService: md, ep: ep}
}

func (f *FolderWorker) getConnection() (conn *IMAPConnection, err error) {
	conn = f.conn
	if conn == nil {
		log.Printf("Conn was nil. Creating new  imap conection")
		conn, err = NewIMAPConnection(f.Account)
	}
	err = conn.VerifyConnected()
	if err != nil {
		conn = f.conn
	}
	f.conn = conn

	return
}
func (f *FolderWorker) run() (err error) {
	c, err := f.getConnection()
	defer c.Close()
	if err != nil {
		log.Printf("Error getting connection. %v\n", err)
		return
	}

	// do the syncing. that means:
	// ** From: http://www.ietf.org/rfc/rfc4549.txt **
	// 1 check the UIDVALIDITY of the mailbox.
	//     if it doesn't match:
	//        empty everythign in the mailbox
	//        remove any pending action that refer to the UIDs in that mailbox
	//        do a fresh sync (means just go to 2)
	//     else:
	// 2      discover new messages
	//        discover changes to old messages
	// 3 fetch the bodies of any "interesting" messages that the client doesn't already have.
	local_status, first_sync := f.MetadataService.MboxStatus(f.Account, f.Folder)
	remote_status, err := c.Examine(f.Folder)

	if !first_sync && local_status.UIDValidity != remote_status.UIDValidity {
		log.Printf("ALERT: UIDVALIDITY has changed. TODO.")
		return
	}
	old_uidnext := local_status.UIDNext

	f.syncNew(old_uidnext, remote_status.UIDNext, func(uidnext uint32) {
		remote_status.UIDNext = uidnext
		f.MetadataService.SaveMboxStatus(f.Account, remote_status)
	})

	if !first_sync {
		f.syncOld(old_uidnext)
	}

	return
}

func (f *FolderWorker) syncNew(old_uidnext, new_uidnext uint32, fun func(uid uint32)) {
	if old_uidnext == new_uidnext {
		log.Printf("No new messages.")
		fun(old_uidnext)
		return
	}

	log.Printf("Old uid next: %v    New uid next: %v ", old_uidnext, new_uidnext)
	var chunk_size uint32 = 50
	if old_uidnext == 0 {
		old_uidnext = 1
	}

	for i := old_uidnext; i < new_uidnext; i += chunk_size {
		set := &imap.SeqSet{}
		set.AddRange(i, i+chunk_size)
		uidnext, err := f.fetchNewMessages(set)
		if err != nil {
			log.Printf("Error. couldn't fetch uid: [%v] Error: [%v]", i, err)
		} else {
			fun(uidnext)
		}
	}
	return
}

func (f *FolderWorker) fetchNewMessages(uids *imap.SeqSet) (uid_next uint32, err error) {
	c, err := f.getConnection()
	if err != nil {
		return
	}
	log.Printf("Fetching range: %v", uids)
	cmd, err := c.conn.UIDFetch(uids, "RFC822", "FLAGS")
	if err != nil {
		log.Printf("Unable to fetch %v", err)
		return
	}

	for cmd.InProgress() {
		c.conn.Recv(-1)
		for _, rsp := range cmd.Data {
			uid := imap.AsNumber(rsp.MessageInfo().Attrs["UID"])
			uid_next = uint32(uid) + 1
			mime := imap.AsBytes(rsp.MessageInfo().Attrs["RFC822"])
			flags := imap.AsString(rsp.MessageInfo().Attrs["FLAGS"])

			if msg, _ := mail.ReadMessage(bytes.NewReader(mime)); msg != nil {
				f.ep.Add(f.Account, f.Folder, uid,flags, msg)
			}
			cmd.Data = nil
		}

	}
	return
}

func (f *FolderWorker) syncOld(old_uidnext uint32) {
	// discover deleted messages
	// discover changes to old messages
}
