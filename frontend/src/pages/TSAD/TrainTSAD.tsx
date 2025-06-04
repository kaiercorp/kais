import { engine } from "appConstants/trial"
import { LabelInput, LabelSelect, PopoverLabel, StatusRow, LabelSelectTyped } from "components"
import { RadioGPUWithCPU, SelectDataModal, TableCLSColumnsDnd, TargetMetricReg, TimestampSelector } from "features"
import { objDeepCopy, ApiFetchDatasets, ApiCreateTrain } from "helpers"
import { useToggle } from "hooks"
import { forwardRef, useEffect, useImperativeHandle, useState, useContext } from "react"
import { Button, Card, Col, Form, InputGroup, Row } from "react-bootstrap"
import { useTranslation } from "react-i18next"
import { tableIndexTypes, CreateTrainType, emptyTSADRequest } from 'common'
import { emptyErrors, validateTSAD } from "./TrainValidator"
import { ProjectContext, TrialContext } from 'contexts'

const TrainTSAD = forwardRef(({ toggle }: any, ref) => {
    const [t] = useTranslation('translation')

    /**
     * 전체 컬럼 = 데이터 번호 열 + 예측 대상 열 + 학습 입력 열 + 학습 제외 열
     */
    //데이터 번호 열
    const [indexOptions, setIndexOptions] = useState<any>([])
    const [selectedIndex, setSelectedIndex] = useState<any>([])
    // 예측 대상 열
    const [dateOptions, setDateOptions] = useState<any>([])
    const [selectedDate, setSelectedDate] = useState<any>([])
    // 학습 입력 열
    const [selectedInclude, setSelectedInclude] = useState<string[]>([])
    // 학습 제외 열
    const [selectedExclude, setSelectedExclude] = useState<string[]>([])

    const [formErrors, setFormErrors] = useState<emptyErrors>({ hasError: false })
    const [directoryId, setDirectoryId] = useState<number | undefined>()

    const { projectContextValue } = useContext(ProjectContext)
    const { trialContextValue, updateTrialContextValue } = useContext(TrialContext)

    const { columns } = ApiFetchDatasets('table', directoryId)
    const createTrain = ApiCreateTrain()

    const [isDataModalOpened, toggleDataModal] = useToggle()
    useImperativeHandle(ref, () => ({
        handleSubmit
    }))

    const requestData = trialContextValue.requestData ? trialContextValue.requestData : objDeepCopy(emptyTSADRequest)
    const trainType = trialContextValue.trainMode || 'auto'

    useEffect(() => {
        // 데이터셋 변경되면 열 선택 초기화
        setSelectedIndex([])
        setSelectedDate([])
        setSelectedInclude([])
        setSelectedExclude([])

        if (!columns || columns.length < 1) {
            return
        }

        // tableIndexTypes에 정의된 것과 같은 컬럼을 데이터 번호 열로 자동 선택
        let _indexCols: any = []
        let _cols = columns.slice()
        _cols.forEach((col: string) => {
            _indexCols.push({ value: col, label: col })
            if (tableIndexTypes && tableIndexTypes.includes(col.toLowerCase())) {
                setSelectedIndex([{ value: col, label: col }])
            }
        })
        setIndexOptions(_indexCols)

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [columns])

    useEffect(() => {
        if (!columns || columns.length < 1) {
            setIndexOptions([])
            return
        }

        let _newOptions: any = []
        columns.forEach((col: string) => {
            if (selectedDate.filter((y: any) => y['value'] === col).length > 0) {
                return false
            }
            if (selectedExclude.includes(col)) {
                return false
            }
            _newOptions.push({ value: col, label: col })
        })
        setIndexOptions(_newOptions)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedDate, selectedInclude])

    useEffect(() => {
        if (!columns) {
            setDateOptions([])
            return
        }

        let _newOptions: any = []
        columns.forEach((col: string) => {
            if (selectedIndex.filter((y: any) => y['value'] === col).length > 0) {
                return false
            }
            if (selectedExclude.includes(col)) {
                return false
            }
            _newOptions.push({ value: col, label: col })
        })
        setDateOptions(_newOptions)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedIndex, selectedInclude])

    useEffect(() => {
        if (!columns) {
            setDateOptions([])
            return
        }

        let _includeCols = columns.filter((col: string) => {
            if (selectedIndex.filter((y: any) => y['value'] === col).length > 0) {
                return false
            }
            if (selectedDate.filter((y: any) => y['value'] === col).length > 0) {
                return false
            }
            if (selectedExclude.includes(col)) {
                return false
            }
            return true
        })
        setSelectedInclude(_includeCols)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [columns, selectedIndex, selectedDate, selectedExclude])

    const onChangeIndexCol = (option: Record<string, Object>[]) => {
        setSelectedIndex(option.slice(-1))
    }

    const onChangeDateCol = (option: Record<string, Object>[]) => {
        setSelectedDate(option.slice(-1))
    }

    const onChangeIncludeCol = (cols: string[]) => {
        let _newCols: string[] = []
        selectedExclude.forEach((selected: string) => {
            if (cols.includes(selected)) {
                _newCols.push(selected)
            }
        })

        setSelectedExclude(cols)
    }

    const handleRequestData = (key: string, value: any) => {
        let newRequestData = objDeepCopy(requestData)
        if ([
            'epochs',
            'save_top_k',
            'train_batch_size',
            'auto_stop',
            'date_format',
            'threshold_mode'
        ].includes(key)
        ) {
            newRequestData.train_config[key] = value
        } else if (key === 'threshold_k') {
            if (/^$|^\d+\.$|^\d+(\.\d{1,5})?$/g.test(value)) {
                newRequestData.train_config[key] = value
            }
        } else if (key.startsWith('target_metric')) {
            const ks = value.split('.')
            newRequestData.train_config.target_metric = { 'mse': 0, 'rmse': 0, 'mae': 0 }
            newRequestData.train_config.target_metric[ks[1]] = 100
        } else {
            newRequestData[key] = value
        }

        updateTrialContextValue({ requestData: newRequestData })
    }

    const handleRequestDataset = (path: string, id: number) => {
        let newRequestData = objDeepCopy(requestData)

        newRequestData.data_path = path
        newRequestData.test_db = t('train.config.test_db', { value: path })
        newRequestData.trial_name = t(`train.config.trial_name.${trainType}`, { value: path })
        newRequestData['dataset_id'] = id

        updateTrialContextValue({ requestData: newRequestData })
    }

    const setConfig = (): CreateTrainType => {
        let newRequestData = objDeepCopy(requestData)

        if (selectedIndex.length > 0) newRequestData.train_config['index_column'] = selectedIndex[0]['value']
        else newRequestData.train_config['index_column'] = 'none'
        newRequestData.train_config['date_column'] = selectedDate[0]?.value || ''
        newRequestData.train_config['input_column'] = selectedInclude

        if (newRequestData.train_config['index_column'] === '') {
            newRequestData.train_config['index_column'] = 'none'
        }

        if (newRequestData.gpus[0] === 'none') {
            newRequestData.gpus = []
        }

        newRequestData.train_config.target_metric.mse = newRequestData.train_config.target_metric.mse / 100
        newRequestData.train_config.target_metric.rmse = newRequestData.train_config.target_metric.rmse / 100
        newRequestData.train_config.target_metric.mae = newRequestData.train_config.target_metric.mae / 100

        return {
            ...newRequestData,
            project_id: projectContextValue.selectedProject.project_id as number,
            trial_name: newRequestData.trial_name,
            dataset_id: newRequestData.dataset_id,
            train_type: trainType,
            data_type: 'table',
            engine_type: engine.ts_ad,
            trial_id: newRequestData.trial_id,
            params: JSON.stringify(newRequestData)
        }
    }

    const handleSubmit = () => {
        const _reqData = setConfig()
        const errors = validateTSAD(trainType, _reqData, t)

        setFormErrors(errors)

        if (!errors.hasError) {
            createTrain.mutate(_reqData)
            return true
        }

        return false
    }

    const handleDirectoryIdChange = (directoryId: number) => {
        setDirectoryId(directoryId)
    }

    return (
        <Form noValidate validated={formErrors.hasError}>
            <Row>
                <Col xs={12} sm={6}>
                    <Card>
                        <Card.Header>{t('ui.train.title.common')}</Card.Header>
                        <Card.Body>
                            <Form.Group>
                                <Row>
                                    <Form.Label column='sm' sm={4}>
                                        {t('ui.train.data_path')}
                                    </Form.Label>
                                    <Col>
                                        {formErrors && formErrors['data_path'] !== '' ? (
                                            <Form.Control.Feedback type='invalid' className='d-block'>
                                                {formErrors['data_path']}
                                            </Form.Control.Feedback>
                                        ) : null}
                                        <InputGroup className='mb-1'>
                                            <Form.Control value={requestData.data_path} readOnly />
                                            <Button variant='info' onClick={toggleDataModal}>
                                                {t('button.select')}
                                            </Button>
                                        </InputGroup>
                                    </Col>
                                </Row>
                            </Form.Group>

                            <RadioGPUWithCPU selectGPU={handleRequestData} errors={formErrors} />

                            <LabelInput
                                title={t('ui.train.name')}
                                name='trial_name'
                                value={requestData ? requestData.trial_name : ''}
                                onChange={(e: any) => handleRequestData('trial_name', e.target.value)}
                                errors={formErrors}
                            />

                            <LabelInput
                                title={t('ui.train.name.testdbname')}
                                name='test_db'
                                value={requestData ? requestData.test_db : ''}
                                onChange={(e: any) => handleRequestData('test_db', e.target.value)}
                                errors={formErrors}
                            />

                            <TargetMetricReg requestData={requestData} handleRequestData={handleRequestData} />

                            {trainType === 'auto' && (
                                <Form.Group className={'mb-1'}>
                                    <Row>
                                        <Form.Label column='sm' sm={4}>{t('ui.train.autostop')}</Form.Label>
                                        <Col column='sm' sm={8} style={{ marginTop: '5px' }}>
                                            <Form.Switch
                                                type='switch'
                                                checked={requestData.train_config.auto_stop}
                                                label={<Form.Label>{requestData.train_config.auto_stop ? t('ui.train.autostop.auto') : t('ui.train.autostop.user')}</Form.Label>}
                                                onChange={() => handleRequestData('auto_stop', !requestData.train_config.auto_stop)}
                                            />
                                        </Col>
                                    </Row>
                                </Form.Group>
                            )}
                        </Card.Body>
                    </Card>

                    <Card>
                        <Card.Header>
                            {t('ui.train.threshold')}
                            <PopoverLabel>
                                <StatusRow><span>{t(`ui.train.threshold.desc`)}</span></StatusRow>
                            </PopoverLabel>
                        </Card.Header>
                        <Card.Body>
                            <LabelSelect
                                title={<span>{t('ui.train.threshold.threshold_mode')}<PopoverLabel>
                                    <StatusRow><span>{t(`ui.train.threshold.threshold_mode.desc`)}</span></StatusRow>
                                </PopoverLabel></span>}
                                name={'threshold_mode'}
                                onChange={(e: any) => handleRequestData('threshold_mode', e.target.value)}
                                value={requestData.train_config.threshold_mode}
                                required={false}
                            >
                                <option key='threshold_mode_ratio' value='percentile'>
                                    {t('ui.train.threshold.threshold_mode.ratio')}
                                </option>
                                <option key='threshold_mode_std' value='std'>
                                    {t('ui.train.threshold.threshold_mode.std')}
                                </option>
                                <option key='threshold_mode_custom' value='custom'>
                                    {t('ui.train.threshold.threshold_mode.custom')}
                                </option>
                            </LabelSelect>

                            <LabelInput
                                title={t('ui.train.threshold.k')}
                                name='threshold_k'
                                value={requestData.train_config.threshold_k}
                                onChange={(e: any) => handleRequestData('threshold_k', e.target.value)}
                                errors={formErrors}
                            />
                        </Card.Body>
                    </Card>
                </Col>

                <Col xs={12} sm={6}>
                    <Card>
                        <Card.Header>{t('ui.train.title.column')}</Card.Header>
                        <Card.Body>
                            <LabelSelectTyped
                                title={t('ui.train.indexCol')}
                                name={'index_column'}
                                options={indexOptions}
                                onChange={(option: Record<string, Object>[]) => onChangeIndexCol(option)}
                                value={selectedIndex}
                                required={true}
                                errors={formErrors}
                            />
                            <LabelSelectTyped
                                title={t('ui.train.dateCol')}
                                name={'date_column'}
                                options={dateOptions}
                                onChange={(option: Record<string, Object>[]) => onChangeDateCol(option)}
                                value={selectedDate}
                                required={false}
                                errors={formErrors}
                            />
                            <TimestampSelector
                                onChange={handleRequestData}
                                value={requestData.train_config.date_format}
                            />
                            <TableCLSColumnsDnd
                                input_columns={selectedInclude}
                                except_columns={selectedExclude}
                                onMoveItem={(cols: string[]) => onChangeIncludeCol(cols)}
                            />
                        </Card.Body>
                    </Card>
                </Col>
            </Row>

            <SelectDataModal show={isDataModalOpened} selectData={handleRequestDataset} toggle={toggleDataModal} isTest={false} dataType='table' directoryId={directoryId} onDirectoryIdChange={handleDirectoryIdChange} />
        </Form>
    )
})

export default TrainTSAD