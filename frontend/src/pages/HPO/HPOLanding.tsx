import { useEffect, useState, useContext } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'

import { logger, objDeepCopy } from 'helpers'
import { Button, Card, Col, Form, Nav, Row, Tab } from 'react-bootstrap'
import { LocationContext } from 'contexts'

import { ApiFetchHpos, ApiInitHpos, ApiUpdateHpos } from 'helpers/api'
import { ButtonArea, CardHeaderLeft, CardHeaderRight } from 'components/Containers'
import HPCardList from './HPCardList'
import { HPOModelType, HPOParamsDistType, HPOParamsType } from 'common'

const HPOLanding = () => {
    const [t] = useTranslation('translation')
    const { updateLocationContextValue } = useContext(LocationContext)

    const { hpo } = ApiFetchHpos()
    const updateConfigs = ApiUpdateHpos()
    const forceInitConfigs = ApiInitHpos()

    useEffect(() => {
        logger.log(`Change Location to ${t('title.hpo')}`)
        updateLocationContextValue({ location: 'hpo' })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const [inits, setInits] = useState<any[]>([])

    useEffect(() => {
        initializeConfigs()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [hpo])

    const initializeConfigs = () => {
        if (!hpo || hpo.length < 1) return

        let initval = objDeepCopy(hpo)
        setInits(initval)
    }

    const onForceInit = () => {
        forceInitConfigs.mutate()
    }

    const emptyModel: HPOModelType = {
        engine_type: "",
        model_name: "MODEL",
        class_type: "both",
        is_use: true,
        params: [
            {
                name: "PARAM",
                suggest_type: "category",
                data_type: "int",
                is_use: true,
                dists: [
                    {
                        dist: "dist,param",
                        is_use: true,
                        cond: { key: '', value: '', operator: '' }
                    }
                ]
            }
        ]
    }

    const emptyParam: HPOParamsType = {
        name: "PARAM",
        suggest_type: "category",
        data_type: "int",
        is_use: true,
        dists: [
            {
                dist: "dist,param",
                is_use: true,
                cond: { key: '', value: '', operator: '' }
            }
        ]
    }

    const emptyDist: HPOParamsDistType = {
        dist: "dist,param",
        is_use: true,
        cond: { key: '', value: '', operator: '' }
    }

    const onChangeItem = (e: any, engine:string, op: string, type: string, idx_model: number, idx_param: number, idx_dist: number, key: string, value: any) => {
        e.stopPropagation()

        let newModels = inits.slice().filter((hp: any) => hp.engine_type === engine)

        if (type === "model") {
            if (op === "add") {
                let newModel = objDeepCopy(emptyModel)
                newModel.engine_type = engine
                newModels.push(newModel)
            } else if (op === "delete") {
                newModels.splice(idx_model, 1)
            } else if (op === "edit") {
                newModels[idx_model][key] = value
            }
        } else if (type === "param") {
            if (op === "add") {
                let newParam = objDeepCopy(emptyParam)
                newModels[idx_model].params.push(newParam)
            } else if (op === "delete") {
                newModels[idx_model].params.splice(idx_param, 1)
            } else if (op === "edit") {
                newModels[idx_model].params[idx_param][key] = value
            }
        } else if (type === "dist") {
            if (op === "add") {
                let newDist = objDeepCopy(emptyDist)
                newModels[idx_model].params[idx_param].dists.push(newDist)
            } else if (op === "delete") {
                newModels[idx_model].params[idx_param].dists.splice(idx_dist, 1)
            } else if (op === "edit") {
                if (key.includes("cond.")) {
                    const k = key.replace("cond.", "")
                    newModels[idx_model].params[idx_param].dists[idx_dist].cond[k] = value
                } else {
                    newModels[idx_model].params[idx_param].dists[idx_dist][key] = value
                }
            }
        } 
        newModels = newModels.concat(inits.slice().filter((hp: any) => hp.engine_type !== engine))

        setInits(newModels)
    }

    const onSaveConfigs = () => {
        let newConfigs = objDeepCopy(inits)

        updateConfigs.mutate(newConfigs)
    }

    const [eventKey, setEventKey] = useState<string>('tabs-vcls')

    return (
        <Form>
            <Col>
                <Row>
                    <Card>
                        <Card.Header>
                            <CardHeaderLeft>Manage Hyper Parameters</CardHeaderLeft>
                            <CardHeaderRight><Button variant={'danger'} onClick={onForceInit}>Force Init</Button></CardHeaderRight>
                        </Card.Header>
                        <Card.Body>
                            <Row>
                                <ButtonArea>
                                    <Button variant={'danger'} onClick={initializeConfigs}>{t('button.cancel')}</Button>
                                    <Button onClick={onSaveConfigs}>{t('button.save')}</Button>
                                </ButtonArea>
                            </Row>
                            <Tab.Container defaultActiveKey='tabs-vcls'>
                                <Nav variant='tabs' className='nav-borderd' as='ul'>
                                    <Nav.Item as='li'>
                                        <Nav.Link as={Link} to='#' eventKey={'tabs-vcls'} onClick={() => setEventKey('tabs-vcls')}>
                                            <i className='d-md-none d-block me-1' />
                                            <span className='d-none d-md-block'>V.CLS</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item as='li'>
                                        <Nav.Link as={Link} to='#' eventKey={'tabs-vad'} onClick={() => setEventKey('tabs-vad')}>
                                            <i className='d-md-none d-block me-1' />
                                            <span className='d-none d-md-block'>V.AD</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item as='li'>
                                        <Nav.Link as={Link} to='#' eventKey={'tabs-tcls'} onClick={() => setEventKey('tabs-tcls')}>
                                            <i className='d-md-none d-block me-1' />
                                            <span className='d-none d-md-block'>T.CLS</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item as='li'>
                                        <Nav.Link as={Link} to='#' eventKey={'tabs-treg'} onClick={() => setEventKey('tabs-treg')}>
                                            <i className='d-md-none d-block me-1' />
                                            <span className='d-none d-md-block'>T.REG</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item as='li'>
                                        <Nav.Link as={Link} to='#' eventKey={'tabs-tsad'} onClick={() => setEventKey('tabs-tsad')}>
                                            <i className='d-md-none d-block me-1' />
                                            <span className='d-none d-md-block'>TS.AD</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item as='li'>
                                        <Nav.Link as={Link} to='#' eventKey={'tabs-tsfc'} onClick={() => setEventKey('tabs-tsfc')}>
                                            <i className='d-md-none d-block me-1' />
                                            <span className='d-none d-md-block'>TS.FC</span>
                                        </Nav.Link>
                                    </Nav.Item>
                                </Nav>

                                <Tab.Content>
                                    <Tab.Pane eventKey='tabs-vcls' id='tabs-vcls'>
                                        {eventKey === 'tabs-vcls' && inits && HPCardList({ hps: inits.filter((hp: any) => hp.engine_type === 'vision-cls'), engine: 'vision-cls', onChange: onChangeItem })}
                                    </Tab.Pane>
                                    <Tab.Pane eventKey='tabs-vad' id='tabs-vad'>
                                        {eventKey === 'tabs-vad' && inits && HPCardList({ hps: inits.filter((hp: any) => hp.engine_type === 'vision-ad'), engine: 'vision-ad', onChange: onChangeItem })}
                                    </Tab.Pane>
                                    <Tab.Pane eventKey='tabs-tcls' id='tabs-tcls'>
                                        {eventKey === 'tabs-tcls' && inits && HPCardList({ hps: inits.filter((hp: any) => hp.engine_type === 'table-cls'), engine: 'table-cls', onChange: onChangeItem })}
                                    </Tab.Pane>
                                    <Tab.Pane eventKey='tabs-treg' id='tabs-treg'>
                                        {eventKey === 'tabs-treg' && inits && HPCardList({ hps: inits.filter((hp: any) => hp.engine_type === 'table-reg'), engine: 'table-reg', onChange: onChangeItem })}
                                    </Tab.Pane>
                                    <Tab.Pane eventKey='tabs-tsad' id='tabs-tsad'>
                                        {eventKey === 'tabs-tsad' && inits && HPCardList({ hps: inits.filter((hp: any) => hp.engine_type === 'ts-ad'), engine: 'ts-ad', onChange: onChangeItem })}
                                    </Tab.Pane>
                                    <Tab.Pane eventKey='tabs-tsfc' id='tabs-tsfc'>
                                        {eventKey === 'tabs-tsfc' && inits && HPCardList({ hps: inits.filter((hp: any) => hp.engine_type === 'ts-fc'), engine: 'ts-fc', onChange: onChangeItem })}
                                    </Tab.Pane>
                                </Tab.Content>
                            </Tab.Container>
                        </Card.Body>
                    </Card>
                </Row>
            </Col>
        </Form>
    )
}

export default HPOLanding