package pubsub

type Topic int

const (

	// 🟦 Published by mediator on behalf of UI, received by emulator

	CLOSE Topic = iota
	CONFIG
	FOCUS
	KEYSTROKE
	OUTBOUND

	// 🟦 Published by emulator, received by mediator & fowarded via dispatchEvent to UI

	DUMP
)
