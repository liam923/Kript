package encode

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestToMap(t *testing.T) {
	nonErrorTests := []struct {
		in  interface{}
		out map[string]interface{}
	}{
		{
			in: map[string]struct {
				A string `firestore:"a,omitempty"`
				B int    `firestore:"b"`
				C float32
			}{
				"key": {
					A: "hello",
					B: 10,
					C: 2.4,
				},
			},
			out: map[string]interface{}{
				"key": map[string]interface{}{
					"a": "hello",
					"b": 10,
					"C": float32(2.4),
				},
			},
		},
		{
			in: map[string]struct {
				A string   `firestore:"a,omitempty"`
				B int      `firestore:"b"`
				C []string `firestore:"c,omitempty"`
				D int      `firestore:"d"`
				E int      `firestore:"e,omitempty"`
			}{
				"key": {
					A: "",
					B: 0,
					C: []string{"hello", "hi"},
				},
			},
			out: map[string]interface{}{
				"key": map[string]interface{}{
					"b": 0,
					"c": []interface{}{"hello", "hi"},
					"d": 0,
				},
			},
		},
		{
			in: struct {
				A struct {
					C string `firestore:"c,omitempty"`
				} `firestore:"a,omitempty"`
				B struct {
					C string `firestore:"c,omitempty"`
				} `firestore:"b,omitempty"`
			}{
				A: struct {
					C string `firestore:"c,omitempty"`
				}{
					C: "",
				},
				B: struct {
					C string `firestore:"c,omitempty"`
				}{
					C: "hello",
				},
			},
			out: map[string]interface{}{
				"b": map[string]interface{}{
					"c": "hello",
				},
			},
		},
		{
			in: struct {
				A struct {
					B string `firestore:"b"`
				} `firestore:"a"`
			}{
				A: struct {
					B string `firestore:"b"`
				}{
					B: "hi",
				},
			},
			out: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "hi",
				},
			},
		},
		{
			in: struct {
				a string
			}{
				a: "hi",
			},
			out: map[string]interface{}{},
		},
		{
			in:  struct{}{},
			out: map[string]interface{}{},
		},
		{
			in: struct {
				A time.Time
			}{
				A: time.Unix(0, 0),
			},
			out: map[string]interface{}{
				"A": time.Unix(0, 0),
			},
		},
	}

	errorTests := []interface{}{
		"hello",
		0,
		func() {},
		&struct{}{},
		struct {
			A *string
		}{},
	}

	for i, tt := range nonErrorTests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			r, e := ToMap(tt.in, "firestore")
			if e != nil {
				t.Errorf("unexpected error: %v", e)
			} else if !reflect.DeepEqual(r, tt.out) {
				t.Errorf("expected %v, got %v", tt.out, r)
			}
		})
	}
	for i, tt := range errorTests {
		t.Run(fmt.Sprintf("%v", i+len(nonErrorTests)), func(t *testing.T) {
			r, e := ToMap(tt, "firestore")
			if e == nil || r != nil {
				t.Errorf("expected error, got: %v", r)
			}
		})
	}
}
