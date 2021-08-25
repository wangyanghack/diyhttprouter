package diyhttprouter

import (
	"testing"
)

func Test_updatePath(t *testing.T) {
	type args struct {
		path string
		l    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				path: "abcdef",
				l:    3,
			},
			want: "def",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updatePath(tt.args.path, tt.args.l); got != tt.want {
				t.Errorf("updatePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFromIndices(t *testing.T) {
	type args struct {
		indices string
		b       byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				indices: "abcdefg",
				b:       'c',
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFromIndices(tt.args.indices, tt.args.b); got != tt.want {
				t.Errorf("getFromIndices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_node_insertNode(t *testing.T) {
	type fields struct {
		path     string
		indices  string
		children []*node
		handle   Handle
	}
	type args struct {
		fullPath string
		handle   Handle
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				path: "abcdef",
			},
			args: args{
				fullPath: "abc",
			},
		},
		{
			name: "test2",
			fields: fields{
				path: "abc",
			},
			args: args{
				fullPath: "abcdef",
			},
		},
		{
			name: "test3",
			fields: fields{
				path: "abcdef",
			},
			args: args{
				fullPath: "abchgi",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node{
				path:     tt.fields.path,
				indices:  tt.fields.indices,
				children: tt.fields.children,
				handle:   tt.fields.handle,
			}
			if err := n.insertNode(tt.args.fullPath, tt.args.handle); (err != nil) != tt.wantErr {
				t.Errorf("node.insertNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
