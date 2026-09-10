package main

import "testing"

func TestPublishRequiresApproval(t *testing.T) {
	cases := []struct {
		name     string
		approved bool
		status   string
	}{{"draft", false, "rejected"}, {"approved", true, "published"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			event, _ := Publish(Template{Name: "login", Content: "Code {{code}}", Approved: tc.approved})
			if event.Status != tc.status {
				t.Fatalf("status=%q, want %q", event.Status, tc.status)
			}
		})
	}
}
