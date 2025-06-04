import { useEffect, useState, useContext } from 'react'
import { useTranslation } from 'react-i18next'

import { logger, objDeepCopy } from 'helpers'
import { Button, Card, Col, Form, Row } from 'react-bootstrap'
import { CustomConfirm } from 'components'
import { LocationContext } from 'contexts'

import { ApiFetchDRs, ApiDeleteDR, ApiDeleteDS } from 'helpers/api'
import { DatasetRootType } from 'common'
import { ButtonArea, CardHeaderLeft } from 'components/Containers'
import classNames from 'classnames'
import styled from 'styled-components'
import DatasetModal from './DatasetModal'

const StyledButtonArea = styled.div`

    & .btn-icon {
        padding: 0 5px;
        font-size: 15px;
    }

    & button {
        margin-left: 5px;
    }
`

const StyledDSArea = styled(Row)`
    padding-left: 20px;
    margin-top: 5px;
`

const DatasetLanding = () => {
    const [t] = useTranslation('translation')
    const { updateLocationContextValue } = useContext(LocationContext)

    useEffect(() => {
        logger.log(`Change Location to ${t('title.dataset')}`)
        updateLocationContextValue({ location: 'dataset' })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const { datasetroots } = ApiFetchDRs()
    const deleteDR = ApiDeleteDR()
    const deleteDS = ApiDeleteDS()

    const [inits, setInits] = useState<any[]>([])

    useEffect(() => {
        initializeDRs()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [datasetroots])

    const initializeDRs = () => {
        if (!datasetroots) return

        let initval = objDeepCopy(datasetroots)
        setInits(initval)
    }

    const closeModal = () => {
        setModal({
            ...modal,
            show: false,
        })
    }

    const [modal, setModal] = useState({
        show: false,
        isEdit: false,
        data: {},
        toggle: closeModal,
    })

    const toggleAdd = () => {
        setModal({
            ...modal,
            show: true,
            isEdit: false,
            data: {},
        })
    }

    const toggleEdit = (selected: DatasetRootType) => {
        setModal({
            ...modal,
            show: true,
            isEdit: true,
            data: selected || {},
        })
    }

    const onDeleteDR = (e: any, selected: DatasetRootType) => {
        e.stopPropagation()

        if (!selected) return

        CustomConfirm({
            onConfirm: () => {
                deleteDR.mutate(Number(selected.id))
            },
            onCancel: () => { },
            message: t('ui.confirm.delete.dataroot', { name: selected.name })
        })
    }

    const onDeleteDS = (e: any, selected: any) => {
        e.stopPropagation()

        if (!selected) return

        CustomConfirm({
            onConfirm: () => {
                deleteDS.mutate(Number(selected.id))
                
            },
            onCancel: () => { },
            message: t('ui.confirm.delete.dataset', { name: selected.name })
        })
    }

    return (
        <Form>
            <Card>
                <Card.Header>
                    <CardHeaderLeft>Manage Dataset root driectories</CardHeaderLeft>
                </Card.Header>
                <Card.Body>
                    {
                        inits && inits.map((dr: DatasetRootType, index: number) => {
                            return (
                                <Row key={`datasetroot-${index}`}>
                                    <Col>
                                        <Row>
                                            <Col>
                                                <Row>
                                                    <Col sm={4}>
                                                        <Form.Label>Name</Form.Label>
                                                    </Col>
                                                    <Col sm={8}>
                                                        <Form.Control
                                                            required
                                                            type={'text'}
                                                            size='sm'
                                                            value={dr.name}
                                                            disabled={true}
                                                        />
                                                    </Col>
                                                </Row>
                                            </Col>
                                            <Col>
                                                <Row>
                                                    <Col sm={4}>
                                                        <Form.Label>Path</Form.Label>
                                                    </Col>
                                                    <Col sm={8}>
                                                        <Form.Control
                                                            required
                                                            type={'text'}
                                                            size='sm'
                                                            value={dr.path}
                                                            disabled={true}
                                                        />
                                                    </Col>
                                                </Row>
                                            </Col>
                                            <Col>
                                                <Form.Group className={'mb-3'}>
                                                    <Form.Check
                                                        type={'checkbox'}
                                                        label={"Is Use"}
                                                        name={"is_use"}
                                                        checked={dr.is_use}
                                                        disabled={true}
                                                    />
                                                </Form.Group>
                                            </Col>
                                            <Col>
                                                <StyledButtonArea>
                                                    <Button variant={'primary'} className="btn-icon" onClick={() => toggleEdit(dr)}>
                                                        <i className={classNames('mdi', 'mdi-pencil')} /> edit
                                                    </Button>
                                                    <Button variant={'danger'} className="btn-icon" onClick={(e) => onDeleteDR(e, dr)}>
                                                        <i className={classNames('mdi', 'mdi-trash-can')} /> delete
                                                    </Button>
                                                </StyledButtonArea>
                                            </Col>
                                        </Row>

                                        {
                                            dr.datasets && dr.datasets.map((ds: any, jndex: number) => {
                                                return (
                                                    <StyledDSArea key={`datasetroot-${index}-dataset-${jndex}`}>
                                                        <Col sm={3}>
                                                            <Form.Control
                                                                required
                                                                type={'text'}
                                                                size='sm'
                                                                value={ds.name}
                                                                disabled={true}
                                                            />
                                                        </Col>
                                                        <Col sm={4}>
                                                            <Form.Control
                                                                required
                                                                type={'text'}
                                                                size='sm'
                                                                value={ds.path}
                                                                disabled={true}
                                                            />
                                                        </Col>
                                                        <Col sm={2}>
                                                            <Button variant={'danger'} className="btn-icon" onClick={(e) => onDeleteDS(e, ds)}>
                                                                <i className={classNames('mdi', 'mdi-trash-can')} /> delete
                                                            </Button>
                                                        </Col>
                                                    </StyledDSArea>
                                                )
                                            })
                                        }


                                    </Col>
                                </Row>
                            )
                        })
                    }
                </Card.Body>
            </Card>
            <Row>
                <ButtonArea>
                    <Button variant={'danger'} onClick={initializeDRs}>{t('button.cancel')}</Button>
                    <Button className="btn-icon" onClick={toggleAdd}>
                        <i className={classNames('mdi', 'mdi-plus')} /> <span>Add Path</span>
                    </Button>
                </ButtonArea>
            </Row>
            <DatasetModal {...modal} />
        </Form>
    )
}

export default DatasetLanding