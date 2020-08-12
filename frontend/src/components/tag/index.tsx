// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import {tagColorNameToColor, disabledGray} from 'src/helpers/tag_colors'
const cx = classnames.bind(require('./stylesheet'))

// Returns an appropriate foreground color (white or black) for the given background color
function getFgColorFromColor(color: number) {
  const r = Math.floor(color / 0x10000);
  const g = Math.floor(color / 0x100 % 0x100);
  const b = color % 0x100;
  const yiq = (r * 299 + g * 587 + b * 114) / 1000;
  return yiq < 128 ? "#fff" : "#000";
}

export function tagColorStyle(colorName: string): React.CSSProperties {
  const color = tagColorNameToColor(colorName)
  return {
    backgroundColor: '#' + ('000000' + color.toString(16)).substr(-6),
    color: getFgColorFromColor(color),
  }
}

export default (props: {
  name: string,
  color: string,
  selected?: boolean,
  onClick?: () => void,
  className?: string,
  disabled?: boolean,
}) => (
  <span
    children={props.name}
    className={cx('root', { clickable: props.onClick != null, selected: !!props.selected }, props.className)}
    style={tagColorStyle(props.disabled ? disabledGray : props.color)}
    onClick={props.onClick}
  />
)
