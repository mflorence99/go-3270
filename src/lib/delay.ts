// 📘 Promise-based delay

export function delay(ms: number): Promise<number> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function nextTick(): Promise<number> {
  return delay(0);
}
