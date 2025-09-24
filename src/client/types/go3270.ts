export type Go3270 = {
  close: () => Uint8Array;
  datastream: (bytes: Uint8Array) => Uint8Array;
  restore: (bytes: Uint8Array) => void;
  testPattern: () => void;
};
