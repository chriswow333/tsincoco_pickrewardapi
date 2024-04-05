package model

type Status int32

var AllStatuses = []Status{Active, Pilot, Inactive}

const (
	Inactive Status = iota
	Pilot
	Active
)
