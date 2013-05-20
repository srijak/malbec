package main

import (
  "log"
  "github.com/kuroneko/gosqlite3"
)


type FolderWorker struct {
  Folder string
  Account * IMAPAccount
  MetadataService MetadataService
  conn * IMAPConnection
}

// Things a sync needs:
//  lastseenuuid
//
//  With no more info the client can start fetching the new messages
//  UID FETCH <lastseenuuid+1>:* <descriptors>
//   You need to do this only if the UIDNEXT has changed since the last time.
//

type MetadataService interface {
  MboxStatus(name string) (status *MboxStatus, first_sync bool)
}

type SqliteMetadata struct {
  filename string
}
func (s *SqliteMetadata) NewSqliteMetadata(filename string){
  db, _ := sqlite3.Open(filename) 
  // TODO: store metadate in sqlite.
  // need to create a table to store MboxStatus
}

func NewFolderWorker(folder string, acct *IMAPAccount, md MetadataService ) *FolderWorker {
  return &FolderWorker{Folder: folder, Account: acct, MetadataService:md }
}

func (f *FolderWorker) getConnection () (conn *IMAPConnection, err error){
  if conn == nil {
    f.conn, err = NewIMAPConnection(f.Account)
  }
  err = f.conn.VerifyConnected()

  return
}
func (f *FolderWorker) run() (err error) {
  c, err := f.getConnection()
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
  local_status, first_sync := f.MetadataService.MboxStatus(f.Folder)
  remote_status, err := c.Examine(f.Folder)
  log.Printf("Local Status: %v \nRemote Status: %v \n", local_status, remote_status)

  if !first_sync && local_status.UIDValidity != remote_status.UIDValidity {
    log.Printf("ALERT: UIDVALIDITY has changed. TODO.")
    return
  }
  old_uidnext := local_status.UIDNext

  f.syncNew(old_uidnext, remote_status.UIDNext)
  if !first_sync {
    f.syncOld(old_uidnext)
  }
  return
}

func (f* FolderWorker) syncNew(old_uidnext, new_uidnext uint32){
  // should sync newest first
  // chunk fetches.
  //conn, err := f.getConnection()
  for i := old_uidnext; i < new_uidnext; i += 50 {
    log.Printf("FETCH %v:%v", i, i+50)
  }
}

func (f* FolderWorker) syncOld(old_uidnext uint32) {
  // discover deleted messages
  // discover changes to old messages
}
