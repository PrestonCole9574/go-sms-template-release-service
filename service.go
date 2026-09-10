package main

import "fmt"

type Template struct {
	Name, Content string
	Approved      bool
}
type ReleaseEvent struct{ TemplateName, Status, Detail string }

func Publish(t Template) (ReleaseEvent, error) {
	if t.Name == "" || t.Content == "" {
		return ReleaseEvent{TemplateName: t.Name, Status: "rejected", Detail: "template needs a name and content"}, fmt.Errorf("template is incomplete")
	}
	if !t.Approved {
		return ReleaseEvent{TemplateName: t.Name, Status: "rejected", Detail: "approval is required"}, fmt.Errorf("template is not approved")
	}
	return ReleaseEvent{TemplateName: t.Name, Status: "published", Detail: "approved template released"}, nil
}
