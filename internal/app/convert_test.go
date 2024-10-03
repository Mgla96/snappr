package app

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_extractJSON(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "extract json with resp that has other text",
			args: args{
				"hello {'foo': 'bar'}",
			},
			want: "{'foo': 'bar'}",
		},
		{
			name: "no json { character found",
			args: args{
				"foobar",
			},
			want: "",
		},
		{
			name: "no json } character found",
			args: args{
				"{'foo': 'bar'",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractJSON(tt.args.response); got != tt.want {
				t.Errorf("extractJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unmarshalTo(t *testing.T) {
	examplePrReviewMap := map[string][]PRCommentInfo{
		"foo": {
			{
				CommentBody: "bar",
				StartLine:   1,
				Line:        2,
			},
		},
	}
	prReviewMapBytes, err := json.Marshal(examplePrReviewMap)
	if err != nil {
		t.Fatalf("failed to marshal examplePrReviewMap: %v", err)
	}

	examplePRCreation := PRCreation{
		Title: "foo",
		Body:  "bar",
		UpdatedFiles: []PRCreationFile{
			{
				Path:          "foo",
				FullContent:   "bar",
				CommitMessage: "baz",
			},
		},
	}
	prCreationBytes, err := json.Marshal(examplePRCreation)
	if err != nil {
		t.Fatalf("failed to marshal examplePrReviewMap: %v", err)
	}

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "unmarshalToMap",
			args: args{
				data: []byte(`{"key": "value"}`),
			},
			want: map[string]string{
				"key": "value",
			},
		},
		{
			name: "unmarshalToMap failure",
			args: args{
				data: []byte(`not json`),
			},
			wantErr: true,
		},
		{
			name: "unsucessful unmarshal",
			args: args{
				data: []byte(`not json`),
			},
			wantErr: true,
		},
		{
			name: "json but not a PRReviewMap",
			args: args{
				data: []byte(`{"foo": "bar"}`),
			},
			wantErr: true,
		},
		{
			name: "unmarshal to PRReviewMap",
			args: args{
				data: prReviewMapBytes,
			},
			want:    examplePrReviewMap,
			wantErr: false,
		},
		{
			name: "unsucessful unmarshal",
			args: args{
				data: []byte(`not json`),
			},
			want:    PRCreation{},
			wantErr: true,
		},
		{
			name: "json but not PRCreation",
			args: args{
				data: []byte(`{"foo": "bar"}`),
			},
			want:    PRCreation{},
			wantErr: false,
		},
		{
			name: "unmarshal to PRCreation",
			args: args{
				data: prCreationBytes,
			},
			want:    examplePRCreation,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.want.(type) {
			case map[string]string:
				got, err := unmarshalTo[map[string]string](tt.args.data)
				if (err != nil) != tt.wantErr {
					t.Errorf("unmarshalTo() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("unmarshalTo() = %v, want %v", got, tt.want)
				}

			case PRReviewMap:
				got, err := unmarshalTo[PRReviewMap](tt.args.data)
				if (err != nil) != tt.wantErr {
					t.Errorf("unmarshalTo() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("unmarshalTo() = %v, want %v", got, tt.want)
				}
			case PRCreation:
				got, err := unmarshalTo[PRCreation](tt.args.data)
				if (err != nil) != tt.wantErr {
					t.Errorf("unmarshalTo() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("unmarshalTo() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}
