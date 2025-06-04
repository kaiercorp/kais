import { Button, Modal } from "react-bootstrap"
import { useTranslation } from 'react-i18next'

import { PopoverLabel } from 'components'
import { BaseModalTitleType } from 'common'

type ModalProps = {
    show: boolean
    title: BaseModalTitleType
    modalBody: JSX.Element
    toggle: () => void
    onSubmit?: () => void
    submitText?: string
}

const BaseModal = ({ show, title, modalBody, toggle, onSubmit }: ModalProps) => {
    const [t] = useTranslation('translation')
    const titlestr = t(title.title)
    const iconstr = t(title.icon)
    const desstr = t(title.description)

    const handleSubmit = () => {
        if (onSubmit) onSubmit()
    }

    return (
        <Modal
            show={show}
            keyboard={true}
            onHide={toggle}
            size={title.size}
            style={{ zIndex: 1055}}
            centered={true}
        >
            <Modal.Header closeButton>
                <span>
                    <i className={iconstr} />
                    {titlestr}
                    {
                        desstr.length > 0 
                        ? <PopoverLabel name='train-description'>{desstr}</PopoverLabel>
                        : null
                    }
                    
                </span>
            </Modal.Header>
            <Modal.Body>
                {modalBody}
            </Modal.Body>
            <Modal.Footer>
                {
                    title.submitText.length > 0 && (
                        <Button variant='success' className='btn-rounded btn-sm' onClick={handleSubmit}>
                            {t(title.submitText)}
                        </Button>
                    )
                }
                <Button variant='dark' className='btn-rounded btn-sm' onClick={toggle}>
                    {t('button.close')}
                </Button>
            </Modal.Footer>
        </Modal>
    )
}

export default BaseModal