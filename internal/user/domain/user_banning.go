package domain

import (
	"errors"
	"time"
)

type banning struct {
	isBanned       bool
	reason         string
	isInDefinitely bool
	from           time.Time
	to             time.Time
}

type banTimeline struct {
	from time.Time
	to   time.Time
}

func NewBanTimeline(from, to time.Time) *banTimeline {
	return &banTimeline{from: from, to: to}
}

// Keep error messages here so we can import and reuse them for unit testing
var (
	ErrEmptyReason         = errors.New("reason to ban a user can't be empty")
	ErrUserAlreadyBanned   = errors.New("user is already banned")
	ErrBanTimelineRequired = errors.New("you must pass correct ban timeline if user ban is not indefinitely")
	ErrUserIsNotBanned     = errors.New("user you are trying to unban is not banned")
)

func (u *User) Ban(reason string, isInDefinitely bool, timeline *banTimeline) error {
	if reason == "" {
		return ErrEmptyReason
	}
	if u.banStatus.isBanned {
		return ErrUserAlreadyBanned
	}
	if !isInDefinitely && timeline == nil {
		return ErrBanTimelineRequired
	}
	u.banStatus.isBanned = true
	u.banStatus.reason = reason
	u.banStatus.isInDefinitely = isInDefinitely
	if timeline != nil {
		u.banStatus.from = timeline.from
		u.banStatus.to = timeline.to
	}

	return nil
}

func (u *User) UnBan() error {
	if !u.banStatus.isBanned {
		return ErrUserIsNotBanned
	}

	u.banStatus.isBanned = false

	return nil
}

func (u *User) IsBanned() bool {
	return u.banStatus.isBanned
}

func (u *User) ReasonForBan() string {
	return u.banStatus.reason
}

func (u *User) IsBanIndefinitely() bool {
	return u.banStatus.isInDefinitely
}

func (u *User) BanStartDate() time.Time {
	return u.banStatus.from
}

func (u *User) BanEndDate() time.Time {
	return u.banStatus.to
}
