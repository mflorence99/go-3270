// 📘 3270 device types

export const Colors: Record<string, string> = {
  green: '#61b064',
  blue: '#42a5f5',
  orange: '#eb9a25',
  white: '#f9f9f9'
};

export const Dimensions: Record<string, [number, number]> = {
  // 👇 [width, var(--toolbar-height)]
  '1': [80, 12],
  '2': [80, 24],
  '3': [80, 32],
  '4': [80, 43],
  '5': [132, 27]
};

export const Emulators: Record<string, string> = {
  '1': 'IBM-3278-1',
  '2': 'IBM-3278-2',
  '3': 'IBM-3278-3',
  '4': 'IBM-3278-4',
  '5': 'IBM-3278-5'
};
