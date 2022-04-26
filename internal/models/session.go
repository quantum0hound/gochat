package models

import "time"

/*
CREATE TABLE sessions
(
    id              serial    not null unique,
    user_id         int references users (id) on delete cascade not null,
    refresh_token   varchar(36) not null unique,
    expires_in      timestamp not null,
    fingerprint     varchar(64) not null
)
*/

type Session struct {
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in" db:"expires_in"`
	UserId       int       `json:"user_id" db:"user_id"`
	Fingerprint  string    `json:"fingerprint" db:"fingerprint"`
}
