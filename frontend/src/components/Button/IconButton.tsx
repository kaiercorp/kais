import classNames from 'classnames'
import { Button } from 'react-bootstrap'

const IconButton = ({ icon, color, onClick }: any) => {
  return (
    <Button variant={`outline-${color}`} onClick={onClick} style={{width: '26px', height: '26px', padding: '0px'}}>
      <i className={classNames('mdi', 'ms-1', 'me-1', icon)}></i>
    </Button>
  )
}

export default IconButton
