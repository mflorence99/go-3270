import { LitElement } from 'lit';
import { ReactiveController } from 'lit';
import { ReactiveControllerHost } from 'lit';

import { nextTick } from '$lib/delay';

// 📘 manage startup tasks

export class StartupController implements ReactiveController {
  host: ReactiveControllerHost;

  constructor(host: ReactiveControllerHost) {
    (this.host = host).addController(this);
  }

  hostConnected(): void {
    // 👇 tag the host with the "ready" class
    if (this.host instanceof LitElement) {
      nextTick().then(() =>
        (this.host as LitElement).classList.add('ready')
      );
    }
  }

  hostDisconnected(): void {}
}
