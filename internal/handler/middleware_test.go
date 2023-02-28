package handler

import (
	"bytes"
	"errors"
	"mfuss/configs"
	"mfuss/internal/mock"
	"mfuss/internal/repositories"
	"mfuss/internal/service"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandler_checkURLSID(t *testing.T) {
	type fields struct {
		Service *service.Service
		Router  *gin.Engine
	}
	type args struct {
		c    *gin.Context
		body []byte
	}
	type want struct {
		statusCode         int
		urlIDSliceCtxExist bool
		urlIDSliceCtx      any
	}
	cfg := configs.Config{SrvAddr: "localhost:8080", BaseURL: "http://localhost:8080"}

	storage := mock.NewStorageMock(&cfg)
	rep := &repositories.Repository{
		URLStorager: storage,
		Config:      &cfg,
	}
	service := service.NewService(rep)
	rr := httptest.NewRecorder()
	ctx, Router := gin.CreateTestContext(rr)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "cant decode",
			fields: fields{Service: service, Router: Router},
			args:   args{c: ctx, body: []byte(`1928719827491`)},
			want:   want{statusCode: http.StatusBadRequest, urlIDSliceCtx: nil, urlIDSliceCtxExist: false},
		},
		{
			name:   "wrong input",
			fields: fields{Service: service, Router: Router},
			args:   args{c: ctx, body: []byte(`["asdff","asdfff"]`)},
			want:   want{statusCode: http.StatusBadRequest, urlIDSliceCtx: nil, urlIDSliceCtxExist: false},
		},
		{
			name:   "correct input",
			fields: fields{Service: service, Router: Router},
			args:   args{c: ctx, body: []byte(`["asdff","asdfw"]`)},
			want:   want{statusCode: http.StatusBadRequest, urlIDSliceCtx: []string{"asdff", "asdfw"}, urlIDSliceCtxExist: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Service: tt.fields.Service,
				Router:  tt.fields.Router,
			}
			body := tt.args.body
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", bytes.NewReader(body))

			tt.args.c.Request = req

			h.checkURLSID(tt.args.c)
			got := rr.Result().StatusCode
			if got != tt.want.statusCode {
				t.Errorf("checkURLSID() status = %v, want %v", got, tt.want.statusCode)
			}
			_, ok := tt.args.c.Get(urlIDSliceCtx)
			if ok != tt.want.urlIDSliceCtxExist {
				t.Errorf("checkURLSID() urlIDSliceCtxExist = %v, want %v", ok, tt.want.urlIDSliceCtxExist)
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
