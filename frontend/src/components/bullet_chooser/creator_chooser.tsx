import * as React from 'react'

import { User } from 'src/global_types'
import { listEvidenceCreators } from 'src/services'
import BulletChooser, { BulletProps } from 'src/components/bullet_chooser'
import { FilterModified } from 'src/helpers'
import { isNotUndefined } from 'src/helpers/is_not_undefined'

export const CreatorChooser = (props: {
  label: string
  value: Array<BulletProps>
  operationSlug: string
  options: Array<BulletProps>
  onChange: (creators: Array<BulletProps>) => void
  className?: string
  disabled?: boolean
  enableNot?: boolean
}) => {
  return (
    <BulletChooser
      label={props.label}
      options={props.options}
      value={props.value}
      onChange={props.onChange}
      enableNot={props.enableNot}
    />
  )
}

export const ManagedCreatorChooser = (props: {
  operationSlug: string,
  className?: string,
  disabled?: boolean,
  label: string,
  onChange: (creators: Array<BulletProps>) => void,
  value: Array<BulletProps>,
  enableNot?: boolean
}) => {
  const [allCreators, setAllCreators] = React.useState<Array<User>>([])

  const reloadCreators = () => { listEvidenceCreators({ operationSlug: props.operationSlug }).then(setAllCreators) }
  React.useEffect(reloadCreators, [props.operationSlug])

  return (
    <CreatorChooser
      {...props}
      options={allCreators.map(creatorToBulletProps).filter(isNotUndefined)}
    />
  )
}

export const creatorToBulletProps = (creator: FilterModified<User> | undefined): BulletProps | undefined => {
  if (!creator) {
    return undefined
  }
  return {
    id: creator.slug,
    name: [creator.firstName, creator.lastName].join(' '),
    modifier: creator.modifier == 'not' ? "not" : undefined
  }
}

export default CreatorChooser
