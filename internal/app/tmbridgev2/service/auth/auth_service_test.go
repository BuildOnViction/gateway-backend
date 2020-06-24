package auth

import (
	"fmt"
	"testing"
)

func Test_verifySig(t *testing.T) {
	type args struct {
		from   string
		sigHex string
		msg    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				from:   "0xd106159eC58BD2EAf5B62eF4e9cDb286170B0Bb9",
				msg:    "855d4d6b2eaf8997b3c2f3e790b0b42f4c9fdaf8de0d728986493af2f70e0db3",
				sigHex: "0xe05f295e3ded0b25fb4a0d3329afa0690feacfd2057607ed3275c4c0ea01cb092a892d31c1254a1e405d1673ac140082c359796c547c40bda5deb18501ddf2d61c",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifySig(tt.args.from, tt.args.sigHex, tt.args.msg)
			fmt.Println("err ", err)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifySig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("verifySig() = %v, want %v", got, tt.want)
			}
		})
	}
}
