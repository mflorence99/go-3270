package core

import (
	"emulator/types"
	"emulator/utils"
)

// 🟧 Basic (but type-safe) pubsub implementation

type Bus struct {
	handlers map[Topic][]interface{}
}

//go:generate stringer -type=Topic
type Topic int

// 🔥 run go generate emulator/core if any of these change

const (
	attn Topic = iota
	close
	focus
	keystroke
	inbound
	initialize
	outbound
	panic
	probe
	q
	ql
	rb
	render
	renderDeltas
	reset
	rm
	rma
	status
	tick
	trace
	wcchar
)

// 🟦 Constructor

func NewBus() *Bus {
	b := new(Bus)
	b.handlers = make(map[Topic][]interface{})
	return b
}

// 🟦 Type-safe publishers

func (b *Bus) PubAttn(aid types.AID) {
	b.Publish(attn, aid)
}

func (b *Bus) PubClose() {
	b.Publish(close)
}

func (b *Bus) PubFocus(focussed bool) {
	b.Publish(focus, focussed)
}

func (b *Bus) PubKeystroke(key types.Keystroke) {
	b.Publish(keystroke, key)
}

type PubInboundHints struct{ RB, RM, Short, WSF bool }

func (b *Bus) PubInbound(chars []byte, hints PubInboundHints) {
	b.Publish(inbound, chars, hints)
}

func (b *Bus) PubInitialize() {
	b.Publish(initialize)
}

func (b *Bus) PubOutbound(chars []byte) {
	b.Publish(outbound, chars)
}

func (b *Bus) PubPanic(msg string) {
	b.Publish(panic, msg)
}

func (b *Bus) PubProbe(addr uint) {
	b.Publish(probe, addr)
}

func (b *Bus) PubQ() {
	b.Publish(q)
}

func (b *Bus) PubQL(qcodes []types.QCode) {
	b.Publish(ql, qcodes)
}

func (b *Bus) PubRB(aid types.AID) {
	b.Publish(rb, aid)
}

func (b *Bus) PubRender() {
	b.Publish(render)
}

func (b *Bus) PubRenderDeltas(deltas *utils.Stack[uint]) {
	b.Publish(renderDeltas, deltas)
}

func (b *Bus) PubReset() {
	b.Publish(reset)
}

func (b *Bus) PubRM(aid types.AID) {
	b.Publish(rm, aid)
}

func (b *Bus) PubRMA(aid types.AID) {
	b.Publish(rma, aid)
}

func (b *Bus) PubStatus(stat *types.Status) {
	b.Publish(status, stat)
}

func (b *Bus) PubTick(counter int) {
	b.Publish(tick, counter)
}

func (b *Bus) PubWCChar(wcc types.WCC) {
	b.Publish(wcchar, wcc)
}

// 🟦 Type-safe subscribers

func (b *Bus) SubAttn(fn func(aid types.AID)) {
	b.Subscribe(attn, fn)
}

func (b *Bus) SubClose(fn func()) {
	b.Subscribe(close, fn)
}

func (b *Bus) SubFocus(fn func(focus bool)) {
	b.Subscribe(focus, fn)
}

func (b *Bus) SubKeystroke(fn func(key types.Keystroke)) {
	b.Subscribe(keystroke, fn)
}

func (b *Bus) SubInbound(fn func(chars []byte, hints PubInboundHints)) {
	b.Subscribe(inbound, fn)
}

func (b *Bus) SubInitialize(fn func()) {
	b.Subscribe(initialize, fn)
}

func (b *Bus) SubOutbound(fn func(chars []byte)) {
	b.Subscribe(outbound, fn)
}

func (b *Bus) SubPanic(fn func(msg string)) {
	b.Subscribe(panic, fn)
}

func (b *Bus) SubProbe(fn func(addr uint)) {
	b.Subscribe(probe, fn)
}

func (b *Bus) SubQ(fn func()) {
	b.Subscribe(q, fn)
}

func (b *Bus) SubQL(fn func(qcodes []types.QCode)) {
	b.Subscribe(ql, fn)
}

func (b *Bus) SubRB(fn func(aid types.AID)) {
	b.Subscribe(rb, fn)
}

func (b *Bus) SubRender(fn func()) {
	b.Subscribe(render, fn)
}

func (b *Bus) SubRenderDeltas(fn func(deltas *utils.Stack[uint])) {
	b.Subscribe(renderDeltas, fn)
}

func (b *Bus) SubReset(fn func()) {
	b.Subscribe(reset, fn)
}

func (b *Bus) SubRM(fn func(aid types.AID)) {
	b.Subscribe(rm, fn)
}

func (b *Bus) SubRMA(fn func(aid types.AID)) {
	b.Subscribe(rma, fn)
}

func (b *Bus) SubStatus(fn func(stat *types.Status)) {
	b.Subscribe(status, fn)
}

func (b *Bus) SubTick(fn func(counter int)) {
	b.Subscribe(tick, fn)
}

func (b *Bus) SubWCChar(fn func(wcc types.WCC)) {
	b.Subscribe(wcchar, fn)
}

// 🔥 Debug only

func (b *Bus) SubTrace(fn func(topic Topic, handler interface{})) {
	b.Subscribe(trace, fn)
}

// 🟦 Brute force cleanup

func (b *Bus) UnsubscribeAll() {
	b.handlers = make(map[Topic][]interface{})
}

// 🟦 Public functions (public just for test cases)

func (b *Bus) Publish(topic Topic, args ...any) {
	debuggers := b.handlers[trace]
	handlers := b.handlers[topic]
	for _, handler := range handlers {
		// 🔥 trace BEFORE action
		for _, debugger := range debuggers {
			utils.Call(debugger, topic, handler)
		}
		utils.Call(handler, args...)
	}
}

func (b *Bus) Subscribe(topic Topic, fn interface{}) {
	b.handlers[topic] = append(b.handlers[topic], fn)
}
