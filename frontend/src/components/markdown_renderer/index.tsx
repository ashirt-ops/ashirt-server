import SyncMarkdownRenderer from './sync_renderer'

const MarkdownRenderer = (props: { className?: string; children: string }) => (
  <div className={props.className}>
    <SyncMarkdownRenderer children={props.children} />
  </div>
)
export default MarkdownRenderer
