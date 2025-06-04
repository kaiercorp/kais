import { useEffect, useState, useContext } from "react"
import { Button, Card, Col, Row } from "react-bootstrap"
import { useTranslation } from "react-i18next"
import { useLocation } from "react-router-dom"
import styled from 'styled-components'

import { TableClsTrialConfigsTable, PerfTable } from "features"
import { logger, objDeepCopy, ApiFetchTrial, ApiDownloadFile } from "helpers"
import { Link } from "react-router-dom"
import { CFMatrix, ModelsArea, TestModelCard, TrainModelContainer } from "features/TrainModel"

import { FeatureImportanceChart } from "features/Chart"
import { LocationContext } from 'contexts'
import { useQueryClient } from '@tanstack/react-query'
import { useSocket } from "hooks"

const CardHeaderLeft = styled.div`
float: left;
color: #ffffff;
font-weight: 600;
`

const CardHeaderRight = styled.div`
float: right;
`

const TableCLSTestDetail = () => {
    const [t] = useTranslation('translation')
    const location = useLocation()

    const [prevLocation, setPrevLocation] = useState(location.pathname)
    const [selectedModel, setSelectedModel] = useState<any>()
    const [trialId, setTrialId] = useState<number | undefined>()

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

        logger.log(`Change Location to ${t('title.vision.cls.test', { trial: trial.trial_name })}`)
        updateLocationContextValue({location: 'vision.cls.test', locationValue: { trial: `[${trial.trial_id}] ${trial.trial_name}` }})

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
            if (msg.hasOwnProperty('origin')) {
                
            } else {
                msg.forEach((m: any) => {
                    m.probs = JSON.parse(m.data)
                })
            }
        } catch (e) {
            logger.error(e)
        }
    }

    useSocket(`/trials/test/result/${trial?.test?.id}`, 'Result Data', handleSocketMessage, {shouldCleanup: true, shouldConnect: !!selectedModel})

    const onSelectModel = (model: any) => {
        setSelectedModel(model)
    }

    const handleDownloadTestFiles = () => {
        ApiDownloadFile(queryClient, `/trial/test/download/${selectedModel.id}`, `${trial.trial_name}_${selectedModel.name}.zip`)
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
                                                <PerfTable perfStr={selectedModel.perf.String} title={t('ui.test.title.perf', { perf: selectedModel.name })} isTestDetail={true}/>
                                            </TrainModelContainer>
                                            <Row style={{marginTop:'10px'}}>
                                                <Col sm={2}></Col>
                                                <Col><Button style={{width:'140px', height:'30px'}} onClick={handleDownloadTestFiles}>{t('button.download.resultfile')}</Button></Col>
                                                <Col sm={2}></Col>
                                            </Row>
                                        </Col>
                                        <Col sm={10} style={{paddingRight: '24px'}}>
                                            <CFMatrix cf_matrixStr={selectedModel.cf_matrix.String} />
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

export default TableCLSTestDetail