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

// 📘 base state class

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

// 📘 the entire state of the app

export type CLUT = Record<string, [string, string]>;

export const defaultCLUT: CLUT = {
  // 👇 [normal, highlight]
  black: ['#505050', '#111138'],
  blue: ['#3366CC', '#0078FF'],
  red: ['#E06666', '#D40000'],
  pink: ['#FFB3DA', '#FF69B4'],
  green: ['#88DD88', '#00AA00'],
  turquoise: ['#99E8DD', '#00C8AA'],
  yellow: ['#F0F080', '#FFB266'],
  white: ['#888888', '#FFFFFF']
};

export type Config = {
  color: string;
  device: string;
  dims: [number, number];
  fontSize: string;
  host: string;
  model: string;
  port: string;
};

export const defaultConfig: Config = {
  color: 'green',
  device: '3279',
  // 👇 [rows, cols]
  dims: [24, 80],
  fontSize: '14',
  host: 'localhost',
  model: '2',
  port: '3270'
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
  color = computed((): string[] => {
    const clut = this.model.get().clut;
    const model = this.model.get().config.model;
    // @ts-ignore 🔥 we know this is always valid
    return model === '3278' ? clut['green'] : clut['white'];
  });

  cursorAt = computed(() => {
    const cursorAt = this.model.get().status.cursorAt;
    if (cursorAt >= 0) {
      const dims = this.model.get().config.dims;
      // @ts-ignore 🔥 we know this is always valid
      return `${String(Math.trunc(cursorAt / dims[0]) + 1).padStart(3, '0')}/${String((cursorAt % dims[0]) + 1).padStart(3, '0')}`;
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
