package common

import (
	"errors"

	"github.com/speps/go-hashids"
)

var h *hashids.HashID

func InitHashID(salt string) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 8
	hs, _ := hashids.NewWithData(hd)

	h = hs
}

// ID2UID ...
func ID2UID(id uint) string {
	e, _ := h.EncodeInt64([]int64{int64(id)})
	return e
}

// UID2ID ...
func UID2ID(uid string) (int, error) {
	d, err := h.DecodeWithError(uid)
	if err != nil {
		return 0, err
	}

	if len(d) == 1 {
		return int(d[0]), nil
	}

	return 0, errors.New("uid parse fail")
}

// EncPidAndGid ...
func EncPidAndGid(pid, gid uint) string {
	e, _ := h.EncodeInt64([]int64{int64(pid), int64(gid)})
	return e
}

// DecPidAndGid ...
func DecPidAndGid(uid string) (uint, uint, error) {
	d, err := h.DecodeWithError(uid)
	if err != nil {
		return 0, 0, err
	}

	if len(d) == 2 {
		return uint(d[0]), uint(d[1]), nil
	} else if len(d) == 1 {
		return uint(d[0]), 0, nil
	}

	return 0, 0, errors.New("uid parse fail")
}

func EncChatShareInfo(roomID int, expireTime, shareTime int64) string {
	numbers := []int64{int64(roomID), expireTime, shareTime}
	e, _ := h.EncodeInt64(numbers)
	return e
}

func DecChatShareInfo(uid string) (int, int64, int64, error) {
	d, err := h.DecodeWithError(uid)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(d) == 3 {
		return d[0], int64(d[1]), int64(d[2]), nil
	}

	return 0, 0, 0, errors.New("uid parse fail")
}
