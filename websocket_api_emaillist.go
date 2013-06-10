package main

import (
  l4g "code.google.com/p/log4go"
  "strconv"
  "fmt"
)

var (
  ep = &SqliteEmailProcessor{folder: "storage"}
)

// EmailList sort types:
const (
  EL_BY_DATE = iota
  EL_BY_ACTION
)

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
  sort := EL_BY_DATE

  if _, ok := cmd.Params["target"]; ok {
    target = cmd.Params["target"]
  }
  if _, ok := cmd.Params["start"]; ok {
    start, _ = strconv.Atoi(cmd.Params["start"])
  }
  if _, ok := cmd.Params["sort"]; ok {
    sort, _ = strconv.Atoi(cmd.Params["sort"])
  }
  fmt.Sprintf("target: %v  start: %v  sort: %v", target, start, sort)
  res := getEmailsListUnified(start, sort)


  l4g.Info("\n\n\nReturning %#v\n\n", res)
  return CommandResult{Response: res,
                      Callback_Id: cmd.Callback_Id}
}

type SparseEmail struct {
}
type SparseEmailList []SparseEmail

func getEmailsListUnified(start, sort int) *SparseEmailList {
  // hmm. yeah, maybe I want to always shard the
  //  uids/emails db by dates. ie, store 1000 uids per db
  //  that makes it easy to stream results back in sorted
  //   date order, which is what people want anyway.
 
  return ep.SparseEmailListUnified(start, sort)
  
}
