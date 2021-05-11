package notifications

import "fmt"

type NType string

func (n NType) String() string {
	return string(n)
}

const (
	Success NType = "success"
	Error   NType = "danger"
	Info    NType = "info"
)

type Notification struct {
	Text string `json:"text"`
	Type NType  `json:"type"`
}

type Notifications struct {
	Count int
	List  []Notification `json:"notifications"`
}

func (n *Notifications) Error(text string) {
	n.Count++
	n.List = append(n.List, Notification{
		Text: text,
		Type: Error,
	})
}

func (n *Notifications) Errorf(format string, args ...interface{}) {
	n.Count++
	n.List = append(n.List, Notification{
		Text: fmt.Sprintf(format, args...),
		Type: Error,
	})
}

func (n *Notifications) Info(text string) {
	n.Count++
	n.List = append(n.List, Notification{
		Text: text,
		Type: Info,
	})
}

func (n *Notifications) Infof(format string, args ...interface{}) {
	n.Count++
	n.List = append(n.List, Notification{
		Text: fmt.Sprintf(format, args...),
		Type: Info,
	})
}

func (n *Notifications) Success(text string) {
	n.Count++
	n.List = append(n.List, Notification{
		Text: text,
		Type: Success,
	})
}

func (n *Notifications) Successf(format string, args ...interface{}) {
	n.Count++
	n.List = append(n.List, Notification{
		Text: fmt.Sprintf(format, args...),
		Type: Success,
	})
}

func (n *Notifications) Add(t NType, text string) {
	n.Count++
	n.List = append(n.List, Notification{
		Text: text,
		Type: t,
	})
}

func (n *Notifications) Clear() {
	n.Count = 0
	n.List = nil
}
