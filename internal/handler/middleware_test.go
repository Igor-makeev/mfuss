package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGzipCompress(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GzipCompress(tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GzipCompress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGzipUnpack(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GzipUnpack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GzipUnpack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserCheck(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserCheck(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shouldCompress(t *testing.T) {
	type args struct {
		req *http.Request
	}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	wrongReq := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "correct header",
			args: args{req: req},
			want: true,
		},
		{
			name: "no header",
			args: args{req: wrongReq},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldCompress(tt.args.req); got != tt.want {
				t.Errorf("shouldCompress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shouldUnpack(t *testing.T) {
	type args struct {
		req *http.Request
	}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	req.Header.Set("Content-Encoding", "gzip")
	wrongReq := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "correct header",
			args: args{req: req},
			want: true,
		},
		{
			name: "no header",
			args: args{req: wrongReq},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldUnpack(tt.args.req); got != tt.want {
				t.Errorf("shouldUnpack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateCook(t *testing.T) {
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
			if got := generateCook(); reflect.TypeOf(got).Kind() != tt.want {
				t.Errorf("generateCook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkCook(t *testing.T) {
	type args struct {
		id string
	}
	correctCook := generateCook()
	wrongCook := "k;agf23123"

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "correct cook",
			args: args{correctCook},
			want: true,
		},
		{
			name: "correct cook",
			args: args{wrongCook},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkCook(tt.args.id); got != tt.want {
				t.Errorf("checkCook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUserID(t *testing.T) {
	type args struct {
		c *gin.Context
	}

	rr := httptest.NewRecorder()

	validCtx, _ := gin.CreateTestContext(rr)
	validInput := "1231qweqwe"
	validCtx.Set(userCtx, validInput)

	emptyGinCtx, _ := gin.CreateTestContext(rr)

	ctxWithWrongType, _ := gin.CreateTestContext(rr)

	var wronginput []int
	ctxWithWrongType.Set(userCtx, wronginput)

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "all ok",
			args:    args{c: validCtx},
			want:    validInput,
			wantErr: nil,
		},
		{
			name:    "empty ctx",
			args:    args{c: emptyGinCtx},
			want:    "",
			wantErr: ErrNoUserID,
		},
		{
			name:    "invalid ctx data",
			args:    args{c: ctxWithWrongType},
			want:    "",
			wantErr: ErrInvalidUserID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUserID(tt.args.c)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("getUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUrlsArray(t *testing.T) {
	type args struct {
		c *gin.Context
	}

	rr := httptest.NewRecorder()

	validCtx, _ := gin.CreateTestContext(rr)
	var validInput []string
	validCtx.Set(urlIDSliceCtx, validInput)

	emptyGinCtx, _ := gin.CreateTestContext(rr)

	ctxWithWrongType, _ := gin.CreateTestContext(rr)

	var wronginput []int
	ctxWithWrongType.Set(urlIDSliceCtx, wronginput)
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr error
	}{
		{
			name:    "empty context",
			args:    args{c: emptyGinCtx},
			want:    nil,
			wantErr: ErrNoDataArray,
		},
		{
			name:    "wrong data array",
			args:    args{c: ctxWithWrongType},
			want:    nil,
			wantErr: ErrInvalidDataArray,
		},
		{
			name:    "ok data",
			args:    args{c: validCtx},
			want:    validInput,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUrlsArray(tt.args.c)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("getUrlsArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUrlsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genetareUserID(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Kind
	}{
		{
			name: "return value correctness",
			want: reflect.String,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genetareUserID(); reflect.TypeOf(got).Kind() != tt.want {
				t.Errorf("genetareUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLSIDCheck(t *testing.T) {
	//TODO
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLSIDCheck(tt.args.c)
		})
	}
}
