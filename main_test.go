package jsontomd

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func Test_decode(t *testing.T) {
	type args struct {
		decoder *json.Decoder
	}
	tests := []struct {
		name    string
		args    args
		want    jsonArray
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				decoder: json.NewDecoder(strings.NewReader(`[
					{"name": "John", "age": 30},
					{"name": "ssss", "age": 10, "hoge": "fuga"},
					{"name": "ssss", "age": 10}
				]`)),
			},
			want: jsonArray{
				{
					{"name", "John"},
					{"age", float64(30)},
				},
				{
					{"name", "ssss"},
					{"age", float64(10)},
					{"hoge", "fuga"},
				},
				{
					{"name", "ssss"},
					{"age", float64(10)},
				},
			},
		},
		{
			name: "key is not quoted, return error",
			args: args{
				decoder: json.NewDecoder(strings.NewReader(`[
					{"name": "John", age: 30},
				]`)),
			},
			wantErr: true,
		},
		{
			name: "duplicate key, return all key value pair",
			args: args{
				decoder: json.NewDecoder(strings.NewReader(`[
					{"name": "John", "name": "John2"}
				]`)),
			},
			want: jsonArray{
				{
					{"name", "John"},
					{"name", "John2"},
				},
			},
		},
		{
			name: "start not [, return error",
			args: args{
				decoder: json.NewDecoder(strings.NewReader(`{
					{"key", "value"}
				]`)),
			},
			wantErr: true,
		},
		{
			name: "end not [, return error",
			args: args{
				decoder: json.NewDecoder(strings.NewReader(`[
					{"key", "value"}
				]`)),
			},
			wantErr: true,
		},
		{
			name: "start invalid start token, return error",
			args: args{
				decoder: json.NewDecoder(strings.NewReader(`a
					{"key", "value"}
				]`)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeArray(tt.args.decoder)
			if (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
