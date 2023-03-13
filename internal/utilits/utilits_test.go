package utilits

import (
	"reflect"
	"testing"
)

func TestCheckURL(t *testing.T) {
	type args struct {
		shortURLId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "vallid URL",
			args:    args{shortURLId: "https://kanobu.ru/"},
			wantErr: false,
		},
		{
			name:    "invallid URL",
			args:    args{shortURLId: "asgasgasgasgq/"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckURL(tt.args.shortURLId); (err != nil) != tt.wantErr {
				t.Errorf("CheckURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenetareID(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Kind
	}{
		{
			name: "type check",
			want: reflect.String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenetareID(); reflect.TypeOf(got).Kind() != tt.want {
				t.Errorf("GenetareID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGenetareID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenetareID()
	}
}
