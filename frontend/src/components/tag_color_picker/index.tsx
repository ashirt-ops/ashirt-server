// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import WithLabel from 'src/components/with_label'
import classnames from 'classnames/bind'
import {tagColorNames, tagColorNameToColor} from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

const colorNameToStyle = (colorName: string) => ({
  background: '#' + ('000000' + tagColorNameToColor(colorName).toString(16)).substr(-6)
})

export default (props: {
  disabled: boolean,
  label: string,
  onChange: (v: string) => void,
  value: string,
}) => (
  <WithLabel label={props.label}>
    <div className={cx('colors')}>
      <div className={cx('current')} style={colorNameToStyle(props.value)} />
      {tagColorNames.map(colorName => (
        <div
          key={colorName}
          className={cx('color', {selected: props.value === colorName})}
          style={colorNameToStyle(colorName)}
          onClick={() => props.onChange(colorName)}
        />
      ))}
    </div>
  </WithLabel>
)
