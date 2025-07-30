export function $textSplit(text: string): string[] {
  return (text || '').split('&n&').filter((n) => n)
}
