import { useEffect, useRef, useState, useContext, useCallback } from 'react'
import { Card, Col, Row } from 'react-bootstrap'

import { useToggle, useSocket } from 'hooks'
import { GPUArea, BaseModal, TrialTable, FilterTrials, ColumnModeling } from 'features'
import { customFilterTrials, customTrialSort, logger, objDeepCopy } from 'helpers'
import { useTranslation } from 'react-i18next'
import OneClickButton from 'components/Button/OneClickButton'
import { 
    BaseModalTitleType, 
    compareTrials, 
    emptyBaseModalTitle, 
    filterTrials, 
    tsADAdditionalTrain, 
    tsADAutoTrain, 
    tsADMultiTest,
    emptyTSADRequest,  
    ProjectType 
} from 'common'

import { useLocation } from 'react-router-dom'
import TrainTSAD from './TrainTSAD'
import TSADMultiSampleTest from './TSADMultiSampleTest'
import TrainTSADAdditional from './TrainTSADAdditional'
import CompareTrial from './CompareTrial'
import { LocationContext, ProjectContext, TrialContext, FilterContext } from 'contexts'
import { useQueryClient } from '@tanstack/react-query'
import { QUERY_KEY } from 'helpers'


const TSADLanding = () => {
    const [t] = useTranslation('translation')
    const location = useLocation()

    const { updateLocationContextValue } = useContext(LocationContext)
    const { updateProjectContextValue } = useContext(ProjectContext)
    const { trialContextValue, updateTrialContextValue } = useContext(TrialContext)
    const { filterContextValue } = useContext(FilterContext)
    const { trials } = trialContextValue
    const { filter, useFilter } = filterContextValue

    const queryClient = useQueryClient()
    const projects = queryClient.getQueryData<ProjectType[]>([`${QUERY_KEY.projects}_ts.ad`])

    const handleSocketMessage = useCallback((e: MessageEvent<any>) => {
        try {
		if (!e.data) return
          let trials = JSON.parse(e.data)
          
          trials.forEach(function (t: any) {
            if (t.params) {
              t.params = JSON.parse(t.params)
            }
            if (t.params_parent && t.params_parent !== '') {
              t.params_parent = JSON.parse(t.params_parent)
            }
            if (t.perf && t.perf !== '') {
              t.perf = JSON.parse(t.perf)
              if (t.train_type === 'test') {
                if (t.params_parent.train_config) t.perf['target_metric'] = t.params_parent.train_config.target_metric
              } else {
                t.perf['target_metric'] = t.target_metric
              }
            }
          })
          
          updateTrialContextValue({ trials })
        } catch (error) {
          logger.error(`Error parsing: ${error}`)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [trialContextValue])
    useSocket('/trials/' + location.pathname.split('/').pop(), 'Trials', handleSocketMessage, { shouldCleanup: true, shouldConnect: (location.pathname.split('/').pop() !== undefined) })

    useEffect(() => {
        logger.log(`Change Location to ${t('title.ts.ad.trials')}`)
        updateLocationContextValue({location: 'ts.ad.trials'})

        const pathVariables = location.pathname.split('/')
        const prjID = Number(pathVariables[pathVariables.length - 1])
        const selectedProject = projects?.filter((project: ProjectType) => Number(project.project_id) === prjID)
        updateProjectContextValue({ selectedProject: selectedProject ? selectedProject[0] : {project_id: prjID} })

        updateTrialContextValue({selectedRows: []})
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [t])

    const [showModal, toggleModal, openModal] = useToggle()
    const [modalTitle, setModalTitle] = useState<BaseModalTitleType>(emptyBaseModalTitle)
    const [modalBody, setModalBody] = useState<JSX.Element>(<></>)

    const childComponentRef = useRef<any>()
    const onSubmit = () => {
        if (childComponentRef) {
            let result = childComponentRef.current.handleSubmit()
            if (result) toggleModal()
        } else {
            toggleModal()
        }
    }

    const openAutoTrainModal = () => {
        updateTrialContextValue({trainMode: 'auto', requestData: objDeepCopy(emptyTSADRequest)})

        setModalTitle(tsADAutoTrain)
        setModalBody(<TrainTSAD ref={childComponentRef} />)
        openModal()
    }

    const openMultiTestModal = () => {
        setModalTitle(tsADMultiTest)
        setModalBody(<TSADMultiSampleTest ref={childComponentRef} />)
        openModal()
    }

    const openCompareModal = () => {
      setModalTitle(compareTrials)
      setModalBody(<CompareTrial ref={childComponentRef} />)
      openModal()
    }

    const openAdditionalTrainModal = (row: any) => {
        updateTrialContextValue({trainMode: 'additional', requestData: objDeepCopy(emptyTSADRequest)})

        setModalTitle(tsADAdditionalTrain)
        setModalBody(<TrainTSADAdditional trialId={row.trial_id} ref={childComponentRef} />)
        openModal()
    }

    const [filteredTrials, setFilteredTrials] = useState<any[]>(trials? trials : [])
    useEffect(() => {
        if (!trials) return

        let newTrials = trials.map((trial: any) => {
            if (trial.state === 'fail' && trial.best_model_download_path) {
                trial.state = 'finish-fail'
            }
            trial['state-search'] = t(`state.${trial.state}`)
            return trial
        })

        newTrials = customTrialSort(newTrials)

        if (useFilter) {
            newTrials = customFilterTrials(filter, newTrials)
        }

        setFilteredTrials(newTrials)
    }, [trials, filter, useFilter, t])

    const openFilterModal = () => {
        setModalTitle(filterTrials)
        setModalBody(<FilterTrials ref={childComponentRef} />)
        openModal()
    }

    return (
        <>
            <Row>
                <Col xs={12} md={12} lg={8} xl={8} xxl={4}>
                    <Row>
                        <Col>
                            <Card.Body className='p-0'>
                                <Row className='g-0'>
                                    <OneClickButton
                                        title={t('oneclick.autotrain')}
                                        subTitle={t('oneclick.sub.train')}
                                        icon={'mdi mdi-cursor-pointer'}
                                        onClick={openAutoTrainModal}
                                    />
                                    <OneClickButton
                                        title={t('oneclick.foldertest')}
                                        subTitle={t('oneclick.sub.test')}
                                        icon={'mdi mdi-flask'}
                                        onClick={openMultiTestModal}
                                    />
                                </Row>
                            </Card.Body>
                        </Col>
                    </Row>
                </Col>
                <Col xs={12} md={12} lg={4} xl={4} xxl={8}>
                    <GPUArea />
                </Col>
            </Row>

            <Row style={{marginTop: '5px'}}>
                <TrialTable
                    CustomColumn={ColumnModeling}
                    filteredTrials={filteredTrials}
                    openCompareModal={openCompareModal}
                    openFilterModal={openFilterModal}
                    openAdditionalTrainModal={openAdditionalTrainModal}
                    hideAcc={true}
                />
            </Row>

            <BaseModal show={showModal} title={modalTitle} modalBody={modalBody} onSubmit={onSubmit} toggle={toggleModal} />
        </>
    )
}

export default TSADLanding