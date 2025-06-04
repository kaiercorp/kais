import { useEffect, useState, useContext } from "react"
import { Button, Card, Col, Row } from "react-bootstrap"
import { useTranslation } from "react-i18next"
import { useLocation, Link } from "react-router-dom"
import styled from 'styled-components'

import { TableClsTrialConfigsTable, PerfTable } from "features"
import { logger, objDeepCopy, ApiFetchTrial, ApiDownloadFile } from "helpers"
import { ModelsArea, TableErrorGraphImage, TestModelCard, TrainModelContainer } from "features/TrainModel"
import { FeatureImportanceChart } from "features/Chart"
import { LocationContext } from 'contexts'
import { useQueryClient } from '@tanstack/react-query'
import { useSocket } from 'hooks'
import { engine } from 'appConstants/trial'

const CardHeaderLeft = styled.div`
float: left;
color: #ffffff;
font-weight: 600;
`

const CardHeaderRight = styled.div`
float: right;
`
const TableREGTestDetail = () => {
    const [t] = useTranslation('translation')
    const location = useLocation()

    const [trialId, setTrialId] = useState<number | undefined>()
    const [prevLocation, setPrevLocation] = useState(location.pathname)
    const [selectedModel, setSelectedModel] = useState<any>()
    const [socketConnected, setSocketConnected] = useState(false)
    const [errorGraph, setErrorGraph] = useState<any>()

    const { updateLocationContextValue } = useContext(LocationContext)

    const { trial } = ApiFetchTrial(trialId) 
    const queryClient = useQueryClient()

    useEffect(() => {
        const pathVariables = location.pathname.split('/')
        setPrevLocation(pathVariables.slice(0, pathVariables.length - 2).join('/'))
        setTrialId(Number(pathVariables.slice(-2, -1)[0]))

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const [models, setModels] = useState<any[]>()
    useEffect(() => {
        if (!trial || trial.project_id === 0) {
            setModels([])
            return
        }

        logger.log(`Change Location to ${t('title.table.reg.test', { trial: trial.trial_name })}`)
        updateLocationContextValue({location: 'table.reg.test', locationValue: { trial: `[${trial.trial_id}] ${trial.trial_name}` }})

        if (trial.test && trial.test.models) {
            let _models = objDeepCopy(trial.test.models)

            setModels(_models)
            setSelectedModel(_models[0])
        }

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trial])

    const handleSocketMessage = (e: MessageEvent<any>) => {
        try {
            if (e.data === 'Invalid request') {
                throw e
            }

            let msg = JSON.parse(e.data)
            if (msg.hasOwnProperty('error_graph')) {
                setErrorGraph(msg['error_graph'])
            } else {
                msg.forEach((m: any) => {
                    m.probs = JSON.parse(m.data)
                })
            }
        } catch (e) {
            logger.error(e)
        }
    }

    const ws = useSocket(`/trials/test/result/${trial?.test?.id}`, 'Result Data', handleSocketMessage, {setSocketConnected, shouldCleanup: true, shouldConnect: !!selectedModel})

    useEffect(() => {
        if (!socketConnected || !ws || !ws.current) return  

        ws.current.send(JSON.stringify({ sample_id: selectedModel.id, engine_type: engine.table_reg }))   
        // eslint-disable-next-line react-hooks/exhaustive-deps 
    }, [socketConnected])

    const onSelectModel = (model: any) => {
        setSelectedModel(model)
        if (!socketConnected || !ws || !ws.current) return
        ws.current.send(JSON.stringify({ model_name: model.name }))
    }

    const handleDownloadTestFiles = () => {
        if (!selectedModel) return

        ApiDownloadFile(queryClient, `/trial/test/download/${selectedModel.id}`, `${trial.trial_name}_${selectedModel.model}.zip`)
    }

    return (
        <Col>
            <Row>
                <Col>
                    <Card>
                        <Card.Header>
                            <CardHeaderLeft>{t('ui.train.title.info')}</CardHeaderLeft>
                            <CardHeaderRight><Link to={prevLocation}>{t('button.go.list')}</Link></CardHeaderRight>
                        </Card.Header>
                        <Card.Body>
                            {
                                trial&&trial.parent_trial&&<TableClsTrialConfigsTable trial={trial} config={trial.parent_trial.params} />
                            }
                        </Card.Body>
                    </Card>
                </Col>
            </Row>

            <Row>
                <Col>
                    <Card>
                        <TrainModelContainer isVertical={false}>
                            <ModelsArea isVertical={false}>
                                <div style={{display:'flex'}}>
                                {models && models.map((model: any) => {
                                    return (
                                        <TestModelCard
                                            key={`model-list-${model.id}`}
                                            model={model}
                                            isSelected={selectedModel && (model.id === selectedModel.id)}
                                            onClick={() => onSelectModel(model)}
                                        />
                                    )
                                })}
                                </div>
                            </ModelsArea>
                        </TrainModelContainer>
                    </Card>
                    <Card>
                        <Card.Body>
                            {
                                selectedModel && (
                                    <Row style={{ marginBottom: '10px' }}>
                                        <Col sm={2}>
                                            <TrainModelContainer isVertical={false}>
                                                <PerfTable perfStr={selectedModel.perf.String} title={t('ui.test.title.perf', { perf: selectedModel.name })} isTestDetail={true} />
                                            </TrainModelContainer>
                                            <Row style={{marginTop:'10px'}}>
                                                <Col sm={2}></Col>
                                                <Col><Button style={{width:'140px', height:'30px'}} onClick={handleDownloadTestFiles}>{t('button.download.resultfile')}</Button></Col>
                                                <Col sm={2}></Col>
                                            </Row>
                                        </Col>
                                        <Col sm={10}>
                                            <TableErrorGraphImage error_graph={errorGraph} />
                                        </Col>
                                    </Row>
                                )
                            }
                            {/* <Row>
                                <Col>
                                    {resultList && <TestDataResult resultList={resultList} onClick={onSelectData} />}
                                </Col>
                                <Col>
                                    
                                </Col>
                            </Row> */}
                            <Row>
                                {selectedModel&&<FeatureImportanceChart feature_importance={selectedModel.chart.String} />}
                            </Row>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </Col>
    )
}

export default TableREGTestDetail