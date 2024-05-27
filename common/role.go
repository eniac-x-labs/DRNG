package common

import (
	"strings"
	"sync/atomic"
)

type RoleType uint32

const (
	LEADER RoleType = iota
	GENERATOR
	UNKNOWN
)

const (
	LEADERSTR    = "leader"
	GENERATORSTR = "generator"
	UNKNOWNSTR   = "unknown"
)

func (r RoleType) String() string {
	switch r {
	case LEADER:
		return LEADERSTR
	case GENERATOR:
		return GENERATORSTR
	default:
		return UNKNOWNSTR
	}
}

func String2RoleType(s string) RoleType {
	switch strings.ToLower(s) {
	case LEADERSTR:
		return LEADER
	case GENERATORSTR:
		return GENERATOR
	default:
		return UNKNOWN
	}
}

func (r *RoleType) GetRole() RoleType {
	return RoleType(atomic.LoadUint32((*uint32)(r)))
}

func (r *RoleType) SetRole(role RoleType) {
	atomic.StoreUint32((*uint32)(r), uint32(role))
}
