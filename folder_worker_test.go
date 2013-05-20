package main

import (
  . "launchpad.net/gocheck"
)

type FolderWorkerSuite struct{}
var _ = Suite(&FolderWorkerSuite{})

func (s *FolderWorkerSuite) TestFolderWorker(c *C){
  // use this to run a FolderWorker
  // so, need to set up a MetadataService
  // to return the test accounts.
  // then make the syncs work.
}
