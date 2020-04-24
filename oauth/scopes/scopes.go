package scopes

import (
	"fmt"
	"strings"
)

// Scope represents a set of requested scopes
// https://forge.autodesk.com/en/docs/oauth/v2/developers_guide/scopes/
type Scope uint64

const (
	// UserProfileRead allow the app to view the user's profile
	UserProfileRead = Scope(1 << iota)

	// UserRead allows the app to view the user's profile
	UserRead

	// UserWrite allow the app to write to the user's profile
	UserWrite

	// ViewablesRead View your viewable data
	ViewablesRead

	// DataRead view your data
	DataRead

	// DataWrite manage your data
	DataWrite

	// DataCreate write your data
	DataCreate

	// DataSearch search across your data
	DataSearch

	// BucketCreate creates new buckets
	BucketCreate

	// BucketRead view your buckets
	BucketRead

	// BucketUpdate update your buckets
	BucketUpdate

	// BucketDelete delete your buckets
	BucketDelete

	// CodeAll author or execute your code
	CodeAll

	// AccountRead view your product and service account
	AccountRead

	// AccountWrite manage your product and service accounts
	AccountWrite

	scopeEnd // should be last element
)

var scopes = [...]string{
	"user-profile:read",
	"user:read",
	"user:write",
	"viewables:read",
	"data:read",
	"data:write",
	"data:create",
	"data:search",
	"bucket:create",
	"bucket:read",
	"bucket:update",
	"bucket:delete",
	"code:all",
	"account:read",
	"account:write",
}

var invalidMask Scope

func init() {
	if 1<<len(scopes) != scopeEnd {
		panic(
			fmt.Sprintf(
				"More Scopes(%v) defined then in scopes array(%v)",
				uint64(scopeEnd),
				1<<len(scopes),
			),
		)
	}
	// setup invalidMask
	invalidMask = ^invalidMask
	for i := 0; i < len(scopes); i++ {
		invalidMask &= ^Scope(1 << i)
	}
}

// For takes a space seprated list of scopes and returns the Scope
// set for it.
func For(val string) Scope {
	var scope Scope
	scps := strings.Split(" ", strings.ToLower(val))
	if len(scps) == 0 {
		return scope
	}
	for i := range scps {
		s := strings.TrimSpace(scps[i])
		for j, scp := range scopes {
			if scp != s {
				continue
			}
			scope |= Scope(1 << j)
			break
		}
	}
	return scope

}

// Allows checks to see if s allows for all scopes in s1
func (s Scope) Allows(s1 Scope) bool { return s&s1 == s1 }

// IsValid checks to see if at least one known scope is encoded
func (s Scope) IsValid() bool { return s != 0 && s&invalidMask == 0 }

// String will return the string version of the set of scopes
func (s Scope) String() string {
	if s == 0 {
		// no scopes defined
		return ""
	}
	var str strings.Builder
	var first = true
	for i := 0; i < len(scopes); i++ {
		if ((1 << i) & s) == 0 {
			continue
		}
		if !first {
			str.WriteRune(' ')
		}
		str.WriteString(scopes[i])
		first = false
	}
	return str.String()
}
