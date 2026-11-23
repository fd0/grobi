//go:build !linux && !freebsd && !openbsd

package main

func readLidState() LidState { return LidUnknown }

func subscribeLidEvents(_ chan<- struct{}, _ <-chan struct{}) {}
