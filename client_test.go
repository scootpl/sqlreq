package sqlreq

import (
	"reflect"
	"testing"

	"github.com/scootpl/sqlreq/parser"
)

func Test_doSelect(t *testing.T) {
	type args struct {
		mode   int
		query  string
		result *result
		params []any
	}
	tests := []struct {
		name       string
		args       args
		wantResult *result
		wantErr    bool
	}{
		{
			name: "test1",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http",
				result: &result{},
			},
			wantResult: &result{
				jsonParam:   make(map[string]any),
				headerParam: make(map[string]string),
				url:         "http",
				method:      "GET",
				timeout:     30,
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where",
				result: &result{},
			},
			wantResult: &result{},
			wantErr:    true,
		},
		{
			name: "test3",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where header test1 = test2",
				result: &result{},
			},
			wantResult: &result{
				jsonParam: make(map[string]any),
				headerParam: map[string]string{
					"test1": "test2",
				},
				url:     "http",
				method:  "GET",
				timeout: 30,
			},
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where test1 = test2",
				result: &result{},
			},
			wantResult: &result{
				jsonParam: map[string]any{
					"test1": "test2",
				},
				headerParam: make(map[string]string),
				url:         "http",
				method:      "GET",
				timeout:     30,
			},
			wantErr: false,
		},
		{
			name: "test5",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where test1 = test2 and header h1 = h2",
				result: &result{},
			},
			wantResult: &result{
				jsonParam: map[string]any{
					"test1": "test2",
				},
				headerParam: map[string]string{
					"h1": "h2",
				},
				url:     "http",
				method:  "GET",
				timeout: 30,
			},
			wantErr: false,
		},
		{
			name: "test6",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where test1 = test2 and header h1 = h2 and header hh1 = hh2",
				result: &result{},
			},
			wantResult: &result{
				jsonParam: map[string]any{
					"test1": "test2",
				},
				headerParam: map[string]string{
					"h1":  "h2",
					"hh1": "hh2",
				},
				url:     "http",
				method:  "GET",
				timeout: 30,
			},
			wantErr: false,
		},
		{
			name: "test7",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where test1 = test2 and header h1 = h2 and header hh1 = hh2 with post",
				result: &result{},
			},
			wantResult: &result{
				jsonParam: map[string]any{
					"test1": "test2",
				},
				headerParam: map[string]string{
					"h1":  "h2",
					"hh1": "hh2",
				},
				url:     "http",
				method:  "POST",
				timeout: 30,
			},
			wantErr: false,
		},
		{
			name: "test8",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where test1 = test2 and header h1 = h2 and header hh1 = hh2 with post with timeout 10",
				result: &result{},
			},
			wantResult: &result{
				jsonParam: map[string]any{
					"test1": "test2",
				},
				headerParam: map[string]string{
					"h1":  "h2",
					"hh1": "hh2",
				},
				url:     "http",
				method:  "POST",
				timeout: 10,
			},
			wantErr: false,
		},
		{
			name: "test9",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where %s = test2 and header h1 = h2 and header hh1 = hh2 with post with timeout 10",
				result: &result{},
				params: []any{"test1"},
			},
			wantResult: &result{
				jsonParam: map[string]any{
					"test1": "test2",
				},
				headerParam: map[string]string{
					"h1":  "h2",
					"hh1": "hh2",
				},
				url:     "http",
				method:  "POST",
				timeout: 10,
			},
			wantErr: false,
		},
		{
			name: "test10",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where %s = test2 and header h1 = '%s' and header hh1 = hh2 with post with timeout 10",
				result: &result{},
				params: []any{"test1", "h2"},
			},
			wantResult: &result{
				jsonParam: map[string]any{
					"test1": "test2",
				},
				headerParam: map[string]string{
					"h1":  "h2",
					"hh1": "hh2",
				},
				url:     "http",
				method:  "POST",
				timeout: 10,
			},
			wantErr: false,
		},
		{
			name: "test11",
			args: args{
				mode:   parser.SelectBody,
				query:  "from http where %s = test2 header h1 = '%s' and header hh1 = hh2 with post with timeout 10",
				result: &result{},
				params: []any{"test1", "h2"},
			},
			wantResult: &result{
				jsonParam: map[string]any{
					"test1": "test2",
				},
				headerParam: map[string]string{
					"h1":  "h2",
					"hh1": "hh2",
				},
				url:     "http",
				method:  "POST",
				timeout: 10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := doSelect(tt.args.mode, tt.args.query, tt.args.result, tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("doSelect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(tt.args.result, tt.wantResult) {
				t.Errorf("doSelect() result = %#v, wantResult %#v", tt.args.result, tt.wantResult)
			}
		})
	}
}
