package languages

import (
	"testing"

	"github.com/ddkwork/golibrary/assert"
	"github.com/ddkwork/golibrary/mylog"
)

func TestNewLanguagesObj(t *testing.T) {
	mylog.Call(func() {
		path := "languages.go"
		l := NewLanguages()

		assert.Equal(t, GoKind, l.CodeFile2Language(path))
		assert.Equal(t, GoKind, CodeFile2Language(path))

		assert.Equal(t, GoKind, Code2Language(goCode))
	})
}

var goCode = `
package main


import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}
`

var pyCode = `
def main():
    print("Hello, world!")
`

var javaScriptCode = `
const para = document.querySelector("p");

para.addEventListener("click", updateName);

function updateName() {
  const name = prompt("Enter a new name");
}
`

var javaCode = `
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, world!");
    }
}
`

var jsCode = `
var _typeof = require("./typeof.js")["default"];
function _toPrimitive(input, hint) {
  if (_typeof(input) !== "object" || input === null) return input;
  var prim = input[Symbol.toPrimitive];
  if (prim !== undefined) {
    var res = prim.call(input, hint || "default");
    if (_typeof(res) !== "object") return res;
    throw new TypeError("@@toPrimitive must return a primitive value.");
  }
  return (hint === "string" ? String : Number)(input);
}
module.exports = _toPrimitive, module.exports.__esModule = true, module.exports["default"] = module.exports;
`
