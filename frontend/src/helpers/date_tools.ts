export function nullableDateStrToNullableDate(d: string | null): Date | null {
  if (d) {
    return new Date(d)
  }
  return null
}
