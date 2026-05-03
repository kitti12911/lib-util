package main

import "testing"

func TestSnake(t *testing.T) {
	tests := map[string]string{
		"UserID":    "user_id",
		"HTTPSPort": "https_port",
		"URLValue":  "url_value",
		"Line1":     "line1",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			if got := snake(input); got != want {
				t.Fatalf("snake(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

func TestIsChildRelation(t *testing.T) {
	tests := map[string]bool{
		"rel:has-one,join:id=user_id":        true,
		"rel:has-many,join:id=user_id":       true,
		"rel:belongs-to,join:user_id=id":     false,
		"rel:many-to-many,join:user_id=id":   false,
		"column_name,type:uuid,default:null": false,
	}

	for tag, want := range tests {
		t.Run(tag, func(t *testing.T) {
			if got := isChildRelation(tag); got != want {
				t.Fatalf("isChildRelation(%q) = %t, want %t", tag, got, want)
			}
		})
	}
}
