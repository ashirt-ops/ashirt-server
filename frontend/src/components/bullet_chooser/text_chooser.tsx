// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'

import BulletChooser, { BulletProps } from 'src/components/bullet_chooser'
import Tag from 'src/components/tag'

export const TextChooser = (props: {
  label: string
  value: Array<BulletProps>
  onChange: (value: Array<BulletProps>) => void
  className?: string
  tagColorName?: string
  disabled?: boolean
}) => {
  return (
    <BulletChooser
      label={props.label}
      options={[]}
      value={props.value}
      onChange={props.onChange}
      onNoValueSelected={async (v) => {
        const bullet = textToBulletProps(v.trim())
        return bullet ?? null
      }}
      noValueRenderer={(v) => {
        if (v.trim() == '') {
          return <>Type to add a term</>
        }
        return <>Add Term: <Tag name={v} color='' /></>
      }}
    />
  )
}

export const textToBulletProps = (text: string | undefined): BulletProps | undefined => {
  if (!text) {
    return undefined
  }
  return {
    id: text,
    name: text,
    modifier: undefined
  }
}

export default TextChooser
