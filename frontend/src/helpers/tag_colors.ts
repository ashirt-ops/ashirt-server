export type TagColor =
  | 'blue'
  | 'yellow'
  | 'green'
  | 'indigo'
  | 'orange'
  | 'pink'
  | 'red'
  | 'teal'
  | 'vermilion'
  | 'violet'
  | 'lightBlue'
  | 'lightYellow'
  | 'lightGreen'
  | 'lightIndigo'
  | 'lightOrange'
  | 'lightPink'
  | 'lightRed'
  | 'lightTeal'
  | 'lightVermilion'
  | 'lightViolet'
  | 'disabledGray'

// We store color names instead of hex codes in the
// database to allow color scheme changes in the future
const TAG_COLORS: Record<TagColor, number> = {
  blue: 0x0e5a8a,
  yellow: 0xa67908,
  green: 0x0a6640,
  indigo: 0x5642a6,
  orange: 0xa66321,
  pink: 0xa82255,
  red: 0xa82a2a,
  teal: 0x008075,
  vermilion: 0x9e2b0e,
  violet: 0x5c255c,

  lightBlue: 0x48aff0,
  lightYellow: 0xffc940,
  lightGreen: 0x3dcc91,
  lightIndigo: 0xad99ff,
  lightOrange: 0xffb366,
  lightPink: 0xff66a1,
  lightRed: 0xff7373,
  lightTeal: 0x2ee6d6,
  lightVermilion: 0xff6e4a,
  lightViolet: 0xc274c2,

  disabledGray: 0x606060,
}

export const disabledGray = 'disabledGray'

export const tagColorNames = Object.keys(TAG_COLORS)

export const tagColorNameToColor = (key: string): number => {
  // typescript complains that a string won't match, but we're accounting for that here.
  return TAG_COLORS[key as TagColor] ?? 0x000000
}

export const shiftColor = (key: TagColor): TagColor => {
  const shifted: Record<TagColor, TagColor> = {
    blue: 'lightBlue',
    yellow: 'lightYellow',
    green: 'lightGreen',
    indigo: 'lightIndigo',
    orange: 'lightOrange',
    pink: 'lightPink',
    red: 'lightRed',
    teal: 'lightTeal',
    vermilion: 'lightVermilion',
    violet: 'lightViolet',

    lightBlue: 'blue',
    lightYellow: 'yellow',
    lightGreen: 'green',
    lightIndigo: 'indigo',
    lightOrange: 'orange',
    lightPink: 'pink',
    lightRed: 'red',
    lightTeal: 'teal',
    lightVermilion: 'vermilion',
    lightViolet: 'violet',

    disabledGray: 'disabledGray',
  }
  return shifted[key]
}

export const randomTagColorName = (): string =>
  tagColorNames[Math.round(Math.random() * (tagColorNames.length - 1))]
