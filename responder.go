package responder

import (
	"net/http"
	"path/filepath"
	"strings"
)

// MimeTypes registered mime types
var MimeTypes = map[string]string{}

// Register new mime type and format
func Register(mime string, format string) {
	MimeTypes[mime] = format
}

func init() {
	for mimeType, format := range map[string]string{
		"text/html":        "html",
		"application/json": "json",
		"application/xml":  "xml",
	} {
		Register(mimeType, format)
	}
}

type Responder struct {
	responds map[string]func()
}

// With support string or []string as formats, With("html", func() {
//   writer.Write([]byte("this is a html request"))
// }).With([]string{"json", "xml"}, func() {
//   writer.Write([]byte("this is a json or xml request"))
// })
func With(format interface{}, fc func()) *Responder {
	rep := &Responder{responds: map[string]func(){}}
	return rep.With(format, fc)
}

func (rep *Responder) With(format interface{}, fc func()) *Responder {
	if f, ok := format.(string); ok {
		rep.responds[f] = fc
	} else if fs, ok := format.([]string); ok {
		for _, f := range fs {
			rep.responds[f] = fc
		}
	}
	return rep
}

// Respond differently according to request's accepted mime type
func (rep *Responder) Respond(request *http.Request) {
	// get request format from url
	if ext := filepath.Ext(request.URL.Path); ext != "" {
		if respond, ok := rep.responds[strings.TrimPrefix(ext, ".")]; ok {
			respond()
			return
		}
	}

	// get request format from Accept
	for _, accept := range strings.Split(request.Header.Get("Accept"), ",") {
		if format, ok := MimeTypes[accept]; ok {
			if respond, ok := rep.responds[format]; ok {
				respond()
				return
			}
		}
	}

	// use first format as default
	for _, respond := range rep.responds {
		respond()
		break
	}
	return
}
