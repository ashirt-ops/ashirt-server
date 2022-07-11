// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { ServiceWorker } from 'src/global_types'
import { listServiceWorkers } from 'src/services'
import { isNotUndefined } from 'src/helpers/is_not_undefined'
import { FilterModified } from 'src/helpers'

import BulletChooser, { BulletProps } from 'src/components/bullet_chooser'

export const ServiceWorkerChooser = (props: {
  label: string
  value: Array<BulletProps>
  operationSlug: string
  options: Array<BulletProps>
  onChange: (workers: Array<BulletProps>) => void
  className?: string
  disabled?: boolean
  enableNot?: boolean
}) => {
  return (
    <BulletChooser
      className={props.className}
      label={props.label}
      options={props.options}
      value={props.value}
      onChange={props.onChange}
      enableNot={props.enableNot}
    />
  )
}

export const ManagedServiceWorkerChooser = (props: {
  operationSlug: string,
  className?: string,
  disabled?: boolean,
  label: string,
  onChange: (workers: Array<BulletProps>) => void,
  value: Array<BulletProps>,
  enableNot?: boolean
}) => {
  const [allWorkers, setAllWorkers] = React.useState<Array<ServiceWorker>>([])

  const reloadWorkers = () => {
    listServiceWorkers()
      .then(list => setAllWorkers(list.filter(item => !item.deleted)))
  }
  React.useEffect(reloadWorkers, [props.operationSlug])

  return (
    <ServiceWorkerChooser
      {...props}
      options={allWorkers.map(workerToBulletProps).filter(isNotUndefined)}
    />
  )
}

export const workerToBulletProps = (worker: FilterModified<ServiceWorker> | undefined): BulletProps | undefined => {
  if (!worker) {
    return undefined
  }
  return {
    id: worker.id,
    name: worker.name,
    modifier: worker.modifier == 'not' ? "not" : undefined
  }
}

export default ServiceWorkerChooser
