import { OverlayTrigger, Tooltip } from 'react-bootstrap'

function PopoverLabel({ children, name, marginLeft, marginRight, size, ...props }: any) {
  const renderTooltip = (props: any) => {
    return (
      <Tooltip id={`PopoverFocus-${name}`} {...props}>
        {children}
      </Tooltip>
    )
  }
  return (
    <>
      <OverlayTrigger placement='bottom' overlay={renderTooltip}>
        <i
          className='mdi mdi-chat-question-outline'
          style={{ marginLeft: `${marginLeft || 0}px`, marginRight: `${marginRight || 0}px`, fontSize: `${size}px` }}
        />
      </OverlayTrigger>
    </>
  )
}

export default PopoverLabel
