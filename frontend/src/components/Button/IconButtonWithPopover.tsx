import classNames from 'classnames'
import { Button, OverlayTrigger, Tooltip } from 'react-bootstrap'

function IconButtonWithPopover({ popTitle, icon, name, ...props }: any) {
  const renderTooltip = (props: any) => {
    return (
      <Tooltip id={`PopoverFocus-${name}`} {...props}>
        {popTitle}
      </Tooltip>
    )
  }
  return (
    <>
      <OverlayTrigger placement='bottom' overlay={renderTooltip}>
        <Button {...props} style={{padding: '1px'}} size="sm">
          <i className={classNames('mdi', 'ms-1', 'me-1', 'font-17', icon)}></i>
        </Button>
      </OverlayTrigger>
    </>
  )
}

export default IconButtonWithPopover
