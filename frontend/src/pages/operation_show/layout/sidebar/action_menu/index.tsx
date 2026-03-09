import classnames from 'classnames/bind'
import { default as Menu, MenuItem, MenuSeparator } from 'src/components/menu'
const cx = classnames.bind(require('./stylesheet'))

const ActionMenu = (props: { name: string; query: string; onDelete: () => void }) => (
  <Menu>
    <div className={cx('top')}>
      <div className={cx('name')}>{props.name}</div>
      <div className={cx('query')}>{props.query}</div>
    </div>
    <MenuSeparator />
    <MenuItem icon={require('./delete.svg')} onClick={props.onDelete}>
      Delete Query
    </MenuItem>
  </Menu>
)
export default ActionMenu
