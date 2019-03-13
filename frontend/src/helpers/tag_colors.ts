// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

// We store color names instead of hex codes in the
// database to allow color scheme changes in the future
const TAG_COLORS: {[key: string]: number} = {
  blue           : 0x0E5A8A,
  yellow         : 0xA67908,
  green          : 0x0A6640,
  indigo         : 0x5642A6,
  orange         : 0xA66321,
  pink           : 0xA82255,
  red            : 0xA82A2A,
  teal           : 0x008075,
  vermilion      : 0x9E2B0E,
  violet         : 0x5C255C,

  lightBlue      : 0x48AFF0,
  lightYellow    : 0xFFC940,
  lightGreen     : 0x3DCC91,
  lightIndigo    : 0xAD99FF,
  lightOrange    : 0xFFB366,
  lightPink      : 0xFF66A1,
  lightRed       : 0xFF7373,
  lightTeal      : 0x2EE6D6,
  lightVermilion : 0xFF6E4A,
  lightViolet    : 0xC274C2,
}

export const tagColorNames = Object.keys(TAG_COLORS)

export const tagColorNameToColor = (key: string): number =>
  TAG_COLORS[key] || 0x000000

export const randomTagColorName = (): string =>
  tagColorNames[Math.round(Math.random() * (tagColorNames.length - 1))]
