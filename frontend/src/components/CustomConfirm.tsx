import { Button } from 'react-bootstrap'
import { confirmAlert } from 'react-confirm-alert'
import 'assets/scss/confirm-alert.css'

type customConfirmOptionProps = {
  onConfirm: () => void,
  onCancel: () => void,
  message: string
}

const CustomConfirm = ({onConfirm, onCancel, message}: customConfirmOptionProps) => {

  return confirmAlert({
    customUI: ({ onClose }) => {
      const handleConfirm = () => {
        onConfirm()
        onClose()
      }
      const handleCancel = () => {
        onCancel()
        onClose()
      }
      return (
        <div className='react-confirm-alert-body'>
          <p>
            {message}
          </p>
          <div className='react-confirm-alert-button-group'>
            <Button variant='success' className='btn-rounded btn-sm' onClick={handleConfirm}>
              OK
            </Button>
            <Button variant='light' className='btn-rounded btn-sm' onClick={handleCancel}>
              cancel
            </Button>
          </div>
        </div>
      )
    },
  })
}

export default CustomConfirm