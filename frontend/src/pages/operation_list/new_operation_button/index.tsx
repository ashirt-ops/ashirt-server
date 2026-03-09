import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

const NewOperationButton = (props: { onClick: () => void }) => (
  <button className={cx('root')} onClick={props.onClick}>
    <div className={cx('circle', 'plus')} />
    New Operation
  </button>
)
export default NewOperationButton
