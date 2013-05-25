package main

import (
	"code.google.com/p/go-imap/go1/imap"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

type MboxSet map[string]MboxStatus

type MboxInfo struct {
	Attrs Flags
	Delim string
	Name  string
}

func NewMboxInfoFromMailboxInfo(m *imap.MailboxInfo) (ret *MboxInfo) {
	ret = &MboxInfo{Delim: m.Delim, Name: m.Name}
	ret.Attrs = make(Flags, len(m.Attrs))
	for k, v := range m.Attrs {
		ret.Attrs[k] = v
	}
	return
}

type AccountData struct {
	Name   string
	Mboxes MboxSet
}

func (a *AccountData) Equals(o *AccountData) bool {
	if a.Name != o.Name {
		return false
	}
	for m, mbox := range a.Mboxes {
		other_m, present := o.Mboxes[m]
		if !present {
			return false
		}

		if !other_m.Equals(&mbox) {
			return false
		}
	}
	return true
}

func (a *AccountData) SetMbox(mbox MboxStatus) {
	if a.Mboxes == nil {
		a.Mboxes = make(MboxSet, 10)
	}
	a.Mboxes[mbox.Name] = mbox
}

func (a *AccountData) Save(to_file string) (err error) {
	// should lock the file
	fh, err := os.OpenFile(to_file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fh.Close()

	gob.NewEncoder(fh).Encode(a)
	return
}

func LoadAccountData(from_file string) (a *AccountData, err error) {
	fh, err := os.OpenFile(from_file, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	decoder := gob.NewDecoder(fh)
	err = decoder.Decode(&a)
	return a, err
}

//
//S: * FLAGS (\Answered \Flagged \Draft \Deleted \Seen $NotJunk $Junk)
//S: * OK [PERMANENTFLAGS ()] Flags permitted.
//S: * OK [UIDVALIDITY 1] UIDs valid.
//S: * 21 EXISTS
//S: * 0 RECENT
//S: * OK [UIDNEXT 22] Predicted next UID.
//
type Flags map[string]bool
type Uids map[uint]bool

type MboxStatus struct {
	Name        string
	UIDValidity uint32
	UIDNext     uint32
	Messages    uint32
	Recent      uint32
	Unseen      uint32
	Flags       Flags
	PermFlags   Flags
	Uids        Uids
}

func (f Flags) Equals(o Flags) bool {
	if len(f) != len(o) {
		return false
	}
	for k, _ := range f {
		_, present := (o)[k]
		if !present {
			return false
		}
	}
	return true
}
func (f Uids) Equals(o Uids) bool {
	if len(f) != len(o) {
		return false
	}
	for k, _ := range f {
		_, present := (o)[k]
		if !present {
			return false
		}
	}
	return true
}

func (m *MboxStatus) Equals(o *MboxStatus) bool {
	if m.Name != o.Name {
		return false
	}
	if m.UIDValidity != o.UIDValidity {
		return false
	}
	if m.UIDNext != o.UIDNext {
		return false
	}
	if m.Messages != o.Messages {
		return false
	}
	if m.Recent != o.Recent {
		return false
	}
	if m.Unseen != o.Unseen {
		return false
	}
	if !m.Flags.Equals(o.Flags) {
		log.Printf("Flags mismatch")
		return false
	}
	if !m.PermFlags.Equals(o.PermFlags) {
		return false
	}
	if !m.Uids.Equals(o.Uids) {
		return false
	}

	return true
}

func (m *MboxStatus) String() string {
	return fmt.Sprintf("--- %+q ---\n"+
		"Flags:        %v\n"+
		"PermFlags:    %v\n"+
		"Messages:     %v\n"+
		"Recent:       %v\n"+
		"Unseen:       %v\n"+
		"UIDNext:      %v\n"+
		"UIDValidity:  %v\n",
		m.Name, m.Flags, m.PermFlags, m.Messages, m.Recent,
		m.Unseen, m.UIDNext, m.UIDValidity)
}

func NewFromMailboxStatus(ms *imap.MailboxStatus) (m *MboxStatus) {
	m = &MboxStatus{Name: ms.Name,
		UIDValidity: ms.UIDValidity,
		UIDNext:     ms.UIDNext,
		Messages:    ms.Messages,
		Recent:      ms.Recent,
	}
	m.Flags = make(Flags, len(ms.Flags))
	for k, v := range ms.Flags {
		m.Flags[k] = v
	}
	m.PermFlags = make(Flags, len(ms.PermFlags))
	for k, v := range ms.PermFlags {
		m.PermFlags[k] = v
	}
	m.Uids = make(map[uint]bool, 100)
	return m
}

func NewMboxStatus() (m *MboxStatus) {
	m = &MboxStatus{}

	m.Flags = make(Flags, 5)
	m.PermFlags = make(Flags, 5)

	m.Uids = make(map[uint]bool, 100)
	return m
}
