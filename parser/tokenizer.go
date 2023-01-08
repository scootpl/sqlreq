package parser

import "strings"

type tokenType string

type token struct {
	Type  tokenType
	Value string
}

func (p *Parser) Tokenizer(query string) []token {
	tokens := []token{}

	fields := strings.Fields(strings.ToLower(query))
	for _, field := range fields {
		var t token
		switch tokenType(field) {
		case from:
			t = token{Type: from, Value: field}
		case where:
			t = token{Type: where, Value: field}
		case and:
			t = token{Type: and, Value: field}
		case payload:
			t = token{Type: payload, Value: field}
		case header:
			t = token{Type: header, Value: field}
		case with:
			t = token{Type: with, Value: field}
		case get:
			t = token{Type: get, Value: field}
		case post:
			t = token{Type: post, Value: field}
		case delete:
			t = token{Type: delete, Value: field}
		case patch:
			t = token{Type: patch, Value: field}
		case put:
			t = token{Type: put, Value: field}
		case comma:
			t = token{Type: comma, Value: field}
		case equal:
			t = token{Type: equal, Value: field}
		case question:
			t = token{Type: question, Value: field}
		case basicauth:
			t = token{Type: basicauth, Value: field}
		case timeout:
			t = token{Type: timeout, Value: field}

		default:
			if len(field) > 0 && strings.HasPrefix(field, "'") {
				field = field[1:]
			}
			if len(field) > 0 && strings.HasSuffix(field, "'") {
				field = field[:len(field)-1]
			}
			if len(field) > 0 && strings.HasPrefix(field, ",") {
				tokens = append(tokens, token{
					Type:  comma,
					Value: ",",
				})
				field = field[1:]
			}
			if len(field) > 0 && strings.HasSuffix(field, ",") {
				field = field[:len(field)-1]
				tokens = append(tokens, token{
					Type:  value,
					Value: field,
				})
				t = token{
					Type:  comma,
					Value: ",",
				}
			} else {
				t = token{
					Type:  value,
					Value: field,
				}
			}
		}

		tokens = append(tokens, t)
	}

	return tokens
}
