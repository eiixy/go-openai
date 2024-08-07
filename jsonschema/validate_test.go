package jsonschema_test

import (
	"github.com/sashabaranov/go-openai/jsonschema"
	"testing"
)

func Test_Validate(t *testing.T) {
	type args struct {
		data   interface{}
		schema jsonschema.Definition
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// string integer number boolean
		{"", args{data: "ABC", schema: jsonschema.Definition{Type: jsonschema.String}}, true},
		{"", args{data: 123, schema: jsonschema.Definition{Type: jsonschema.String}}, false},
		{"", args{data: 123, schema: jsonschema.Definition{Type: jsonschema.Integer}}, true},
		{"", args{data: 123.4, schema: jsonschema.Definition{Type: jsonschema.Integer}}, false},
		{"", args{data: "ABC", schema: jsonschema.Definition{Type: jsonschema.Number}}, false},
		{"", args{data: 123, schema: jsonschema.Definition{Type: jsonschema.Number}}, true},
		{"", args{data: false, schema: jsonschema.Definition{Type: jsonschema.Boolean}}, true},
		{"", args{data: 123, schema: jsonschema.Definition{Type: jsonschema.Boolean}}, false},
		{"", args{data: nil, schema: jsonschema.Definition{Type: jsonschema.Null}}, true},
		{"", args{data: 0, schema: jsonschema.Definition{Type: jsonschema.Null}}, false},
		// array
		{"", args{data: []any{"a", "b", "c"}, schema: jsonschema.Definition{
			Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.String}},
		}, true},
		{"", args{data: []any{1, 2, 3}, schema: jsonschema.Definition{
			Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.String}},
		}, false},
		{"", args{data: []any{1, 2, 3}, schema: jsonschema.Definition{
			Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.Integer}},
		}, true},
		{"", args{data: []any{1, 2, 3.4}, schema: jsonschema.Definition{
			Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.Integer}},
		}, false},
		// object
		{"", args{data: map[string]any{
			"string":  "abc",
			"integer": 123,
			"number":  123.4,
			"boolean": false,
			"array":   []any{1, 2, 3},
		}, schema: jsonschema.Definition{Type: jsonschema.Object, Properties: map[string]jsonschema.Definition{
			"string":  {Type: jsonschema.String},
			"integer": {Type: jsonschema.Integer},
			"number":  {Type: jsonschema.Number},
			"boolean": {Type: jsonschema.Boolean},
			"array":   {Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.Number}},
		},
			Required: []string{"string"},
		}}, true},
		{"", args{data: map[string]any{
			"integer": 123,
			"number":  123.4,
			"boolean": false,
			"array":   []any{1, 2, 3},
		}, schema: jsonschema.Definition{Type: jsonschema.Object, Properties: map[string]jsonschema.Definition{
			"string":  {Type: jsonschema.String},
			"integer": {Type: jsonschema.Integer},
			"number":  {Type: jsonschema.Number},
			"boolean": {Type: jsonschema.Boolean},
			"array":   {Type: jsonschema.Array, Items: &jsonschema.Definition{Type: jsonschema.Number}},
		},
			Required: []string{"string"},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonschema.Validate(tt.args.schema, tt.args.data); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		schema  jsonschema.Definition
		content []byte
		v       any
	}
	var result1 struct {
		String string  `json:"string"`
		Number float64 `json:"number"`
	}
	var result2 struct {
		String string  `json:"string"`
		Number float64 `json:"number"`
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{
			schema: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"string": {Type: jsonschema.String},
					"number": {Type: jsonschema.Number},
				},
			},
			content: []byte(`{"string":"abc","number":123.4}`),
			v:       &result1,
		}, false},
		{"", args{
			schema: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"string": {Type: jsonschema.String},
					"number": {Type: jsonschema.Number},
				},
				Required: []string{"string", "number"},
			},
			content: []byte(`{"string":"abc"}`),
			v:       result2,
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsonschema.Unmarshal(tt.args.schema, tt.args.content, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				t.Logf("Unmarshal() v = %+v\n", tt.args.v)
			}
		})
	}
}