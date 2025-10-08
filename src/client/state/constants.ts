// 📘 state constants

// 🔥 it would be nice to get these colors from the CLUT in the Go code, but I think not worth losing the ability to treat them as static resources, and instead dependent on the load of the WASM code

export const defaultColor = ['#00AA00', '#88DD88'];

export const Colors: Record<string, string[]> = {
  green: defaultColor,
  blue: ['#0078FF', '#3366CC'],
  orange: ['#FF8000', '#FFB266'],
  white: ['#888888', '#FFFFFF']
};

export const defaultDimensions: [number, number] = [80, 24];

export const Dimensions: Record<string, [number, number]> = {
  // 👇 [width, height]
  '1': [40, 12],
  '2': defaultDimensions,
  '3': [80, 32],
  '4': [80, 43],
  '5': [132, 27]
};

export const Emulators: Record<string, string> = {
  '1': 'IBM-3277-1',
  '2': 'IBM-3277-2',
  '3': 'IBM-3278-3',
  '4': 'IBM-3278-4',
  '5': 'IBM-3278-5'
};
