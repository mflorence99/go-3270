import { Colors } from '$lib/types/3270';
import { Config } from '$client/state/state';
import { Dimensions } from '$lib/types/3270';
import { Emulators } from '$lib/types/3270';
import { LitElement } from 'lit';
import { MdDialog } from '@material/web/dialog/dialog.js';
import { Pages } from '$client/pages/root';
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

// 📘 close 3270 connection on exit

window.addEventListener('beforeunload', () => {
  State.theTn3270?.close();
});

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

              .instructions {
                font-weight: bold;
              }
            }

            .color .sample {
              font-family: Terminal;
              font-size: larger;
            }

            .connection {
              align-items: center;
              justify-content: space-between;
            }

            .connection .host {
              display: flex;
              flex-direction: row;
              gap: 0.5rem;
            }

            .emulation .dims {
              font-size: smaller;
            }
          }
        }
      }
    `
  ];

  @state() connecting!: boolean;
  @query('#dialog') dialog!: MdDialog;
  @state() message!: string;
  @consume({ context: stateContext }) theState!: State;

  override render(): TemplateResult {
    const model = this.theState.model;
    return html`
      <main class="stretcher">
        <section class="connector">
          <article class="header">
            <header class="major">3270</header>
            <p class="minor">Go-powered 3270 Emulator</p>
          </article>

          <hr />

          <form @submit=${this.#connect} class="config" name="config">
            <article class="connection">
              <div class="host">
                <md-filled-text-field
                  label="Hostname or IP"
                  name="host"
                  style="width: 10rem"
                  value=${model.get().config
                    .host}></md-filled-text-field>
                <md-filled-text-field
                  label="Port"
                  name="port"
                  style="width: 5rem"
                  value=${model.get().config
                    .port}></md-filled-text-field>
              </div>

              <md-filled-button ?disabled=${this.connecting}>
                ${this.connecting ? 'Connnecting...' : 'Connect'}
                <app-icon
                  icon="hourglass_full"
                  style=${styleMap({
                    display: this.connecting ? 'block' : 'none'
                  })}
                  slot="icon"></app-icon>
              </md-filled-button>
            </article>

            <article class="emulation">
              <p class="instructions">Select 3270 Model to Emulate</p>

              ${repeat(
                Object.entries(Emulators),
                (emulator) => emulator[0],
                (emulator) => html`
                  <label>
                    <md-radio
                      ?checked=${model.get().config.emulator ===
                      emulator[0]}
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
            </article>

            <article class="color">
              <p class="instructions">Select Default 3270 Color</p>

              ${repeat(
                Object.entries(Colors),
                (color) => color[0],
                (color) => html`
                  <label>
                    <md-radio
                      ?checked=${model.get().config.color === color[0]}
                      name="color"
                      value=${color[0]}></md-radio>
                    <span
                      class="sample"
                      style=${styleMap({
                        color: color[1]
                      })}>
                      CUSTOMER NUM: 123456
                    </span>
                  </label>
                `
              )}
            </article>
          </form>
        </section>
      </main>

      <md-dialog id="dialog">
        <header slot="headline">3270 Connection Error</header>
        <section slot="content">
          <p>
            An error occured while connecting to the
            ${model.get().config.emulator} at
            ${model.get().config.host}:${model.get().config.port}.
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

  // 👁️ https://dev.to/blikblum/dry-form-handling-with-lit-19f
  async #connect(e: Event): Promise<void> {
    e.preventDefault();
    const form = e.target as HTMLFormElement;
    if (form) {
      const formData = new FormData(form);
      const config = Object.fromEntries(formData.entries()) as Config;
      this.theState.updateConfig(config);
      try {
        // 👇 try to connect to 3270
        this.connecting = true;
        State.theTn3270?.close();
        State.theTn3270 = await Tn3270.tn3270(
          config.host,
          config.port,
          Emulators[config.emulator] as string
        );
        State.theTn3270.stream$.subscribe({
          // 🔥 need to do a lot more than this
          next: () => {
            if (this.connecting) this.theState.gotoPage(Pages.emulator);
            this.connecting = false;
          },

          // 🔥 WebSocket connection established, but that failed
          error: async (e: any) => {
            console.error(e);
            this.connecting = false;
            this.message = e.reason;
            await this.dialog.show();
            this.theState.gotoPage(Pages.connector);
            State.theTn3270 = null;
          },

          // 👇 normal completion eg: Tn3270.close()
          complete: () => {
            console.log(
              `%c3270 -> Server -> Client %cClosed`,
              'color: palegreen',
              'color: cyan'
            );
            this.theState.gotoPage(Pages.connector);
            State.theTn3270 = null;
          }
        });
      } catch (e: any) {
        // 🔥 tried to upgrade to WebSocket, but that failed
        console.error(e);
        this.connecting = false;
        this.message = `Unable to reach proxy server ${location.hostname}:${location.port}`;
        await this.dialog.show();
        this.theState.gotoPage(Pages.connector);
        State.theTn3270 = null;
      }
    }
  }
}
