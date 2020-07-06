package common

import (
	"reflect"
	"testing"
)

func TestDifference(t *testing.T) {
	type args struct {
		slice1 []string
		slice2 []string
	}
	tests := []struct {
		name         string
		args         args
		wantInSlice1 []string
		wantInSlice2 []string
	}{
		{
			name: "test 1",
			args: args{
				slice1: []string{"hello"},
				slice2: []string{"world"},
			},
			wantInSlice1: []string{"hello"},
			wantInSlice2: []string{"world"},
		},
		{
			name: "test 1",
			args: args{
				slice1: []string{"hello", "1", "2"},
				slice2: []string{"1", "world", "3"},
			},
			wantInSlice1: []string{"hello", "2"},
			wantInSlice2: []string{"world", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInSlice1, gotInSlice2 := Difference(tt.args.slice1, tt.args.slice2)
			if !reflect.DeepEqual(gotInSlice1, tt.wantInSlice1) {
				t.Errorf("Difference() gotInSlice1 = %v, want %v", gotInSlice1, tt.wantInSlice1)
			}
			if !reflect.DeepEqual(gotInSlice2, tt.wantInSlice2) {
				t.Errorf("Difference() gotInSlice2 = %v, want %v", gotInSlice2, tt.wantInSlice2)
			}
		})
	}
}
