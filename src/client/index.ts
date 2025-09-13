export * from '$client/pages/home';
export * from '$client/pages/root';
export * from '$client/pages/screen';

export * from '$client/components/icon';

import '@material/web/button/filled-button.js';
import '@material/web/radio/radio.js';
import '@material/web/textfield/filled-text-field.js';

import { Tn3270 } from '$client/services/tn3270';

const tn3270 = await Tn3270.tn3270('localhost', '3270', 'IBM-3278-4-E');
tn3270?.stream$.subscribe({
  next: () => {},
  error: (error: Error) => console.error(error),
  complete: () => console.log('All done!')
});
