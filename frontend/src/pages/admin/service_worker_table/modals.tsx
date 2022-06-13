// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'

import { ServiceWorker } from 'src/global_types'
import { useForm, useFormField } from 'src/helpers'

import ChallengeModalForm from 'src/components/challenge_modal_form'
import { default as Input } from 'src/components/input'
import ModalForm from 'src/components/modal_form'
import { createServiceWorker, deleteServiceWorkers, restoreServiceWorker, updateServiceWorker } from 'src/services'
import { SourcelessCodeblock } from 'src/components/code_block'
import Label from 'src/components/with_label'

export const DeleteServiceModal = (props: {
  worker: ServiceWorker,
  onRequestClose: () => void,
}) => (
  <ChallengeModalForm
    modalTitle="Delete Service"
    warningText="This will remove the service from the system. Existing metadata will be kept, but no new metadata will be added by this worker. Note that queued work will still be completed and recorded."
    submitText="Delete"
    challengeText={props.worker.name}
    handleSubmit={() => deleteServiceWorkers({ id: props.worker.id })}
    onRequestClose={props.onRequestClose}
  />
)

export const RestoreServiceModal = (props: {
  worker: ServiceWorker,
  onRequestClose: () => void,
}) => (
  <ChallengeModalForm
    modalTitle="Restore Service"
    warningText="This will restore the service to active use. All evidence created after this point can now use this worker."
    submitText="Restore"
    handleSubmit={() => restoreServiceWorker({ id: props.worker.id })}
    onRequestClose={props.onRequestClose}
  />
)


export const AddEditServiceWorkerModal = (props: {
  worker?: ServiceWorker,
  onRequestClose: () => void,
}) => {
  const serviceName = useFormField<string>(props.worker?.name ?? "")
  const serviceConfig = useFormField<string>(props.worker?.config ?? "")

  const handleSubmit = () => {
    if (serviceName.value.trim() === '') {
      return Promise.reject(new Error("Service name should contain some value"))
    }
    if (serviceConfig.value.trim() === '') {
      return Promise.reject(new Error("Please provide a configuration"))
    }
    try {
      JSON.parse(serviceConfig.value)
    }
    catch (err) {
      const rejection = "JSON config was not parsable." + (err instanceof Error
        ? ` Error: ${err.message}`
        : ""
      )
      return Promise.reject(new Error(rejection))
    }

    const commonProps = {
      name: serviceName.value,
      config: serviceConfig.value,
    }
    return props.worker
      ? updateServiceWorker({ ...commonProps, id: props.worker.id })
      : createServiceWorker({ ...commonProps })
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
      <Input label="Service Name" {...serviceName} />
      {/* TODO: this might be a good place to stick a flag for rendering purposes */}
      {/* <TextArea label="Service Config" {...serviceConfig} /> */}
      <Label label='Service Config'>
        <SourcelessCodeblock
          code={serviceConfig.value}
          onChange={serviceConfig.onChange}
          editable
          language='json'
        />
      </Label>
    </ModalForm>
  )
}
