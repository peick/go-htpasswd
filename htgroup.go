// Package htpasswd groups provides an autorisation mechanism using Apache-style group files.
//
// An Apache group file looks like this:
// users: user1 user2 user3
// admins: user1
//
// Basic usage of this package:
//
// userGroups, groupLoadErr := htgroup.NewHTGroup("./my-group-file", nil)
// ok := userGroups.IsUserInGroup(username, "admins")
package htpasswd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
)

// Data structure for users and theirs groups (map).
// The map key is the user, the value is an array of groups.
type userGroupMap map[string][]string

// A HTGroup encompasses an Apache-style group file.
type HTGroup struct {
	filePath   string
	userGroups atomic.Pointer[userGroupMap]
}

// NewHTGroup creates a HTGroup from an Apache-style group file.
//
// The filename must exist and be accessible to the process, as well as being a valid group file.
//
// bad is a function, which if not nil will be called for each malformed or rejected entry in the group file.
func NewHTGroup(filename string) (*HTGroup, error) {
	htGroup := HTGroup{
		filePath: filename,
	}
	return &htGroup, htGroup.Reload()
}

// NewHTGroupsFromReader is like NewHTGroup but reads from r instead of a named file.
func NewHTGroupsFromReader(r io.Reader) (*HTGroup, error) {
	htGroup := HTGroup{}

	readFileErr := htGroup.ReloadFromReader(r)
	if readFileErr != nil {
		return nil, readFileErr
	}

	return &htGroup, nil
}

// Reload rereads the group file.
func (g *HTGroup) Reload() error {
	file, err := os.Open(g.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return g.ReloadFromReader(file)
}

// ReloadFromReader rereads the group file from a Reader.
func (g *HTGroup) ReloadFromReader(r io.Reader) error {
	userGroups := make(userGroupMap)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if err := processLine(&userGroups, line); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanning group file failed: %w", err)
	}

	g.userGroups.Store(&userGroups)

	return nil
}

func processLine(userGroups *userGroupMap, rawLine string) error {
	// ignore empty line
	line := strings.TrimSpace(rawLine)
	if line == "" {
		return nil
	}

	// ignore comment line. Inline comments are not allowed
	if strings.HasPrefix(line, "#") {
		return nil
	}

	groupAndUsers := strings.SplitN(line, ":", 2)
	if len(groupAndUsers) != 2 {
		return fmt.Errorf("malformed line, no colon: %s", line)
	}

	var group = strings.TrimSpace(groupAndUsers[0])
	var users = strings.Fields(groupAndUsers[1])
	for _, user := range users {
		if (*userGroups)[user] == nil {
			(*userGroups)[user] = []string{}
		}
		(*userGroups)[user] = append((*userGroups)[user], group)
	}

	return nil
}

// IsUserInGroup checks whether the user is in a group.
// Returns true of user is in that group, otherwise false.
func (g *HTGroup) IsUserInGroup(user string, group string) bool {
	groups := g.GetUserGroups(user)
	return containsGroup(groups, group)
}

// GetUserGroups reads all groups of a user.
// Returns all groups as a string array or an empty array.
func (g *HTGroup) GetUserGroups(user string) []string {
	groups := (*g.userGroups.Load())[user]

	if groups == nil {
		return []string{}
	}
	return groups
}

func containsGroup(groups []string, group string) bool {
	for _, g := range groups {
		if g == group {
			return true
		}
	}
	return false
}
