import { Config } from '$client/state/state';
import { Dimensions } from '$client/state/consts';
import { Emulators } from '$client/state/consts';
import { LitElement } from 'lit';
import { MdDialog } from '@material/web/dialog/dialog.js';
import { SignalWatcher } from '@lit-labs/signals';
import { State } from '$client/state/state';
import { TemplateResult } from 'lit';
import { Tn3270 } from '$client/services/tn3270';

import { consume } from '@lit/context';
import { css } from 'lit';
import { customElement } from 'lit/decorators.js';
import { globals } from '$client/css/globals/shadow-dom';
import { html } from 'lit';
import { query } from 'lit/decorators.js';
import { repeat } from 'lit/directives/repeat.js';
import { state } from 'lit/decorators.js';
import { stateContext } from '$client/state/state';
import { styleMap } from 'lit/directives/style-map.js';

declare global {
  interface HTMLElementTagNameMap {
    'app-connector': Connector;
  }
}

// 📘 manage 3270 connection

@customElement('app-connector')
export class Connector extends SignalWatcher(LitElement) {
  static override styles = [
    globals,
    css`
      .stretcher {
        align-items: center;
        display: flex;
        flex-direction: column;
        height: 100%;
        justify-content: center;
        min-width: 800px /* 👈 to prevent wrapping */;

        .connector {
          display: flex;
          flex-direction: column;
          gap: 2rem;
          justify-content: stretch;

          .header {
            text-align: center;

            .major {
              color: #00cc00;
              font-family: Terminal;
              font-size: 12rem;
              font-weight: bold;
              line-height: 1;
              text-shadow: 0.33rem 0.33rem #008000;
            }

            .minor {
              font-size: 2rem;
              font-weight: bold;
              text-transform: uppercase;
            }
          }

          .config {
            align-items: stretch;
            display: flex;
            flex-direction: row;
            gap: 1rem;
            justify-content: center;

            .color,
            .connection,
            .emulation {
              display: flex;
              flex-direction: column;
              gap: 0.5rem;
              justify-content: space-between;

              .instructions {
                font-weight: bold;
              }
            }

            .color .sample {
              font-family: Terminal;
            }

            .connection .host {
              display: flex;
              flex-direction: row;
              gap: 0.5rem;
            }

            .emulation .dims {
              font-size: smaller;
            }

            .label {
              align-items: center;
              display: flex;
              flex-direction: row;
              gap: 0.25rem;
            }
          }
        }
      }
    `
  ];

  @state() connecting!: boolean;
  @query('#dialog') dialog!: MdDialog;
  @state() message!: string;
  @consume({ context: stateContext }) state!: State;

  #tn3270: Tn3270 | null = null;

