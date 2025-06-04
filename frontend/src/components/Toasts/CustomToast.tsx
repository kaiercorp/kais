import { useToggle } from 'hooks'
import { Button, Toast } from 'react-bootstrap'

type Props = {
  index: number
  message: string
  onAutoClose: (index: number) => void
}

const CustomToast = ({index, message, onAutoClose}: Props) => {
  const [isOpenToast, , , hideToast] = useToggle(true)
  const handleHideToast = () => {
    onAutoClose(index)
    hideToast()
  }

  return (
    <Toast
      className='d-flex align-items-center'
      show={isOpenToast}
      onClose={handleHideToast}
      delay={1500}
      autohide
    >
      <Toast.Body>{message}</Toast.Body>
      <Button variant="" onClick={handleHideToast} className='btn-close ms-auto me-2'></Button>
    </Toast>
  )
} 

export default CustomToast