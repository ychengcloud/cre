// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package gen

import (
	"go/token"
	"strings"
	"text/template"
	"unicode"

	"github.com/go-openapi/inflect"

	"github.com/ychengcloud/cre/spec"
)

var (
	Funcs = template.FuncMap{
		"receiver": receiver,
		"snake":    snake,
		"pascal":   pascal,
		"camel":    camel,
		"plural":   plural,
		"singular": singular,
	}
	rules    = ruleset()
	acronyms = make(map[string]struct{})
)

func ruleset() *inflect.Ruleset {
	rules := inflect.NewDefaultRuleset()
	for _, w := range []string{
		"ACL", "API", "ASCII", "AWS", "CPU", "CSS", "DNS", "EOF", "GB", "GUID",
		"HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "KB", "LHS", "MAC", "MB",
		"QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "SSO", "TCP",
		"TLS", "TTL", "UDP", "UI", "UID", "URI", "URL", "UTF8", "UUID", "VM",
		"XML", "XMPP", "XSRF", "XSS",
	} {
		rules.AddAcronym(w)
	}
	return rules
}

// receiver returns the receiver name of the given type.
//
//	[]T       => t
//	[1]T      => t
//	User      => u
//	UserQuery => uq
//
func receiver(importPkg []string, s string) (r string) {
	// Trim invalid tokens for identifier prefix.
	s = strings.Trim(s, "[]*&0123456789")
	parts := strings.Split(snake(s), "_")
	min := len(parts[0])
	for _, w := range parts[1:] {
		if len(w) < min {
			min = len(w)
		}
	}

	// 所有可能组合与importPkg中的包名进行比较
	// found为true时，表示与import中的包名冲突
	var found bool
	for i := 1; i < min; i++ {
		found = false
		r := parts[0][:i]
		for _, w := range parts[1:] {
			r += w[:i]
		}

		for _, pkg := range importPkg {
			if pkg == r {
				found = true
				break
			}
		}
		if !found {
			s = r
			break
		}
	}

	name := strings.ToLower(s)
	//如果是保留字，或与 import 包名重复，则加上下划线
	if token.Lookup(name).IsKeyword() || found {
		name = "_" + name
	}
	return name
}

// plural a name.
func plural(name string) string {
	p := rules.Pluralize(name)
	if p == name {
		p += "Slice"
	}
	return p
}

func isSeparator(r rune) bool {
	return r == '_' || r == '-' || unicode.IsSpace(r)
}

func pascalWords(words []string) string {
	for i, w := range words {
		words[i] = rules.Capitalize(w)
	}
	return strings.Join(words, "")
}

// pascal converts the given name into a PascalCase.
//
//	user_info 	=> UserInfo
//	full_name 	=> FullName
//	user_id   	=> UserId
//	full-admin	=> FullAdmin
//
func pascal(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	return pascalWords(words)
}

// camel converts the given name into a camelCase.
//
//	user_info  => userInfo
//	full_name  => fullName
//	user_id    => userId
//	full-admin => fullAdmin
//
func camel(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	if len(words) == 1 {
		return strings.ToLower(words[0])
	}
	return strings.ToLower(words[0]) + pascalWords(words[1:])
}

// snake converts the given struct or field name into a snake_case.
//
//	Username => username
//	FullName => full_name
//	HTTPCode => http_code
//
func snake(s string) string {
	var (
		j int
		b strings.Builder
	)
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		// Put '_' if it is not a start or end of a word, current letter is uppercase,
		// and previous is lowercase (cases like: "UserInfo"), or next letter is also
		// a lowercase and previous letter is not "_".
		if i > 0 && i < len(s)-1 && unicode.IsUpper(r) {
			if unicode.IsLower(rune(s[i-1])) ||
				j != i-1 && unicode.IsLower(rune(s[i+1])) && unicode.IsLetter(rune(s[i-1])) {
				j = i
				b.WriteString("_")
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func singular(s string) string {
	return rules.Singularize(s)
}

func contains(ops []spec.Op, str string) bool {
	for _, op := range ops {
		if op.Name() == str {
			return true
		}
	}

	return false
}
