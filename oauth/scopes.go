package oauth

import (
	"fmt"
	"strings"
)

// Scopes represents a set of requested scopes
// https://forge.autodesk.com/en/docs/oauth/v2/developers_guide/scopes/
type Scopes uint64

const (
	// ScopeUserProfileRead allow the app to view the user's profile
	ScopeUserProfileRead = Scopes(1 << iota)
	// ScopeUserRead allows the app to view the user's profile
	ScopeUserRead
	// ScopeUserWrite allow the app to write to the user's profile
	ScopeUserWrite
	// ScopeViewablesRead View your viewable data
	ScopeViewablesRead
	// ScopeDataRead view your data
	ScopeDataRead
	// ScopeDataWrite manage your data
	ScopeDataWrite

	//ScopeDataCreate write your data
	ScopeDataCreate

	//ScopeDataSearch search across your data
	ScopeDataSearch

	//ScopeBucketCreate creates new buckets
	ScopeBucketCreate

	//ScopeBucketRead view your buckets
	ScopeBucketRead

	//ScopeBucketUpdate update your buckets
	ScopeBucketUpdate

	//ScopeBucketDelete delete your buckets
	ScopeBucketDelete

	// ScopeCodeAll author or execute your code
	ScopeCodeAll

	//ScopeAccountRead view your product and service account
	ScopeAccountRead

	//ScopeAccountWrite manage your product and service accounts
	ScopeAccountWrite

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

var invalidMask Scopes

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
		invalidMask &= ^Scopes(1 << i)
	}
}

//ScopeFor takes a space seprated list of scopes and returns the Scope
// set for it.
func ScopeFor(val string) Scopes {
	var scope Scopes
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
			scope |= Scopes(1 << j)
			break
		}
	}
	return scope

}

// Allows checks to see if s allows for all scopes in s1
func (s Scopes) Allows(s1 Scopes) bool { return s&s1 == s1 }

// IsValid checks to see if at least one known scope is encoded
func (s Scopes) IsValid() bool { return s != 0 && s&invalidMask == 0 }

// String will return the string version of the set of scopes
func (s Scopes) String() string {
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
