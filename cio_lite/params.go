package cio_lite

import (
	"bytes"
	"fmt"
	"strings"

	"net/url"
)

type Params struct {
	BodyType         string
	Delimiter        string
	Email            string
	IncludeHeaders   string
	Status           string
	StatusOK         string
	IncludeBody      bool
	IncludeFlags     bool
	IncludeNamesOnly bool
	Limit            int
	Offset           int
}

func (p Params) QueryString() string {
	var buffer bytes.Buffer

	if p.BodyType != "" {
		buffer.WriteString(fmt.Sprintf("body_type=%s&", url.QueryEscape(p.Delimiter)))
	}
	if p.Delimiter != "" {
		buffer.WriteString(fmt.Sprintf("delimiter=%s&", url.QueryEscape(p.Delimiter)))
	}
	if p.Email != "" {
		buffer.WriteString(fmt.Sprintf("email=%s&", url.QueryEscape(p.Email)))
	}
	if p.IncludeHeaders != "" {
		buffer.WriteString(fmt.Sprintf("include_headers=%s&", url.QueryEscape(p.Delimiter)))
	}
	if p.Status != "" {
		buffer.WriteString(fmt.Sprintf("status=%s&", url.QueryEscape(p.Status)))
	}
	if p.StatusOK != "" {
		buffer.WriteString(fmt.Sprintf("status_ok=%s&", url.QueryEscape(p.StatusOK)))
	}

	if p.IncludeBody {
		buffer.WriteString("include_body=1&")
	}
	if p.IncludeFlags {
		buffer.WriteString("include_flags=1&")
	}
	if p.IncludeNamesOnly {
		buffer.WriteString("include_names_only=1&")
	}

	if p.Limit != 0 {
		buffer.WriteString(fmt.Sprintf("limit=%d&", p.Limit))
	}
	if p.Offset != 0 {
		buffer.WriteString(fmt.Sprintf("offset=%d&", p.Offset))
	}

	args := strings.TrimRight(buffer.String(), "&")
	if len(args) == 0 {
		return ""
	} else {
		return fmt.Sprintf("?%s", args)
	}
}
