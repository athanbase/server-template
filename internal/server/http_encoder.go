package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Response struct {
	Code    int             `json:"code" form:"code"`
	Message string          `json:"message" form:"message"`
	Data    json.RawMessage `json:"data" form:"data"`
}

type ErrResponse struct {
	Code    int    `json:"code" form:"code"`
	Message string `json:"message" form:"message"`
	Reason  string `json:"reason" form:"reason"`
	Data    any    `json:"data" form:"data"`
}

var marshalOpt = protojson.MarshalOptions{
	UseProtoNames:   false,
	EmitUnpopulated: true,
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	se := errors.FromError(err)
	if se.Code == http.StatusTooManyRequests {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(se.Reason))
		return
	}

	reply := &ErrResponse{
		Code:    int(se.Code),
		Message: se.Message,
		Reason:  se.Reason,
		Data:    nil,
	}

	codec, _ := khttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", strings.Join([]string{"application", codec.Name()}, "/"))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func ResponseEncoder(w http.ResponseWriter, r *http.Request, v any) error {
	if v == nil {
		return nil
	}

	if rd, ok := v.(khttp.Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}

	reply := &Response{
		Code:    0,
		Message: "ok",
	}

	codec, _ := khttp.CodecForRequest(r, "Accept")
	var (
		b    []byte
		err  error
		data []byte
	)
	switch m := v.(type) {
	case proto.Message:
		b, err = marshalOpt.Marshal(m)
	case json.Marshaler:
		b, err = m.MarshalJSON()
	default:
		b, err = json.Marshal(m)
	}
	if err != nil {
		return err
	}

	reply.Data = json.RawMessage(b)

	data, err = codec.Marshal(reply)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", strings.Join([]string{"application", codec.Name()}, "/"))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// RequestDecoder decodes the request body to object.
func RequestDecoder(r *http.Request, v interface{}) error {
	codec, ok := khttp.CodecForRequest(r, "Content-Type")
	if !ok {
		log.Warn(errors.BadRequest("CODEC", fmt.Sprintf("unregister Content-Type: %s", r.Header.Get("Content-Type"))))
	}
	data, err := io.ReadAll(r.Body)

	// reset body.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	if len(data) == 0 {
		return nil
	}

	if err = codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", fmt.Sprintf("body unmarshal %s", err.Error()))
	}

	return nil
}
