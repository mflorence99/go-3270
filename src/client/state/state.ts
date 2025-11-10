import { Signal } from '@lit-labs/signals';

import { computed } from '$client/types/signals';
import { config } from '$client/config';
import { createContext } from '@lit/context';
import { effect } from '$client/types/signals';
import { enablePatches } from 'immer';
import { produce } from 'immer';
import { signal } from '@lit-labs/signals';

import StackTrace from 'stacktrace-js';

enablePatches();

// 🟧 Base state class

abstract class Base<T> {
  delta: Partial<T> = {};

  // 👇 the signal that is the state itself
  model: Signal.State<T>;

  constructor(defaultState: T, key: string, persist: boolean) {
    if (persist) {
      const raw = localStorage.getItem(key);
      const persistedState = raw ? JSON.parse(raw) : defaultState;
      // 👇 watch out! the new state maybe completely different from the old
      //    better to also remove any deprecated properties
      //    but this is simpler and always works, even for optionals
      this.model = signal<T>({ ...defaultState, ...persistedState });
      effect(() =>
        localStorage.setItem(key, JSON.stringify(this.model.get()))
      );
    } else this.model = signal<T>(defaultState);
  }

  mutate(mutator: (state: T) => void): void {
    // 👇 finding the caller is "expensive" so we feature flag logging
    let caller: string | undefined;
    if (config.logStateChanges) {
      const frame = StackTrace.getSync()[1];
      caller = frame?.functionName;
    }
    // 👇 the "old" state
    const prevState = this.model.get();
    if (config.logStateChanges)
      console.log(
        '%c👈 prev state',
        'color: palegreen',
        caller,
        prevState
      );
    // 👇 derive the new state and accumulate deltas for inspection
    const newState = produce(prevState, mutator, (patches) => {
      this.delta = {};
      for (const patch of patches) {
        patch.path.reduce((acc: any, fld, ix, arr) => {
          acc[fld] = ix === arr.length - 1 ? patch.value : {};
          return acc[fld];
        }, this.delta);
      }
      if (config.logStateChanges)
        console.log(
          `%c🆕 patches... %c${caller} %c👉${JSON.stringify(this.delta)}`,
          'color: khaki',
          'color: white',
          'color: wheat'
        );
    });
    // 👇 the "new" state and (potentially) the patches that produced it
    if (config.logStateChanges)
      console.log(
        '%c👉 next state',
        'color: skyblue',
        caller,
        newState
      );
    this.model.set(newState);
  }
}

// 🟧 The entire state of the app

export type CLUT = Record<number, [string, string]>;

export const defaultCLUT: CLUT = {
  0xf0: ['#202020', 'Background'],
  0xf1: ['#4169E1', 'Blue'],
  0xf2: ['#FF0000', 'Red'],
  0xf3: ['#EE82EE', 'Pink'],
  0xf4: ['#04c304', 'Green'],
  0xf5: ['#40E0D0', 'Turquiose'],
  0xf6: ['#FFFF00', 'Yellow'],
  0xf7: ['#FFFFFF', 'Foreground'],
  0xf8: ['#202020', 'Black'],
  0xf9: ['#0000CD', 'Deep Blue'],
  0xfa: ['#FFA500', 'Orange'],
  0xfb: ['#800080', 'Purple'],
  0xfc: ['#90EE90', 'Pale Green'],
  0xfd: ['#AFEEEE', 'Pale Turquoise'],
  0xfe: ['#C0C0C0', 'Grey'],
  0xff: ['#E2E2E9', 'White']
};

export type Config = {
  device: string;
  dims: [number, number];
  fontSize: string;
  host: string;
  model: string;
  port: string;
  screenshot: string;
};

export const defaultConfig: Config = {
  device: '3279',
  // 👇 [rows, cols]
  dims: [24, 80],
  fontSize: '14',
  host: 'localhost',
  model: '2',
  port: '3270',
  screenshot: ''
};

export type Status = {
  alarm: boolean;
  cursorAt: number;
  error: boolean;
  locked: boolean;
  message: string;
  numeric: boolean;
  protected: boolean;
  waiting: boolean;
};

export const defaultStatus: Status = {
  alarm: false,
  cursorAt: 0,
  error: false,
  locked: false,
  message: '',
  numeric: false,
  protected: false,
  waiting: false
};

export type StateModel = {
  clut: CLUT;
  config: Config;
  status: Status;
};

const defaultState: StateModel = {
  clut: defaultCLUT,
  config: defaultConfig,
  status: defaultStatus
};

export class State extends Base<StateModel> {
  cursorAt = computed(() => {
    const cursorAt = this.model.get().status.cursorAt;
    if (cursorAt >= 0) {
      const dims = this.model.get().config.dims;
      // @ts-ignore 🔥 we know this is always valid
      return `${String(Math.trunc(cursorAt / dims[1]) + 1).padStart(3, '0')}/${String((cursorAt % dims[1]) + 1).padStart(3, '0')}`;
    } else return '';
  });

  constructor(key: string) {
    super(defaultState, key, true);
  }

  resetStatus(): void {
    this.mutate((state) => void (state.status = defaultState.status));
  }

  // 🔥 note all or nothinmg for CLUT!
  updateCLUT(clut: CLUT): void {
    this.mutate((state) => void (state.clut = clut));
  }

  updateConfig(config: Partial<Config>): void {
    this.mutate(
      (state) =>
        void ((state.config = { ...state.config, ...config }),
        (state.status = defaultState.status))
    );
  }

  updateStatus(status: Partial<Status>): void {
    this.mutate(
      (state) => void (state.status = { ...state.status, ...status })
    );
  }
}

export const stateContext = createContext<State>(Symbol('state'));
