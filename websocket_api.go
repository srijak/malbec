package main

import (
  l4g "code.google.com/p/log4go"
)

type WebsocketCommandFunc func (req *WebsocketCommand) (res CommandResult)

var (
  availableCommands = map[string]WebsocketCommandFunc{
    "get_account_mailbox_map" : getAccountMailboxMap,
    "get_emails_list" : getEmailsList,
  }

  metadataService = NewSqliteMetadata("metadata.sqlite")
)

// command error codes
const (
  _ = iota
  CE_NOT_FOUND = iota * -1
)


type CommandResult struct {
  Callback_Id int
  Response interface{}
}

func requestHandler(req *WebsocketCommand) (res CommandResult){
  if cmd, ok := availableCommands[req.Type]; ok {
    return cmd(req)
  }
  return CommandResult{Response: map[string]int{
    "errorCode" : CE_NOT_FOUND,
    },
    Callback_Id: req.Callback_Id}
}


func getAccountMailboxMap(cmd *WebsocketCommand) (CommandResult){
  // look up the info from the db, always.
  // Except something else keep the db upto date.
  // as the MboxStatus row always has a email column,
  // we can just get all lines from it and then 
  res := metadataService.AccountsAndMailboxes()
  l4g.Info("\n\n\nReturning %#v\n\n", res)
  return CommandResult{Response: res,
                      Callback_Id: cmd.Callback_Id}
}
