package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type MetadataService interface {
	MboxStatus(account *IMAPAccount, name string) (status *MboxStatus, first_sync bool)
	SaveMboxStatus(account *IMAPAccount, mbox *MboxStatus) (err error)
}

type SqliteMetadata struct {
	filename string
	db       *sql.DB
}

func (s *SqliteMetadata) MboxStatus(acct *IMAPAccount, name string) (status *MboxStatus, first_sync bool) {
	status = NewMboxStatus()
	first_sync = true

	cmd := "select pk, name, flags, perm_flags, messages, recent, unseen, uid_next, uid_validity " +
		"from mailbox_status where name=? and email=? "
	stmt, _ := s.db.Prepare(cmd)
	defer stmt.Close()
	var id int64
	var flags []byte
	var perm_flags []byte
	err := stmt.QueryRow(name, acct.Username).Scan(&id, &(status.Name), &flags, &perm_flags, &(status.Messages), &(status.Recent), &(status.Unseen), &(status.UIDNext), &(status.UIDValidity))
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No rows returned.")
		return status, true
	case err != nil:
		log.Printf("\n\nError querying db: %v\n\n", err)
	}

	err = json.Unmarshal(flags, &(status.Flags))
	if err != nil {
		log.Printf("Error unmarshalling: %v to Flags", string(flags))
		err = nil
	}
	err = json.Unmarshal(perm_flags, &(status.PermFlags))
	if err != nil {
		log.Printf("Error unmarshalling: %v to Flags", string(perm_flags))
		err = nil
	}

	log.Printf("loaded mbox: %v", status)
	return
}
func (s *SqliteMetadata) SaveMboxStatus(acct *IMAPAccount, mbox *MboxStatus) (err error) {
	flags, err := json.Marshal(mbox.Flags)
	if err != nil {
		log.Println("error marshaling message as JSON: ", err.Error()[:100])
		flags = []byte("{'error': 'Error converting to json.'}")
		err = nil
	}
	perm_flags, err := json.Marshal(mbox.PermFlags)
	if err != nil {
		log.Println("error marshaling message as JSON: ", err.Error()[:100])
		flags = []byte("{'error': 'Error converting to json.'}")
		err = nil
	}

	cmd := "insert or replace into mailbox_status(name, flags, perm_flags, messages, recent, unseen, uid_next, uid_validity, email) " +
		" VALUES (?,?,?,?,?,?,?,?,?) "
	_, err = s.db.Exec(cmd, mbox.Name, string(flags), string(perm_flags), mbox.Messages, mbox.Recent, mbox.Unseen, mbox.UIDNext, mbox.UIDValidity, acct.Username)
	if err != nil {
		log.Printf("\n\nError saving mailbox status %v : %v", mbox, err)
	}

	return
}

func NewSqliteMetadata(filename string) *SqliteMetadata {
	db, _ := sql.Open("sqlite3", filename)
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS mailbox_status " +
		"(pk INTEGER PRIMARY KEY, name VARCHAR(100), flags VARCHAR(1024), perm_flags VARCHAR(1024), " +
		"  messages INTEGER, recent INTEGER, unseen INTEGER, uid_next INTEGER, uid_validity INTEGER, " +
		"  email VARCHAR(255)) ")
	_, err = db.Exec("create unique index idx_composite_key on mailbox_status(email, name)")
	if err != nil {
		log.Printf("\n\nError creating table: %v \n\n", err)
	}

	/*
	   int res = "CREATE TABLE IF NOT EXISTS email "
	     "(pk INTEGER PRIMARY KEY, datetime REAL, sender_name VARCHAR(50), sender_address VARCHAR(50), "
	     "tos TEXT, ccs TEXT, bccs TEXT, attachments TEXT, msg_id VARCHAR(50), uid VARCHAR(20), folder VARCHAR(20), folder_num INTEGER, folder_num_1 INTEGER, folder_num_2 INTEGER, folder_num_3 INTEGER, extra INTEGER);"] UTF8String] , NULL, NULL, &errorMsg);
	*/

	// TODO: store metadate in sqlite.
	// need to create a table to store MboxStatus
	return &SqliteMetadata{filename: filename, db: db}
}
