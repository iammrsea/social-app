package domain

import (
	"errors"
	"strings"
	"time"
)

type banning struct {
	isBanned       bool
	reason         string
	isInDefinitely bool
	from           time.Time
	to             time.Time
	bannedAt       time.Time
}

type BanTimeline struct {
	from time.Time
	to   time.Time
}

func NewBanTimeline(from, to time.Time) *BanTimeline {
	return &BanTimeline{from: from, to: to}
}

// Keep error messages here so we can import and reuse them for unit testing
var (
	ErrEmptyReason         = errors.New("reason to ban a user can't be empty")
	ErrUserAlreadyBanned   = errors.New("user is already banned")
	ErrBanTimelineRequired = errors.New("you must pass correct ban timeline if user ban is not indefinitely")
	ErrUserIsNotBanned     = errors.New("user you are trying to unban is not banned")
)

func (u *User) Ban(reason string, isInDefinitely bool, timeline *BanTimeline) error {
	if strings.TrimSpace(reason) == "" {
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
	u.banStatus.bannedAt = time.Now()
	if timeline != nil {
		u.banStatus.from = timeline.from
		u.banStatus.to = timeline.to
		u.banStatus.isInDefinitely = false
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

func (u *User) BannedAt() time.Time {
	return u.banStatus.bannedAt
}
