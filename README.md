# Go-3270

Go-3270 is an [IBM 3270](https://en.wikipedia.org/wiki/IBM_3270) terminal emulator for the modern age. Why on earth is one needed, almost 50 years after the devices were first introduced in 1971? Of course, there is no need at all, I just did this for fun, as a voyage through computer archaeology and a vehicle for learning Go. I don't expect anyone will ever use it as an actual emulator but I hope it showcases how these historically important devices worked.

## Status

Go-3270 is currently under active development. The objective (by no means met yet) is to implement the full specification, as detailed in [GA23-0059-07 3270 Data Stream Programmer's Reference](https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf) 1992 edition.

Consequently, only minimal documentation is presented here.

## Acknowledgements

I am very grateful to [Michael Fogelman](https://github.com/fogleman) and his [Go 3270 App](https://github.com/racingmars/go3270/tree/master) project which have helped me develop controlled and predictable test applications, as well as my old pal [Mark Rafter]() for encouraging me to learn Go in the first place. Michael's [gg](https://github.com/fogleman/gg) Go graphics package is a key dependancy of my small project.

## Development Prerequisites

1. [Bun 1.3.0+](https://bun.com/)
2. [Go 1.25+](https://go.dev/)
3. [Docker 28.5.1+](https://www.docker.com/) (to launch host applications)

## Structure

### Client

Go-3270 presents via a UI written in [Lit](https://lit.dev/) using [Material Web](https://material-web.dev/) components. The UI hosts an HTML `<canvas>` element into which the core Go emulator code displays all the 3270 interactions.

## Go-3270 Emulator

(Almost) all the interactions with the host 3270 application program are handled in Go. The initial telnet "will/do" negotiation is coded in Typescript, as is the Web socket code.

The emulator doesn't directy draw on the `<canvas>` hosted by the UI. Rather, it uses native Go packages to draw into a device context. A `requestAnimationFrame` loop efficiently bitblts the context into the `<canvas>` as it changes.

> Many thanks to [Mark Farnan](https://github.com/markfarnan) and his [go-canvas](https://github.com/markfarnan/go-canvas) project for this idea.

The emulator is modeled on core components, roughly mirroring that of the hardware itself:

1. The `buffer` is an array of cells and fields, matchinmg the discrete positions on the 3270 screen.

2. The `screen` itself which uses the Go [gg](https://github.com/fogleman/gg) package to render into a device context. A `cache` holds pre-rendered `glyphs` for fast rendering.

3. The `keyboard` handles all operator input.

4. The `consumer` ingests outbound (from the host app to nthe 3270) data streams into the `buffer`.

5. The `producer` creates inbound (from the 3270 to the host app) data streams from the `buffer` contents.

All these components communicate via a pubsub `bus`. A separate `logger` component also listens to the `bus` to produce detailed debugging data.

All the Go code is copiled to WASM for execution in the browser. However, the bulk of the code is written in pure Go and doesn't use `syscall/js` in order to make it "Go testable".

Only a single `mediator` uses `syscall/js`. It is responsible for creating a call- and event-based interface with the UI.

## Server

Because the emulator is browser-hosted, _some_ server is necessary to deploy it. The static assets (HTML, CSS, Type/Javascript) could of course be deployed via a CDN, even Github iself.

However, the only kind of socket browser code can use are Web sockets, while communication with the host application must be via TCP sockets. Consequently, the server also acts as a proxy between the emulator and the application. The `Bun.serve` API considerably simplifies this code.

## Development

The `bin` directory contains scripts to facilitate development and testing.

1. `client` launches a process that watches the client code (UI and emulator) and rebuilds it as it changes.

2. `server` does the same for the server code (which in practice rarely changes)

3. `localhost` launches the server on port 3000. It itself watches for client code changes and relaunches itself as necessary.

Typically, all of the above run simultaneously in three separate terminal sessions. A fourth runs one of the following, all of which listen for 3270 TCP connections on port 3270.

4. `mvs` runs [rattydave/docker-ubuntu-hercules-mvs](https://hub.docker.com/r/rattydave/docker-ubuntu-hercules-mvs) to provide MVS as a host application.

5. `vm370` runs [rattydave/docker-ubuntu-hercules-vm370](https://hub.docker.com/r/rattydave/docker-ubuntu-hercules-vm370) to provide VM/370 as a host application.

6. Various other scripts run test applications written in [Go 3270 App](https://github.com/racingmars/go3270/tree/master) as described above.

## Screenshots

### Go 3270 Database App

![Go 3270 Database app](database.png)

### MVS

![MVSp](mvs.png)

### VM/370

![VM/370](vm370.png)
