import { Button } from 'react-bootstrap'
import { confirmAlert } from 'react-confirm-alert'
import 'assets/scss/confirm-alert.css'

type customConfirmOptionProps = {
  message: string
}

const CustomAlert = ({message}: customConfirmOptionProps) => {

  return confirmAlert({
    customUI: ({ onClose }) => {
      return (
        <div className='react-confirm-alert-body'>
          <p>
            {message}
          </p>
          <div className='react-confirm-alert-button-group'>
            <Button variant='success' className='btn-rounded btn-sm' onClick={onClose}>
              OK
            </Button>
          </div>
        </div>
      )
    },
  })
}

export default CustomAlert