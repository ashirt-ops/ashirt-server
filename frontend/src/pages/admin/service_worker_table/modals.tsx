// Copyright 2022, Yahoo Inc
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'

import { ServiceWorker } from 'src/global_types'
import { useForm, useFormField } from 'src/helpers'

import ChallengeModalForm from 'src/components/challenge_modal_form'
import { default as Input, TextArea } from 'src/components/input'
import ModalForm from 'src/components/modal_form'
import { createServiceWorker, deleteServiceWorkers, updateServiceWorker } from 'src/services/service_workers'

export const DeleteServiceModal = (props: {
  worker: ServiceWorker,
  onRequestClose: () => void,
}) => <ChallengeModalForm
    modalTitle="Delete Service"
    warningText="This will remove the service from the system. Existing metadata will be kept, but no new metadata will be added by this worker."
    submitText="Delete"
    challengeText={props.worker.name}
    handleSubmit={() => deleteServiceWorkers({ id: props.worker.id })}
    onRequestClose={props.onRequestClose}
  />

export const AddEditServiceWorkerModal = (props: {
  worker?: ServiceWorker,
  onRequestClose: () => void,
}) => {
  const serviceName = useFormField<string>(props.worker?.name ?? "")
  const serviceConfig = useFormField<string>(props.worker?.config ?? "")

  const handleSubmit = () => {
    const commonProps = {
      name: serviceName.value,
      serviceType: 'aws',
      config: serviceConfig.value,
    }
    return props.worker
      ? updateServiceWorker({...commonProps, id: props.worker.id})
      : createServiceWorker({...commonProps})
  }

  const formComponentProps = useForm({
    fields: [serviceName, serviceConfig],
    handleSubmit,
    onSuccess: props.onRequestClose
  })
  return (
    <ModalForm
      title={props.worker ? "Edit Service Worker" : "Create Service Worker"}
      submitText={props.worker ? "Save" : "Create"}
      cancelText="Cancel"
      onRequestClose={props.onRequestClose}
      {...formComponentProps}
    >
      <Input label="Service name" {...serviceName} />
      <TextArea label="Service Config" {...serviceConfig} />
    </ModalForm>
  )
}