  // 👁️ https://dev.to/blikblum/dry-form-handling-with-lit-19f
  // 👇 "connected" here means socket connection
  async connect(evt: Event): Promise<void> {
    evt.preventDefault();
    const form = evt.target as HTMLFormElement;
    if (form) {
      const formData = new FormData(form);
      const config = Object.fromEntries(formData.entries()) as Config;
      this.state.updateConfig(config);
      try {
        // 👇 try to connect to 3270
        this.connecting = true;
        this.#tn3270?.close();
        this.#tn3270 = await Tn3270.tn3270(
          config.host,
          config.port,
          Emulators[config.emulator] as string
        );
        this.#tn3270.stream$.subscribe({
          next: (bytes: Uint8ClampedArray) => {
            if (this.connecting)
              this.dispatchEvent(new CustomEvent('connected'));
            this.dispatchEvent(
              new CustomEvent('receiveFromApp', {
                detail: { bytes }
              })
            );
            this.connecting = false;
          },

          // 🔥 WebSocket connection established, but that failed
          error: async (e: any) => {
            console.error(e);
            this.connecting = false;
            this.message = e.reason;
            await this.dialog.show();
            this.dispatchEvent(new CustomEvent('disconnected'));
            this.#tn3270 = null;
          },

          // 👇 normal completion eg: Tn3270.close()
          complete: () => {
            console.log(
              `%c3270 -> Server -> Client %cClosed`,
              'color: palegreen',
              'color: cyan'
            );
            this.dispatchEvent(new CustomEvent('disconnected'));
            this.#tn3270 = null;
          }
        });
      } catch (e: any) {
        // 🔥 tried to upgrade to WebSocket, but that failed
        console.error(e);
        this.connecting = false;
        this.message = `Unable to reach proxy server ${location.hostname}:${location.port}`;
        await this.dialog.show();
        this.dispatchEvent(new CustomEvent('disconnected'));
        this.#tn3270 = null;
      }
    }
  }

  disconnect(): void {
    this.#tn3270?.close();
  }

  palette(): void {
    this.dispatchEvent(new CustomEvent('palette'));
  }

  panic(message: string): void {
    this.message = message;
    this.dialog.show();
    this.#tn3270?.close();
  }

  override render(): TemplateResult {
    return html`
      <main class="stretcher">
        <section class="connector">
          <article class="header">
            <header class="major">3270</header>
            <p class="minor">Go-powered 3270 Emulator</p>
          </article>

          <hr />

          <form @submit=${this.connect} class="config" name="config">
            <article class="connection">
              <div class="host">
                <md-filled-text-field
                  label="Hostname or IP"
                  name="host"
                  style="width: 10rem"
                  value=${this.state.model.get().config
                    .host}></md-filled-text-field>
                <md-filled-text-field
                  label="Port"
                  name="port"
                  style="width: 5rem"
                  value=${this.state.model.get().config
                    .port}></md-filled-text-field>
              </div>
            </article>

            <article class="emulation">
              <p class="instructions">Select 3270 Model to Emulate</p>

              ${repeat(
                Object.entries(Emulators),
                (emulator) => emulator[0],
                (emulator) => html`
                  <label class="label">
                    <md-radio
                      ?checked=${this.state.model.get().config
                        .emulator === emulator[0]}
                      name="emulator"
                      value=${emulator[0]}></md-radio>
                    ${emulator[1]} &mdash;
                    <em class="dims">
                      ${Dimensions[emulator[0]]?.[0]} x
                      ${Dimensions[emulator[0]]?.[1]}
                    </em>
                  </label>
                `
              )}

              <br />

              <md-filled-button
                ?disabled=${this.connecting}
                style="align-self: center">
                ${this.connecting ? 'Connnecting...' : 'Connect'}
                <app-icon
                  icon="hourglass_full"
                  style=${styleMap({
                    display: this.connecting ? 'block' : 'none'
                  })}
                  slot="icon"></app-icon>
              </md-filled-button>
            </article>

            <article class="color">
              <p class="instructions">Select Default 3270 Color</p>

              ${repeat(
                ['green', 'blue', 'yellow', 'white'],
                (color) => color,
                (color) => html`
                  <label class="label">
                    <md-radio
                      ?checked=${this.state.model.get().config.color ===
                      color}
                      name="color"
                      value=${color}></md-radio>
                    <span
                      class="sample"
                      style=${styleMap({
                        color: this.state.model.get().clut[color]![0]
                      })}>
                      CUSTOMER NUM:
                      <b>123456</b>
                    </span>
                  </label>
                `
              )}

              <br />

              <md-outlined-button
                @click=${this.palette}
                style="align-self: center"
                type="button">
                Customize
              </md-outlined-button>
            </article>
          </form>
        </section>
      </main>

      <md-dialog id="dialog">
        <header slot="headline">3270 Connection Error</header>
        <section slot="content">
          <p>
            An error occured while connecting to the
            ${Emulators[this.state.model.get().config.emulator]} at
            ${this.state.model.get().config
              .host}:${this.state.model.get().config.port}.
            Please take any necessary corrective action and retry.
            <br />
            <br />
            ${this.message}
          </p>
        </section>
        <form slot="actions" method="dialog">
          <md-outlined-button>OK</md-outlined-button>
        </form>
      </md-dialog>
    `;
  }

  sendToApp(bytes: Uint8ClampedArray): void {
    this.#tn3270?.sendToApp(bytes);
  }
}
