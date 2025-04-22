package domain

import (
	"errors"
	"strings"
	"time"
)

// Keep error messages here so we can import and reuse them for unit testing
var (
	ErrEmptyReason         = errors.New("reason to ban a user can't be empty")
	ErrUserAlreadyBanned   = errors.New("user is already banned")
	ErrBanTimelineRequired = errors.New("you must pass correct ban timeline if user ban is not indefinitely")
	ErrUserIsNotBanned     = errors.New("user you are trying to unban is not banned")
)

type ban struct {
	isBanned     bool
	reason       string
	isIndefinite bool
	from         time.Time
	to           time.Time
	bannedAt     time.Time
}

type BanTimeline struct {
	from time.Time
	to   time.Time
}

func NewBanTimeline(from, to time.Time) *BanTimeline {
	return &BanTimeline{from: from, to: to}
}

func NewBan(isBanned bool, reason string, isIndefinite bool, from, to, bannedAt time.Time) *ban {
	return &ban{
		isBanned:     isBanned,
		reason:       reason,
		isIndefinite: isIndefinite,
		from:         from,
		to:           to,
		bannedAt:     bannedAt,
	}
}

func (u *User) Ban(reason string, isIndefinite bool, timeline *BanTimeline) error {
	if strings.TrimSpace(reason) == "" {
		return ErrEmptyReason
	}
	if u.banStatus.isBanned {
		return ErrUserAlreadyBanned
	}
	if !isIndefinite && timeline == nil {
		return ErrBanTimelineRequired
	}
	u.banStatus.isBanned = true
	u.banStatus.reason = reason
	u.banStatus.isIndefinite = isIndefinite
	u.banStatus.bannedAt = time.Now()
	if timeline != nil {
		u.banStatus.from = timeline.from
		u.banStatus.to = timeline.to
		u.banStatus.isIndefinite = false
	}

	return nil
}

func (u *User) UnBan() error {
	if !u.banStatus.isBanned {
		return ErrUserIsNotBanned
	}

	u.banStatus.isBanned = false
	u.banStatus.reason = ""
	u.banStatus.isIndefinite = false
	u.banStatus.from = time.Time{}
	u.banStatus.to = time.Time{}

	return nil
}

func (u *User) IsBanned() bool {
	return u.banStatus.isBanned
}

func (u *User) ReasonForBan() string {
	return u.banStatus.reason
}

func (u *User) IsBanIndefinite() bool {
	return u.banStatus.isIndefinite
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
