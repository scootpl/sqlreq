package parser

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"

	graphviz "github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type node struct {
	typ  tokenType
	conn []int
	do   func(v string)
}

func (p *Parser) newNode(t tokenType, do func(v string), conn ...int) node {
	return node{
		typ:  t,
		conn: conn,
		do:   do,
	}
}

type Parsed struct {
	OutputHeaders []string
	HeaderKey     []string
	HeaderValue   []string
	JsonKey       []string
	JsonValue     []string
	Url           string
	Method        string
	Payload       bool
	BasicLogin    string
	BasicPass     string
	Timeout       string
}

const (
	SelectBody   = 3
	SelectHeader = 1
)

type Parser struct {
	nodes []node
	Parsed
}

func New(startPoint int) *Parser {
	p := new(Parser)
	p.Method = "GET"
	p.initNodes(startPoint)
	return p
}

func (p *Parser) initNodes(startPoint int) {
	nodes := make([]node, 49)

	if startPoint == 1 {
		nodes[0] = p.newNode(start, nil, 1)
	}

	if startPoint == 3 {
		nodes[0] = p.newNode(start, nil, 3)
	}

	nodes[1] = p.newNode(value, func(v string) { p.OutputHeaders = append(p.OutputHeaders, v) }, 2, 3)
	nodes[2] = p.newNode(comma, nil, 1)
	nodes[3] = p.newNode(from, nil, 4)
	nodes[4] = p.newNode(value, func(v string) { p.Url = v }, 5, 6, 7)
	nodes[5] = p.newNode(eof, nil)
	nodes[6] = p.newNode(where, nil, 9, 10, 11)
	nodes[7] = p.newNode(with, nil, 8, 30, 31, 32, 33, 34, 38)
	nodes[8] = p.newNode(post, func(v string) { p.Method = "POST" }, 5, 7)
	nodes[9] = p.newNode(header, nil, 12)
	nodes[10] = p.newNode(value, func(v string) { p.JsonKey = append(p.JsonKey, v) }, 13)
	nodes[11] = p.newNode(payload, func(v string) { p.Payload = true }, 14)
	nodes[12] = p.newNode(value, func(v string) { p.HeaderKey = append(p.HeaderKey, v) }, 15)
	nodes[13] = p.newNode(equal, nil, 16)
	nodes[14] = p.newNode(equal, nil, 17)
	nodes[15] = p.newNode(equal, nil, 18)
	nodes[16] = p.newNode(value, func(v string) { p.JsonValue = append(p.JsonValue, v) }, 7, 19, 5)
	nodes[17] = p.newNode(question, func(v string) { p.Payload = true }, 7, 20, 5)
	nodes[18] = p.newNode(value, func(v string) { p.HeaderValue = append(p.HeaderValue, v) }, 5, 7, 21)
	nodes[19] = p.newNode(and, nil, 10, 22)
	nodes[20] = p.newNode(and, nil, 23)
	nodes[21] = p.newNode(and, nil, 9, 10, 11)
	nodes[22] = p.newNode(header, nil, 24)
	nodes[23] = p.newNode(header, nil, 25)
	nodes[24] = p.newNode(value, func(v string) { p.HeaderKey = append(p.HeaderKey, v) }, 26)
	nodes[25] = p.newNode(value, func(v string) { p.HeaderKey = append(p.HeaderKey, v) }, 27)
	nodes[26] = p.newNode(equal, nil, 28)
	nodes[27] = p.newNode(equal, nil, 29)
	nodes[28] = p.newNode(value, func(v string) { p.HeaderValue = append(p.HeaderValue, v) }, 7, 19, 5)
	nodes[29] = p.newNode(value, func(v string) { p.HeaderValue = append(p.HeaderValue, v) }, 7, 5, 20)
	nodes[30] = p.newNode(get, func(v string) { p.Method = "GET" }, 5, 7)
	nodes[31] = p.newNode(put, func(v string) { p.Method = "PUT" }, 5, 7)
	nodes[32] = p.newNode(patch, func(v string) { p.Method = "PATCH" }, 5, 7)
	nodes[33] = p.newNode(delete, func(v string) { p.Method = "DELETE" }, 5, 7)
	nodes[34] = p.newNode(basicauth, nil, 35)
	nodes[35] = p.newNode(value, func(v string) { p.BasicLogin = v }, 36)
	nodes[36] = p.newNode(comma, nil, 37)
	nodes[37] = p.newNode(value, func(v string) { p.BasicPass = v }, 5, 7)
	nodes[38] = p.newNode(timeout, nil, 39)
	nodes[39] = p.newNode(value, func(v string) { p.Timeout = v }, 7, 5)
	p.nodes = nodes
}

func (p *Parser) GenerateGraph() {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	nodes := make(map[int]*cgraph.Node)

	for i := range p.nodes {
		node, _ := graph.CreateNode(fmt.Sprintf("%d-%s", i, p.nodes[i].typ))
		nodes[i] = node
	}

	for i, n := range p.nodes {
		for _, c := range n.conn {
			graph.CreateEdge("", nodes[i], nodes[c])
		}
	}

	var buf bytes.Buffer
	if err := g.Render(graph, "dot", &buf); err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
}

func (p *Parser) Parse(tokens []token) error {
	index := 0

	for i, t := range tokens {
		ok := false
		for _, c := range p.nodes[index].conn {
			if t.Type == p.nodes[c].typ {
				ok = true
				index = c
				if p.nodes[c].do != nil {
					p.nodes[c].do(t.Value)
				}
				break
			}
		}
		if !ok {
			var expectedTokens []string
			for _, c := range p.nodes[index].conn {
				expectedTokens = append(expectedTokens, string(p.nodes[c].typ))
			}

			return fmt.Errorf("parse error at %d token (%s), expected: %s", i+1, t.Type, strings.Join(expectedTokens, ","))
		}
	}

	for _, c := range p.nodes[index].conn {
		if p.nodes[c].typ == eof {
			return nil
		}
	}

	return errors.New("command incomplete, unexpected EOF")
}
