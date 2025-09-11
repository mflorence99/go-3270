// 🟦 Conversion between EBCDIC and ASCII

//    @see https://stackoverflow.com/questions/25367120/example-ebcdic-file-for-java-program-to-convert-ebcdic-to-ascii

// 🔥 starts at position 64

const ebcdic: any[] = [
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '¢',
  '.',
  '<',
  '(',
  '+',
  '|',
  '&',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '!',
  '$',
  '*',
  ')',
  ';',
  '¬',
  '-',
  '/',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '|',
  ',',
  '%',
  '_',
  '>',
  '?',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '`',
  ':',
  '#',
  '@',
  "'",
  '=',
  '\"',
  '\u00a0',
  'a',
  'b',
  'c',
  'd',
  'e',
  'f',
  'g',
  'h',
  'i',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  'j',
  'k',
  'l',
  'm',
  'n',
  'o',
  'p',
  'q',
  'r',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  's',
  't',
  'u',
  'v',
  'w',
  'x',
  'y',
  'z',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '`',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '{',
  'A',
  'B',
  'C',
  'D',
  'E',
  'F',
  'G',
  'H',
  'I',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '}',
  'J',
  'K',
  'L',
  'M',
  'N',
  'O',
  'P',
  'Q',
  'R',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\\',
  '\u00a0',
  'S',
  'T',
  'U',
  'V',
  'W',
  'X',
  'Y',
  'Z',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '0',
  '1',
  '2',
  '3',
  '4',
  '5',
  '6',
  '7',
  '8',
  '9',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0'
];

// reverse engineer ascii
const ascii = new Array(256);

for (let i = 0; i < ebcdic.length; i++) {
  if (ebcdic[i] !== null) ascii[ebcdic[i].charCodeAt(0)] = i + 64;
}

/**
 * Convert from ASCII to EBCDIC
 */

export function a2e(a: string): number[] {
  const e: number[] = [];
  for (let i = 0; i < a.length; i++) {
    if (a[i] === ' ' || a[i] === '\u00a0') e[i] = 0x40;
    else e[i] = ascii[a.charCodeAt(i)];
  }
  return e;
}

/**
 * Convert from EBCDIC to ASCII
 */

export function e2a(e: number[]): string {
  let a = '';
  for (let i = 0; i < e.length; i++) {
    // @ts-ignore 🔥 we know e[i] is always valid
    if (e[i] >= 64) a += ebcdic[e[i] - 64];
    else a += '\u2022';
  }
  return a;
}
