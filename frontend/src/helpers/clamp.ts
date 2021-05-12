
// clamp constrains n to a number between min and max, inclusive
// clamp(-1, 0, 4) == 0
// clamp(5, 0, 4) == 4
// clamp(1, 0, 4) == 1
// clamp(2.5, 0, 4) == 2.5
export const clamp = (n: number, min: number, max: number) => Math.max(Math.min(max, n), min)
