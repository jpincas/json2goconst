# json2goconst
Compiles from JSON a Go package consisting of const declarations equal to dot-notation JSON paths.

Useful if you have a JSON file which is a big, nested object ending in text nodes (like a big, hierarchical dictionary) and you want to refer to keys in a type-safe way in your Go code.

Takes a json file like this:

```json
{
  "errors": {
    "test": "not important",
    "test2": "not important"
  },
  "messages": {
    "test": "not important",
    "subMessages": {
      "subMessage1": "not important",
      "subMessage2": "not important"
    }
  }
}
```

and turn it into this:

```go
const (
	Errors_Test = "errors.test"
	Errors_Test2 = "errors.test2"

	Messages_Test = "messages.test"

	Messages_SubMessages_SubMessage1 = "messages.subMessages.subMessage1"
	Messages_SubMessages_SubMessage2 = "messages.subMessages.subMessage2"
)
```

