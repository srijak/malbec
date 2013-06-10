package main

import (
//	"fmt"
//	"net/mail"
	"database/sql"
//	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
  "os"
  "path"
  "log"
  /*
  "io/ioutil"
  "time"
  */
)

type Contact struct {
  Name string
  Emails []string
}

type ContactService interface {
	Add(name, email string) (err error)
  GetName(email string) (name string)
}

type SqliteContactService struct {
  folder string
  conn *sql.DB
}
func (c *SqliteContactService) getDbConn() (db *sql.DB){
  folder := path.Join(c.folder, "index")
  os.MkdirAll(folder, 0700)
  filename := path.Join(folder, "contacts")

  // email address is what we will be updating occurences of/doing actions on.
  // if an email address has a name, then the name will be inserted(if not present)
  //  in contacts table and the contact_id set in the emails table.
  if c.conn  == nil {
    db, _ = sql.Open("sqlite3", filename)
    _, err := db.Exec("CREATE TABLE IF NOT EXISTS contacts " +
      " (id INTEGER PRIMARY KEY, name varchar(80) UNIQUE" +
      " ); ")
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS emails " +
      " (id INTEGER PRIMARY KEY, contact_id INTEGER, email varchar(256) UNIQUE, " +
      "  occurences INTEGER default 1, sent_invite INTEGER default 0 " +
      " ); ")

    _, err = db.Exec("create index if not exists idx_email_occurs on emails(occurences DESC);")
    _, err = db.Exec("create index if not exists idx_contact_id on emails(contact_id);")
    if err != nil {
      log.Printf("\n\nError creating table: %v \n\n", err)
    }
  }else{
    db = c.conn
  }

  return
}


func (c *SqliteContactService) GetName(email string) (name string) {
  conn := c.getDbConn()
  name = ""
  contact_stmt, _  := conn.Prepare( " select name from contacts where id in (select contact_id from emails where email = ? limit 1) ")
  defer contact_stmt.Close()
  err := contact_stmt.QueryRow(email).Scan(&name)
  switch {
    case err == sql.ErrNoRows:
      log.Printf("No rows returned. We don't have a name for email: [%v]", email)
    case err != nil:
      log.Printf("Error querying for contact with name: [%v]. Error: %v", name, err)
  }

  return
}

func (c *SqliteContactService) Add(name, email string) (err error){

  conn := c.getDbConn()
  contact_id := -1
  if len(name)> 0 {
    // try to get the contact id from 
    contact_create := " insert into contacts (name)  values(?);"
    _, err := conn.Exec(contact_create, name)
    if err != nil {
      log.Printf("Error inserting new contact. Error: %v", err)
    }
    contact_stmt, _  := conn.Prepare( " select id from contacts where name = ? limit 1 ")
    defer contact_stmt.Close()
    err = contact_stmt.QueryRow(name).Scan(&contact_id)
    switch {
      case err == sql.ErrNoRows:
        log.Printf("No rows returned. IMPOSSIBLU")
      case err != nil:
        log.Printf("Error querying for contact with name: [%v]. Error: %v", name, err)
    }
  }


  cmd := "insert or ignore into emails (email, contact_id) values(?, ?); ";

  update :=" update emails set occurences = occurences + 1  where email = ? ;";

  _, err = conn.Exec(cmd, email, contact_id)
  _, err = conn.Exec(update, email)
  if err != nil {
    log.Printf("Error inserting contact. Name:[%v]  Email:[%v].  Error: %v", name, email, err )
  }

  log.Printf("ADDED EMAIL: [%v]", email)
  
  return nil
}

func NewSqliteContactService(folder string) (s *SqliteContactService) {
  // if folder doesn't exist, create it
  // then open up the metadata file in it, or create if it doesn't exist 
  // the metadata folder will store 
  os.MkdirAll(folder, 0700)
  s = &SqliteContactService{folder: folder}
  return
}
