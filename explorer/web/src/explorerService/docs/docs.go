// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2018-10-08 14:10:43.384088 +0800 CST m=+0.047220450

package docs

import (
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "Tronscan API",
        "title": "Tronscan API",
        "contact": {
            "name": "tron",
            "url": "http://www.swagger.io/support"
        },
        "license": {},
        "version": "1.0"
    },
    "paths": {}
}`

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}
func init() {
	swag.Register(swag.Name, &s{})
}
