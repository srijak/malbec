package main

import (
  l4g "code.google.com/p/log4go"
  "strconv"
  "fmt"
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

func getEmailsList(cmd *WebsocketCommand) (CommandResult){
  l4g.Info("Params: ", cmd.Params)
  // parse params then call the corresponding function.
  //  start of params spec:
  //    target: "unified" vs "email@address.com"
  //    start: start offset
  //    search: "" vs "from: blah@a.com" etcs
  // Might not be the best way, but lets just do this for now.

  // handle defaults
  target := "unified"
  start := 0
  if _, ok := cmd.Params["target"]; ok {
    target = cmd.Params["target"]
  }
  if _, ok := cmd.Params["start"]; ok {
    start, _ = strconv.Atoi(cmd.Params["start"])
  }

  res := fmt.Sprintf("target: %v  start: %v", target, start) 

  l4g.Info("\n\n\nReturning %#v\n\n", res)
  return CommandResult{Response: res,
                      Callback_Id: cmd.Callback_Id}
}

type EmailList uint
func getEmailsListUnified(start uint) uint {
  var page_size uint
  page_size = 100
  // hmm. yeah, maybe I want to always shard the
  //  uids/emails db by dates. ie, store 1000 uids per db
  //  that makes it easy to stream results back in sorted
  //   date order, which is what people want anyway.
  
  return page_size
  
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
