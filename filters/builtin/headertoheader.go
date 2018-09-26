package builtin

import (
	"fmt"

	"github.com/zalando/skipper/filters"
)

const (
	ActionUnknown         = "actionUnknown"
	HeaderToAddHeaderName = "headerToAddHeader"
	HeaderToSetHeaderName = "headerToSetHeader"
)

type headerUpdateType int

const (
	headerToAddHeader headerUpdateType = iota
	headerToSetHeader
)

type (
	headerToHeaderFilter struct {
		typ           headerUpdateType
		oldHeaderName string
		newHeaderName string
		formatString  string
	}
)

// NewHeaderToAddHeader creates a filter which converts a header from the incoming
// Request to another header
//
//     headerToAddHeader("X-Bar-Header", "X-Foo-Header")
//
// The above filter will set the value of "X-Foo-Header" header to the
// value of "X-Bar-Header" header, and add it to the request and will not
// override the value if the header exists already
//
// The header value can be created by a formatstring with an optional third parameter
//
//     headerToAddHeader("X-Bar-Header", "X-Foo-Header", "prefix %s postfix")
//     headerToAddHeader("X-Access-Token", "Authorization", "Bearer %s")
//
func NewHeaderToAddHeader() filters.Spec {
	return &headerToHeaderFilter{typ: headerToAddHeader}
}

// NewHeaderToSetHeader creates a filter which converts a header from the incoming
// Request to another header
//
//     headerToSetHeader("X-Bar-Header", "X-Foo-Header")
//
// The above filter will set the value of "X-Foo-Header" header to the
// value of "X-Bar-Header" header, to the request and will override the value if
// the header exists already
//
// The header value can be created by a formatstring with an optional third parameter
//
//     headerToSetHeader("X-Bar-Header", "X-Foo-Header", "prefix %s postfix")
//     headerToSetHeader("X-Access-Token", "Authorization", "Bearer %s")
//
func NewHeaderToSetHeader() filters.Spec {
	return &headerToHeaderFilter{typ: headerToSetHeader}
}

func (s *headerToHeaderFilter) Name() string {
	switch s.typ {
	case headerToAddHeader:
		return HeaderToAddHeaderName
	case headerToSetHeader:
		return HeaderToSetHeaderName
	}
	return ActionUnknown
}

// CreateFilter creates a `headerToHeader` filter instance with below signature
// s.CreateFilter("X-Bar-Header", "X-Foo-Header")
func (*headerToHeaderFilter) CreateFilter(args []interface{}) (filters.Filter, error) {
	if l := len(args); l < 2 || l > 3 {
		return nil, filters.ErrInvalidFilterParameters
	}

	o, ok := args[0].(string)
	if !ok {
		return nil, filters.ErrInvalidFilterParameters
	}

	n, ok := args[1].(string)
	if !ok {
		return nil, filters.ErrInvalidFilterParameters
	}

	formatString := "%s"
	if len(args) == 3 {
		formatString, ok = args[2].(string)
		if !ok {
			return nil, filters.ErrInvalidFilterParameters
		}
	}

	return &headerToHeaderFilter{oldHeaderName: o, newHeaderName: n, formatString: formatString}, nil
}

// String prints nicely the headerToHeaderFilter configuration based on the
// configuration and check used.
func (f *headerToHeaderFilter) String() string {
	switch f.typ {
	case headerToAddHeader:
		return fmt.Sprintf("%s(%s, %s, %s)", HeaderToAddHeaderName, f.oldHeaderName, f.newHeaderName, f.formatString)
	case headerToSetHeader:
		return fmt.Sprintf("%s(%s, %s, %s)", HeaderToSetHeaderName, f.oldHeaderName, f.newHeaderName, f.formatString)
	}
	return ActionUnknown
}

func (f *headerToHeaderFilter) Request(ctx filters.FilterContext) {
	req := ctx.Request()

	v := req.Header.Get(f.oldHeaderName)
	if v == "" {
		return
	}

	switch f.typ {
	case headerToAddHeader:
		req.Header.Add(f.newHeaderName, fmt.Sprintf(f.formatString, v))
	case headerToSetHeader:
		req.Header.Set(f.newHeaderName, fmt.Sprintf(f.formatString, v))
	}
}

func (*headerToHeaderFilter) Response(ctx filters.FilterContext) {}
