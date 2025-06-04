/*eslint-disable*/
import { useEffect, useState, useContext } from 'react'
import { Button, FloatingLabel, Form, Modal } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'

import { trimAndRemoveDoubleSpace } from 'helpers'
import { CustomConfirm } from 'components'
import { LocationContext } from 'contexts'
import { ApiCreateDR, ApiUpdateDR } from 'helpers'
import { useQueryClient } from '@tanstack/react-query'
import { QUERY_KEY } from 'helpers'
import { ProjectType } from 'common'

type ModalProps = {
    show: boolean
    isEdit: boolean
    data: any
    toggle: () => void
}

const DatasetModal = ({ show, isEdit, data, toggle }: ModalProps) => {
    const [t] = useTranslation('translation')
    const { locationContextValue } = useContext(LocationContext)

    const insertDR = ApiCreateDR()
    const updateDR = ApiUpdateDR()

    const queryClient = useQueryClient()

    const [invalidName, setInvalidName] = useState(false)
    const [btnDisabled, setBtnDisabled] = useState(true)
    const [name, setName] = useState('')
    const [path, setPath] = useState('')
    const [isUse, setIsUse] = useState(true)

    useEffect(() => {
        if (show) {
            setName(data.name||'')
            setPath(data.path||'')
            setIsUse(data.is_use||true)
        }
    }, [show])

    const changeName = (e: any) => {
        e.stopPropagation()
        let newName = e.target.value

        setName(newName)
    }

    const changePath = (e: any) => {
        e.stopPropagation()
        setPath(e.target.value)
    }

    const changeIsUse = (e: any) => {
        e.stopPropagation()
        setIsUse(!isUse)
    }

    const onHandleProject = async () => {
        if (isEdit) {
            onEditDR()
        } else {
            onCreateDR()
        }
    }

    const onCreateDR = async () => {
        setBtnDisabled(true)
        const trimedName = trimAndRemoveDoubleSpace(name)

        CustomConfirm({
            onConfirm: () => {
                insertDR.mutate({ name: name, path: path, is_use: isUse })
                setBtnDisabled(false)
                toggle()
            },
            onCancel: () => {
                setBtnDisabled(false)
            },
            message: t('ui.confirm.add', { new: trimedName }),
        })
    }

    const onEditDR = async () => {
        setBtnDisabled(true)
        const trimedName = trimAndRemoveDoubleSpace(name)

        CustomConfirm({
            onConfirm: () => {
                updateDR.mutate({ id: data.id, name: name, path: path, is_use: isUse })
                setBtnDisabled(false)
                toggle()
            },
            onCancel: () => {
                setBtnDisabled(false)
            },
            message: t('ui.confirm.edit', { old: data.name, new: trimedName }),
        })
    }

    return (
        <Modal show={show} keyboard={true} onHide={toggle}>
            <Modal.Header closeButton>{isEdit ? t('ui.dataroot.edit') : t('ui.dataroot.add')}</Modal.Header>
            <Modal.Body>
                <FloatingLabel controlId='name' label='Name' className='mb-1'>
                    <Form.Control type='text' isInvalid={invalidName} onChange={changeName} value={name || ''} />
                </FloatingLabel>
                <FloatingLabel controlId='path' label='Path' className='mb-1'>
                    <Form.Control type='text' onChange={changePath} value={path || ''} />
                </FloatingLabel>
                <Form.Group className={'mb-3'}>
                    <Form.Check
                        type={'checkbox'}
                        label={"Is Use"}
                        name={"is_use"}
                        checked={isUse}
                        onChange={changeIsUse}
                    />
                </Form.Group>
            </Modal.Body>
            <Modal.Footer>
                <Button variant='success' className='btn-rounded btn-sm' onClick={onHandleProject}>
                    {isEdit ? t('button.change') : t('button.add')}
                </Button>
                <Button variant='dark' className='btn-rounded btn-sm' onClick={toggle}>
                    {t('button.cancel')}
                </Button>
            </Modal.Footer>
        </Modal >
    )
}

export default DatasetModal
