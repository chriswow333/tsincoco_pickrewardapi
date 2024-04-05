package model

type Operation int32

const (
	INSERT Operation = iota
	UPDATE
	DELETE
)
