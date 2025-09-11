export * from '$client/components/icon';
export * from '$client/components/root';
export * from '$client/components/test';

import '@material/web/button/filled-button.js';
import '@material/web/checkbox/checkbox.js';

import { Tn3270 } from '$client/services/tn3270';

const tn3270 = await Tn3270.tn3270('localhost', '3270', 'IBM-3278-4-E');
tn3270?.stream$.subscribe({
  next: (data: Uint8Array) => console.log(data),
  error: (error: Error) => console.log(error),
  complete: () => console.log('All done!')
});
