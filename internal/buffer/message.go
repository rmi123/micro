package buffer

import (
	"github.com/zyedidia/micro/v2/internal/config"
	"github.com/zyedidia/tcell"
)

type MsgType int

const (
	MTInfo = iota
	MTWarning
	MTError
)

// Message represents the information for a gutter message
type Message struct {
	// The Msg iteslf
	Msg string
	// Start and End locations for the message
	Start, End Loc
	// The Kind stores the message type
	Kind MsgType
	// The Owner of the message
	Owner string
}

// NewMessage creates a new gutter message
func NewMessage(owner string, msg string, start, end Loc, kind MsgType) *Message {
	return &Message{
		Msg:   msg,
		Start: start,
		End:   end,
		Kind:  kind,
		Owner: owner,
	}
}

// NewMessageAtLine creates a new gutter message at a given line
func NewMessageAtLine(owner string, msg string, line int, kind MsgType) *Message {
	start := Loc{-1, line - 1}
	end := start
	return NewMessage(owner, msg, start, end, kind)
}

func (m *Message) Style() tcell.Style {
	switch m.Kind {
	case MTInfo:
		if style, ok := config.Colorscheme["gutter-info"]; ok {
			return style
		}
	case MTWarning:
		if style, ok := config.Colorscheme["gutter-warning"]; ok {
			return style
		}
	case MTError:
		if style, ok := config.Colorscheme["gutter-error"]; ok {
			return style
		}
	}
	return config.DefStyle
}

func (b *Buffer) AddMessage(m *Message) {
	b.Messages = append(b.Messages, m)
}

func (b *Buffer) removeMsg(i int) {
	copy(b.Messages[i:], b.Messages[i+1:])
	b.Messages[len(b.Messages)-1] = nil
	b.Messages = b.Messages[:len(b.Messages)-1]
}

func (b *Buffer) ClearMessages(owner string) {
	for i := len(b.Messages) - 1; i >= 0; i-- {
		if b.Messages[i].Owner == owner {
			b.removeMsg(i)
		}
	}
}

func (b *Buffer) ClearAllMessages() {
	b.Messages = make([]*Message, 0)
}

//////////////////////////
// !!! PSEUDO CODE !!! ///
//////////////////////////

func (b *Buffer) InitializeOwnerNavigation(owner string) {

	on := OwnerNavigation{
		messages: make([]*Message, 0),
		curMessage: -1
	}

	// Make a view of the messages belonging to the owner.

	for _, m := on.messages {
		if m.Owner == owner {
			on.messages = append(on.messages, m)
		}
	}

	ownerNavigations[owner] = on
}

func (b *Buffer) NavigateToCertainOwnerMessage(owner string, next bool) {

	on := b.ownerNavigations[owner]

	if len(on.messages) == 0 {

		return
	}

	if next { on.curMessage++ } else { on.curMessage-- }

	if on.curMessage < 0 || on.curMessage >= len(on.messages) {

		on.curMessage = 0
	}

	b.GetActiveCursor().GoToLoc(on.messages[on.curMessage].Start)
}

func (b *Buffer) NavigateToPreviousOwnerMessage(owner string) {

	NavigateToCertainOwnerMessage(owner, false)
}

func (b *Buffer) NavigateToNextOwnerMessage(owner string) {

	NavigateToCertainOwnerMessage(owner, true)
}

//////////////////////////
