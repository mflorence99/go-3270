package pubsub

import (
	"go3270/emulator/consts"
	"go3270/emulator/utils"
	"go3270/emulator/wcc"
)

// 🟧 Basic (but type-safe) pubsub implementation

type Bus struct {
	handlers map[string][]interface{}
}

// 🟦 Constructor

func NewBus() *Bus {
	b := new(Bus)
	b.handlers = make(map[string][]interface{})
	return b
}

// 🟦 Type-safe publishers

func (b *Bus) PubAttn(aid consts.AID) {
	b.Publish("attn", aid)
}

func (b *Bus) PubClose() {
	b.Publish("close")
}

func (b *Bus) PubConfig(cfg Config) {
	b.Publish("config", cfg)
}

func (b *Bus) PubFocus(focus bool) {
	b.Publish("focus", focus)
}

func (b *Bus) PubKeystroke(key Keystroke) {
	b.Publish("keystroke", key)
}

func (b *Bus) PubInbound(chars []byte, wsf bool) {
	b.Publish("inbound", chars, wsf)
}

func (b *Bus) PubOutbound(chars []byte) {
	b.Publish("outbound", chars)
}

func (b *Bus) PubPanic(msg string) {
	b.Publish("panic", msg)
}

func (b *Bus) PubProbe(addr int) {
	b.Publish("probe", addr)
}

func (b *Bus) PubQ() {
	b.Publish("q")
}

func (b *Bus) PubQL(qcodes []consts.QCode) {
	b.Publish("ql", qcodes)
}

func (b *Bus) PubRB(aid consts.AID) {
	b.Publish("rb", aid)
}

func (b *Bus) PubRender() {
	b.Publish("render")
}

func (b *Bus) PubRenderDeltas(deltas *utils.Stack[int]) {
	b.Publish("render-deltas", deltas)
}

func (b *Bus) PubReset() {
	b.Publish("reset")
}

func (b *Bus) PubRM(aid consts.AID) {
	b.Publish("rm", aid)
}

func (b *Bus) PubRMA(aid consts.AID) {
	b.Publish("rms", aid)
}

func (b *Bus) PubStatus(stat *Status) {
	b.Publish("status", stat)
}

func (b *Bus) PubTick(counter int) {
	b.Publish("tick", counter)
}

func (b *Bus) PubWCC(wcc wcc.WCC) {
	b.Publish("wcc", wcc)
}

// 🟦 Type-safe subscribers

func (b *Bus) SubAttn(fn func(aid consts.AID)) {
	b.Subscribe("attn", fn)
}

func (b *Bus) SubClose(fn func()) {
	b.Subscribe("close", fn)
}

func (b *Bus) SubConfig(fn func(cfg Config)) {
	b.Subscribe("config", fn)
}

func (b *Bus) SubFocus(fn func(focus bool)) {
	b.Subscribe("focus", fn)
}

func (b *Bus) SubKeystroke(fn func(key Keystroke)) {
	b.Subscribe("keystroke", fn)
}

func (b *Bus) SubInbound(fn func(chars []byte, wsh bool)) {
	b.Subscribe("inbound", fn)
}

func (b *Bus) SubOutbound(fn func(chars []byte)) {
	b.Subscribe("outbound", fn)
}

func (b *Bus) SubPanic(fn func(msg string)) {
	b.Subscribe("panic", fn)
}

func (b *Bus) SubProbe(fn func(addr int)) {
	b.Subscribe("probe", fn)
}

func (b *Bus) SubQ(fn func()) {
	b.Subscribe("q", fn)
}

func (b *Bus) SubQL(fn func(qcodes []consts.QCode)) {
	b.Subscribe("ql", fn)
}

func (b *Bus) SubRB(fn func(aid consts.AID)) {
	b.Subscribe("rb", fn)
}

func (b *Bus) SubRender(fn func()) {
	b.Subscribe("render", fn)
}

func (b *Bus) SubRenderDeltas(fn func(deltas *utils.Stack[int])) {
	b.Subscribe("render-deltas", fn)
}

func (b *Bus) SubReset(fn func()) {
	b.Subscribe("reset", fn)
}

func (b *Bus) SubRM(fn func(aid consts.AID)) {
	b.Subscribe("rm", fn)
}

func (b *Bus) SubRMA(fn func(aid consts.AID)) {
	b.Subscribe("rma", fn)
}

func (b *Bus) SubStatus(fn func(stat *Status)) {
	b.Subscribe("status", fn)
}

func (b *Bus) SubTick(fn func(counter int)) {
	b.Subscribe("tick", fn)
}

func (b *Bus) SubWCC(fn func(wcc wcc.WCC)) {
	b.Subscribe("wcc", fn)
}

// 🔥 Debug only

func (b *Bus) SubTrace(fn func(topic string, handler interface{})) {
	b.Subscribe("$$$", fn)
}

// 🟦 Brute force cleanup

func (b *Bus) UnsubscribeAll() {
	b.handlers = make(map[string][]interface{})
}

// 🟦 Public functions (public just for test cases)

func (b *Bus) Publish(topic string, args ...any) {
	debuggers := b.handlers["$$$"]
	handlers := b.handlers[topic]
	for _, handler := range handlers {
		utils.Call(handler, args...)
		for _, debugger := range debuggers {
			utils.Call(debugger, topic, handler)
		}
	}
}

func (b *Bus) Subscribe(topic string, fn interface{}) {
	b.handlers[topic] = append(b.handlers[topic], fn)
}
