import { useEffect, useRef, useState, useContext } from 'react'
import { Button, Card, Col, Form, Nav, Row, Tab } from 'react-bootstrap'
import { Link, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'
import { useQueryClient } from '@tanstack/react-query'

import { LocationContext } from 'contexts'
import { useToggle, useSocket } from 'hooks'
import { logger, ApiFetchTrial, ApiDownloadFile } from 'helpers'
import { APICore } from 'helpers/api/apiCore'
import { engine } from 'appConstants/trial'

import { RegModelList, TableClsTrialConfigsTable, HyperParameterHistory, LogHistory, TableClsFIContainer, MultiSampleTest, BaseModal } from 'features'
import { BaseModalTitleType, emptyBaseModalTitle, tableRegMultiTest } from 'common'

const CardHeaderLeft = styled.div`
    float: left;
    color: #ffffff;
    font-weight: 600;
`

const CardHeaderRight = styled.div`
    float: right;
`

const CardHeaderButton = styled(Button)`
    float: right;
    margin-right: 10px;
    padding: 2px 5px;
    font-size: 11px;
`

const api = new APICore()

const TableREGTrainDetail = () => {
    const [t] = useTranslation('translation')
    const location = useLocation()

    const [prevLocation, setPrevLocation] = useState(location.pathname)
    const [trains, setTrains] = useState<any[]>()
    const [eventKey, setEventKey] = useState<string>('tabs-summary')
    const [modalTitle, setModalTitle] = useState<BaseModalTitleType>(emptyBaseModalTitle)
    const [modalBody, setModalBody] = useState<JSX.Element>(<></>)
    const [trialId, setTrialId] = useState<number | undefined>()

    const childComponentRef = useRef<any>()

    const { updateLocationContextValue } = useContext(LocationContext)

    const [showModal, toggleModal, openModal] = useToggle()

    const queryClient = useQueryClient()
    const user = api.getLoggedInUser()
    const { trial } = ApiFetchTrial(trialId)

    useEffect(() => {
        const pathVariables = location.pathname.split('/')
        setPrevLocation(pathVariables.slice(0, pathVariables.length - 2).join('/'))
        setTrialId(Number(pathVariables.slice(-2, -1)[0]))

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        if (!trial || trial.project_id === 0) return

        logger.log(`Change Location to ${t('title.table.reg.train', { trial: trial.trial_name })}`)
        updateLocationContextValue({ location: 'table.reg.train', locationValue: { trial: `[${trial.trial_id}] ${trial.trial_name}` } })

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trial])

    const handleSocketMessage = (e: MessageEvent<any>) => {
        try {
            let msg = JSON.parse(e.data)
            setTrains(msg)
        } catch (e) {
            logger.error(e)
        }
    }
    useSocket(`/trials/trains/${trial?.trial_id}`, 'Trains', handleSocketMessage, { shouldCleanup: true, shouldConnect: !!trial && trial.trial_id && trial.trial_id !== 0 })

    const onSubmit = () => {
        if (childComponentRef) {
            let result = childComponentRef.current.handleSubmit()
            if (result) toggleModal()
        } else {
            toggleModal()
        }
    }

    const openMultiTestModal = (row: any) => {
        setModalTitle(tableRegMultiTest)
        setModalBody(<MultiSampleTest ref={childComponentRef} selectedTrial={trial} selectedModels={row} engineType={engine.table_reg} />)
        openModal()
    }

    const handleDownloadModel = (row: any) => {
        ApiDownloadFile(queryClient, `/trial/model/download/${row.train_id}/${row.name}`, `${trial.trial_name}_${row.name}.zip`)
    }

    const handleDownloadReport = (e:any) => {
        e.stopPropagation()
        ApiDownloadFile(queryClient, `/trial/report/download/${trial.trial_id}`, `${trial.trial_id}_${trial.trial_name}_report.zip`)
    }

    return (
        <Form>
            <Col>
                <Row>
                    <Col>
                        <Card>
                            <Card.Header>
                                <CardHeaderLeft>{t('ui.train.title.info')}</CardHeaderLeft>
                                <CardHeaderRight><Link to={prevLocation}>{t('button.go.list')}</Link></CardHeaderRight>
                                {trial && trial.state === 'finish' && <CardHeaderButton onClick={openMultiTestModal}>{t('modal.title.table.reg.multitest.title')}</CardHeaderButton>}
                                {user && <CardHeaderButton variant='success' onClick={handleDownloadReport}>Save Report</CardHeaderButton>}
                            </Card.Header>
                            <Card.Body>
                                <TableClsTrialConfigsTable trial={trial} config={trial?.params} showPerfTable={true}/>
                            </Card.Body>
                        </Card>
                    </Col>
                </Row>

                <Row>
                    <Col>
                        <Card>
                            <Tab.Container defaultActiveKey='tabs-summary'>
                                <Nav variant='tabs' className='nav-bordered' as='ul'>
                                    <Nav.Item as='li'>
                                        <Nav.Link as={Link} to='#' eventKey='tabs-summary' onClick={() => setEventKey('tabs-summary')}>
                                            <i className='d-md-none d-block me-1' />
                                            <span className='d-none d-md-block'>{t('ui.train.title.chart')}</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item as='li'>
                                        <Nav.Link as={Link} to='#' eventKey='tabs-models' onClick={() => setEventKey('tabs-models')}>
                                            <i className='d-md-none d-block me-1' />
                                            <span className='d-none d-md-block'>{t('ui.train.model.title')}</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    {
                                        user &&
                                        <Nav.Item as='li'>
                                            <Nav.Link as={Link} to='#' eventKey='tabs-hphistory' onClick={() => setEventKey('tabs-hphistory')}>
                                                <i className='d-md-none d-block me-1' />
                                                <span className='d-none d-md-block'>{t('ui.train.title.hphistory')}</span>
                                            </Nav.Link>
                                        </Nav.Item>
                                    }
                                    {
                                        user &&
                                        <Nav.Item as='li'>
                                            <Nav.Link as={Link} to='#' eventKey='tabs-loghistory' onClick={() => setEventKey('tabs-loghistory')}>
                                                <i className='d-md-none d-block me-1' />
                                                <span className='d-none d-md-block'>{t('ui.train.title.loghistory')}</span>
                                            </Nav.Link>
                                        </Nav.Item>
                                    }
                                </Nav>

                                <Tab.Content>
                                    <Tab.Pane eventKey='tabs-summary' id='tabs-summary'>
                                        {trial && trial.trial_id > 0 && <TableClsFIContainer trial={trial} trains={trains} />}
                                    </Tab.Pane>
                                    <Tab.Pane eventKey='tabs-models' id='tabs-models'>
                                        {eventKey === 'tabs-models' && trial && trial.trial_id > 0 && <RegModelList trial={trial} openMultiTest={openMultiTestModal} downloadModel={handleDownloadModel} />}
                                    </Tab.Pane>
                                    <Tab.Pane eventKey='tabs-hphistory' id='tabs-hphistory'>
                                        {eventKey === 'tabs-hphistory' && trial && trial.trial_id > 0 && <HyperParameterHistory trial={trial} />}
                                    </Tab.Pane>
                                    <Tab.Pane eventKey='tabs-loghistory' id='tabs-loghistory'>
                                        {eventKey === 'tabs-loghistory' && trial && trial.trial_id > 0 && <LogHistory trial={trial} />}
                                    </Tab.Pane>
                                </Tab.Content>
                            </Tab.Container>
                        </Card>
                    </Col>
                </Row>
                <BaseModal show={showModal} title={modalTitle} modalBody={modalBody} onSubmit={onSubmit} toggle={toggleModal} />
            </Col>
        </Form>
    )
}

export default TableREGTrainDetail