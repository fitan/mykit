package mygorm

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Query struct {
	v url.Values
}

func NewQuery() *Query {
	return &Query{v: url.Values{}}
}

func (q *Query) Q(field string, op string, values ...string) *Query {
	q.v.Add("_q", fmt.Sprintf("%s%s%s", field, op, strings.Join(values, ",")))
	return q
}

func (q *Query) Paging(page string, pageSize string) *Query {
	q.v.Add("_page", page)
	q.v.Add("_pageSize", pageSize)
	return q
}

func (q *Query) Select(fields ...string) *Query {
	q.v.Add("_select", strings.Join(fields, ","))
	return q
}

func (q *Query) Sort(fields ...string) *Query {
	for _, v := range fields {
		q.v.Add("_sort", v)
	}
	return q
}

func (q *Query) Include(fields ...string) *Query {
	q.v.Add("_include", strings.Join(fields, ","))
	return q
}

func (q *Query) Encode() string {
	return q.v.Encode()
}

func (q *Query) NewRequest() *http.Request {
	r, _ := http.NewRequest("GET", "http://localhost:8080?"+q.Encode(), nil)
	return r
}
