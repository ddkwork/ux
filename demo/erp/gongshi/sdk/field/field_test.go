package field

import (
	"testing"

	"github.com/ddkwork/golibrary/std/safemap"
	"github.com/ddkwork/golibrary/std/stream"
)

func TestName(t *testing.T) {
	stream.NewGeneratedFile().EnumTypes("Field", safemap.NewStringerKeys([]string{
		"text",
		"multiLineText",
		"number",
		"singleSelect",
		"multipleSelect",
		"dateTime",
		"formula",
		"attachment",
		"link",
		"user",
		"phone",
		"email",
		"checkbox",
		"url",
		"currency",
		"percent",
	}, true))
}
