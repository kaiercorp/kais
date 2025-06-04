import classNames from 'classnames'
import { Button } from 'react-bootstrap'

const IconLabelButton = ({ icon, label, color, onClick }: any) => {
    return (
        <span>
            <Button variant={`outline-${color}`} onClick={onClick} style={{ width: '22px', height: '22px', padding: '0px' }}>
                <i className={classNames('mdi', icon)}></i>
            </Button>
            <span style={{ marginLeft: '6px'}}>{label}</span>
        </span>
    )
}

export default IconLabelButton
