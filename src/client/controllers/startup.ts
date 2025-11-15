import { LitElement } from 'lit';
import { ReactiveController } from 'lit';
import { ReactiveControllerHost } from 'lit';

import { nextTick } from '$client/utils/delay';

// 🟧 Manage startup tasks

export class Startup implements ReactiveController {
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
