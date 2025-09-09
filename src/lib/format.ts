// 📘 see https://gist.github.com/zentala/1e6f72438796d74531803cc3833c039c

export function formatBytes(
  bytes: number,
  numDP = 2,
  numPad = 10
): string {
  // 👇 default for zero size
  let result = '0b';
  if (bytes > 0) {
    // 👇 reduce raw size to best unit
    const k = 1024,
      sizes = ['b', 'k', 'M', 'G', 'T'],
      i = Math.floor(Math.log(bytes) / Math.log(k));
    let size = String(
      parseFloat((bytes / Math.pow(k, i)).toFixed(numDP))
    );
    // 👇 supply trailing zeroes as necesssary
    const ix = size.indexOf('.');
    if (ix > 0) {
      const numAfterDP = size.length - (ix + 1);
      if (numDP > numAfterDP)
        size = size.padEnd(size.length + (numDP - numAfterDP), '0');
    }
    result = `${size}${sizes[i]}`;
  }
  // 👇 all done!
  return result.padStart(numPad);
}
