package xorm

import (
	"fmt"
)

const (
	defaultLimit  = 10000
	defaultOffset = 0
)

// Condition defines condition elements.
type Condition struct {
	Key   string      `json:"key"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

// Conditions contains a list of condition and it raw query string.
type Conditions struct {
	Conds  []Condition `json:"conditions"`
	Offset int         `json:"offset,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	args   []interface{}
	query  string
}

// NewConditions return a new initialized conditions.
func NewConditions() *Conditions {
	return &Conditions{
		Conds:  make([]Condition, 0),
		Offset: defaultOffset,
		Limit:  defaultLimit,
	}
}

// Parse parses conditions args and raw query.
func (c Conditions) Parse() (string, []interface{}) {
	c.args = make([]interface{}, 0)
	for _, cond := range c.Conds {
		c.query += fmt.Sprintf("%s %s ? and ", cond.Key, cond.Op)
		c.args = append(c.args, cond.Value)
	}

	if c.query != "" {
		c.query = c.query[:len(c.query)-5]
	}

	return c.query, c.args
}
