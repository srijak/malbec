package main

type WebsocketCommandFunc func (req *WebsocketCommand) (response interface{})

var (
  validWebsocketRequests = map[string]WebsocketCommandFunc{
    "blah" : s,
  }
)

// command error codes
const (
  _ = iota
  CE_NOT_IMPLEMENTED = iota * -1
)

type CommandError struct {
  ErrorCode int
  Callback_Id int  // any result we send back *must* have the callback id so the client can do matching.
}

func requestHandler(req *WebsocketCommand) (response interface{}){
  return CommandError{ErrorCode: CE_NOT_IMPLEMENTED, Callback_Id: req.Callback_Id}
}

func s(cmd *WebsocketCommand) (resp interface{}){
  return nil
}
