// Originally from https://gl.oddhunters.com/pub/fiberzap
// Copyright (apparently) by Ozgur Boru <boruozgur@yandex.com.tr>
// and "mert" (https://gl.oddhunters.com/mert)
package fiberzap

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap/zapcore"
)

func getAllowedHeaders() map[string]bool {
	return map[string]bool{
		"User-Agent": true,
		"X-Mobile":   true,
	}
}

type resp struct {
	code  int
	_type string
}

func Resp(r *fasthttp.Response) *resp {
	return &resp{
		code:  r.StatusCode(),
		_type: bytes.NewBuffer(r.Header.ContentType()).String(),
	}
}

func (r *resp) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("type", r._type)
	enc.AddInt("code", r.code)

	return nil
}

type req struct {
	body     string
	fullPath string
	user     string
	ip       string
	method   string
	route    string
	headers  *headerbag
}

func Req(c *fiber.Ctx) *req {
	reqq := c.Request()
	var body []byte
	buffer := new(bytes.Buffer)
	err := json.Compact(buffer, reqq.Body())
	if err != nil {
		body = reqq.Body()
	} else {
		body = buffer.Bytes()
	}

	headers := &headerbag{
		vals: make(map[string]string),
	}
	allowedHeaders := getAllowedHeaders()
	reqq.Header.VisitAll(func(key, val []byte) {
		k := bytes.NewBuffer(key).String()
		if _, exist := allowedHeaders[k]; exist {
			headers.vals[strings.ToLower(k)] = bytes.NewBuffer(val).String()
		}
	})

	var userEmail string
	if u := c.Locals("userEmail"); u != nil {
		userEmail = u.(string)
	}

	return &req{
		body:     bytes.NewBuffer(body).String(),
		fullPath: bytes.NewBuffer(reqq.RequestURI()).String(),
		headers:  headers,
		ip:       c.IP(),
		method:   c.Method(),
		route:    c.Route().Path,
		user:     userEmail,
	}
}

func (r *req) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("fullPath", r.fullPath)
	enc.AddString("ip", r.ip)
	enc.AddString("method", r.method)
	enc.AddString("route", r.route)

	if r.body != "" {
		enc.AddString("body", r.body)
	}

	if r.user != "" {
		enc.AddString("user", r.user)
	}

	err := enc.AddObject("headers", r.headers)
	if err != nil {
		return err
	}

	return nil
}

type headerbag struct {
	vals map[string]string
}

func (h *headerbag) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range h.vals {
		enc.AddString(k, v)
	}

	return nil
}
