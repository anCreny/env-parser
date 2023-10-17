# env-parser

## About
This is a simple parser to run through a struct and get env variables appropriate for 'custom' tag value.


## Installing
```bash
go get github.com/anCreny/env-parser
```

## Usage
To use custom tag you should create **EnvParser** first:
```go
func New(tag string, useName, safe bool) *EnvParser
```
* tag - tag name that will be parsed by.
* useName -  if a field of a structure want to parse doesn't have the tag, you can manage **parser** behaviour:
  * *true* - **parser** use field name like a tag value.
  * *false* - ignore field if tag with the given name doesn't exist.
* safe - getting opportunity to manage fields override. If the specific field has a non-zero value:
  * *true* - value from env var ___won't___ replace the field value
  * *false* - value from env var ___will___ replace the field value

The constructor will return a parser instance of the structure:
```go
type EnvParser struct {
	tag string
	useName bool
}
```

The structure has only one simple method:
```go
func (e *EnvParser) Parse(structure interface{}) error
```

Method get in a pointer to a structure you want to parse in, instead you will get an error.

***

You'r wellcome to contribute and comment it!

Love you <3
