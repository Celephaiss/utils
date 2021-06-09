package utils

import (
	"reflect"
	"testing"
)

func TestApply(t *testing.T) {
	type args struct {
		collections []interface{}
		f           ApplyFunc
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		// TODO: Add test cases.
		{
			"add 1",
			args{
				collections: []interface{}{1, 2, 3},
				f:           func(i interface{}) interface{} { return i.(int) + 1 },
			},
			[]interface{}{2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Apply(tt.args.collections, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestApply2(t *testing.T) {
	type args struct {
		collections []interface{}
		f           ApplyFunc
		nGoroutine  int
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		// TODO: Add test cases.
		{
			"add 1",
			args{
				collections: []interface{}{1, 2, 3},
				f:           func(i interface{}) interface{} { return i.(int) + 1 },
				nGoroutine:  2,
			},
			[]interface{}{2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Apply2(tt.args.collections, tt.args.f, tt.args.nGoroutine); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply2() = %v, want %v", got, tt.want)
			}
		})
	}
}
