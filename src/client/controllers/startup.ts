import { LitElement } from 'lit';
import { ReactiveController } from 'lit';
import { ReactiveControllerHost } from 'lit';

import { config } from '$client/config';
import { enablePatches } from 'immer';
import { nextTick } from '$lib/delay';

// 👇 finding the patches is "expensive" so we feature flag logging
if (config.logStateChanges) enablePatches();

// 📘 manage startup tasks

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
