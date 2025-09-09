import { State } from '$client/state/state';

import { computed } from '$client/types/signals';
import { createContext } from '@lit/context';

// 📘 a conceptual model for real states
//    may morph into a real app-state

export type AppStateModel = {
  x: number;
  y: string;
  z: boolean;
};

const defaultState: AppStateModel = {
  x: 1000,
  y: '2000',
  z: true
};

export class AppState extends State<AppStateModel> {
  // 👇 just an example of a computed property
  asJSON = computed(() => JSON.stringify(this.model.get()));

  constructor(key: string) {
    super(defaultState, key, true);
  }

  // 👇 just an example of a mutator
  incrementX(x: number): void {
    this.mutate((state) => void (state.x += x));
  }
}

export const appStateContext = createContext<AppState>(
  Symbol('app-state')
);
