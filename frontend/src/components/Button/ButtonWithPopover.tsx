import { Button, OverlayTrigger, Tooltip } from 'react-bootstrap'

function ButtonWithPopover({ popTitle, children, name, ...props }: any) {
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
        <Button {...props}>{children}</Button>
      </OverlayTrigger>
    </>
  )
}

export default ButtonWithPopover
