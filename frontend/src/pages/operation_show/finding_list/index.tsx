import * as React from 'react'
import FindingsTable from './findings_table'
import Layout from '../layout'
import { Finding } from 'src/global_types'
import { DeleteFindingModal, EditFindingModal } from '../finding_modals'
import { useNavigate, useLocation, useParams } from 'react-router-dom'
import { getFindings } from 'src/services'
import { useWiredData, useModal, renderModals } from 'src/helpers'
import { mkNavTo } from 'src/helpers/navigate-to-query'

export default () => {
  const { slug } = useParams<{ slug: string }>()
  const operationSlug = slug! // useParams puts everything in a partial, so our type above doesn't matter.
  const location = useLocation()
  const navigate = useNavigate()

  const query: string = new URLSearchParams(location.search).get('q') || ''

  const wiredFindings = useWiredData(React.useCallback(() => getFindings({
    operationSlug,
    query,
  }), [operationSlug, query]))

  const navTo = mkNavTo({
    navTo: navigate,
    slug: operationSlug
  })

  const editFindingModal = useModal<{ finding: Finding }>(modalProps => (
    <EditFindingModal {...modalProps} onEdited={wiredFindings.reload} operationSlug={operationSlug} />
  ))
  const deleteFindingModal = useModal<{ finding: Finding }>(modalProps => (
    <DeleteFindingModal {...modalProps} onDeleted={wiredFindings.reload} operationSlug={operationSlug} />
  ))

  return (
    <Layout
      onFindingCreated={wiredFindings.reload}
      onNavigate={navTo}
      operationSlug={operationSlug}
      query={query}
      view="findings"
    >
      {wiredFindings.render(findings => (
        <div style={{ padding: 20 }}>
          <FindingsTable
            findings={findings}
            onDelete={finding => deleteFindingModal.show({ finding })}
            onEdit={finding => editFindingModal.show({ finding })}
            operationSlug={operationSlug}
          />
        </div>
      ))}

      {renderModals(editFindingModal, deleteFindingModal)}
    </Layout>
  )
}
