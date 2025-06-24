package stream

import (
	"fmt"
	"strings"

	"github.com/ddkwork/golibrary/std/mylog"
)

type (
	object struct {
		Major uint64
		Minor uint64
		Patch uint64
		Build uint64
	}
)

func NewVersion[T string | uint64](s ...T) (v *object) {
	v = &object{
		Major: 0,
		Minor: 0,
		Patch: 0,
		Build: 0,
	}
	switch len(s) {
	case 1:

	case 3:

	case 4:

	default:
		mylog.Check("err arg")
	}
	switch data := any(s[0]).(type) {
	case string:
		split := strings.Split(data, `.`)
		mylog.CheckNil(split)
		o := &object{
			Major: ParseUint(split[0]),
			Minor: ParseUint(split[1]),
			Patch: ParseUint(split[2]),
			Build: 0,
		}
		if len(split) == 4 {
			o.Build = ParseUint(split[3])
		}
		return o
	case uint64:
		o := &object{
			Major: any(s[0]).(uint64),
			Minor: any(s[1]).(uint64),
			Patch: any(s[2]).(uint64),
			Build: 0,
		}
		if len(s) == 4 {
			o.Build = any(s[3]).(uint64)
		}
		return o
	default:
		return
	}
}

func (o *object) String() string {
	sprintf := fmt.Sprintf("%d.%d.%d", o.Major, o.Minor, o.Patch)
	if o.Build > 0 {
		sprintf += fmt.Sprintf(".%d", o.Build)
	}
	return sprintf
}
