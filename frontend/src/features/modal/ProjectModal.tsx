/*eslint-disable*/
import { useEffect, useState } from 'react'
import { Button, FloatingLabel, Form, Modal } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'

import { trimAndRemoveDoubleSpace } from 'helpers'
import { CustomConfirm } from 'components'
import {  } from 'helpers'
import { ProjectType } from 'common'

type ModalProps = {
  show: boolean
  isEdit: boolean
  data: any
  category: string,
  projects: any,
  handleSubmit: ({}:any) => void,
  toggle: () => void
}

const ProjectModal = ({ show, isEdit, data, category, projects, handleSubmit, toggle }: ModalProps) => {
  const [t] = useTranslation('translation')

  const [invalidName, setInvalidName] = useState(false)
  const [btnDisabled, setBtnDisabled] = useState(false)
  const [projectName, setProjectName] = useState('')
  const [projectDescription, setProjectDescription] = useState('')

  useEffect(() => {
    if (show) {
      setProjectName(data.project_name)
      setProjectDescription(data.description)
    }
  }, [show])

  const checkName = (name: string) => {
    return (
      projects ? projects.filter((project: ProjectType) => project.project_name?.toUpperCase() === name.toUpperCase())
        .length > 0 : false
    )
  }

  const changeProjectName = (e: any) => {
    e.stopPropagation()
    let newName = e.target.value

    setProjectName(newName)
    const isDisabled = checkName(newName)

    setBtnDisabled(isDisabled || newName === '')
    setInvalidName(isDisabled)
  }

  const changeProjectDescription = (e: any) => {
    e.stopPropagation()
    setProjectDescription(e.target.value)
  }

  const onHandleProject = async () => {
    if (btnDisabled) return

    if (isEdit) {
      onEditProject()
    } else {
      onCreateProject()
    }
  }

  const onCreateProject = async () => {
    setBtnDisabled(true)
    const trimedProjectName = trimAndRemoveDoubleSpace(projectName)

    CustomConfirm({
      onConfirm: () => {
        handleSubmit({
          project_name: trimedProjectName,
          description: projectDescription,
          category
        })
        setBtnDisabled(false)
        toggle()
      },
      onCancel: () => {
        setBtnDisabled(false)
      },
      message: t('ui.confirm.add', { new: trimedProjectName }),
    })
  }

  const onEditProject = async () => {
    setBtnDisabled(true)
    const trimedProjectName = trimAndRemoveDoubleSpace(projectName)

    CustomConfirm({
      onConfirm: () => {
        handleSubmit({ 
          project_id: data.project_id,
          project_name: trimedProjectName,
          description: projectDescription
        })
        setBtnDisabled(false)
        toggle()
      },
      onCancel: () => {
        setBtnDisabled(false)
      },
      message: t('ui.confirm.edit'),
    })
  }

  return (
    <Modal show={show} keyboard={true} onHide={toggle}>
      <Modal.Header closeButton>{isEdit ? t('ui.project.edit') : t('ui.project.add')}</Modal.Header>
      <Modal.Body>
        <FloatingLabel controlId='projectName' label='Project Name' className='mb-1'>
          <Form.Control type='text' isInvalid={invalidName} onChange={changeProjectName} value={projectName || ''} />
          <Form.Control.Feedback type='invalid' tooltip>
            {t('ui.info.project.exist')}
          </Form.Control.Feedback>
        </FloatingLabel>
        <FloatingLabel controlId='projectDescription' label='Project Description' className='mb-1'>
          <Form.Control type='text' onChange={changeProjectDescription} value={projectDescription || ''} />
        </FloatingLabel>
      </Modal.Body>
      <Modal.Footer>
        <Button variant='success' className='btn-rounded btn-sm' onClick={onHandleProject} disabled={btnDisabled}>
          {isEdit ? t('button.change') : t('button.add')}
        </Button>
        <Button variant='dark' className='btn-rounded btn-sm' onClick={toggle}>
          {t('button.cancel')}
        </Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ProjectModal
