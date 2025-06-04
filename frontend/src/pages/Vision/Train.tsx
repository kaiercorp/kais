import { forwardRef, useEffect, useImperativeHandle, useState, useContext } from 'react'
import { Button, Card, Col, Form, InputGroup, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'

import { LabelInput, LabelSelect } from 'components'
import { ClsClassConfigModal, RadioGPU, SelectDataModal, SelectImageResolution, TargetMetric, TargetMetricMLCls, TargetMetricAD } from 'features'
import { useToggle } from 'hooks'
import { checkOnlyNumber, objDeepCopy, ApiFetchDatasets, ApiCreateTrain } from 'helpers'
import { emptyErrors as clsEmptyErrors, validateCLS } from 'pages/Vision/Cls.SL/TrainValidator'
import { emptyErrors as adEmptyErrors, validateAD } from 'pages/Vision/AD/TrainValidator'
import { engine } from 'appConstants/trial'
import { ProjectContext, TrialContext } from 'contexts'
import { emptySLClsTrial, emptyMLClsTrial, emptyVADTrial, CreateTrainType } from 'common'


const StyledRow = styled(Row)`
height: 200px;
margin-top: 15px;
`
const Container = styled(Col)`
border: 1px solid #4a525d;
padding: 0;
margin: 0 12px;
overflow: auto;
`

const LabelContainer  = styled(Container)`
border: 1px solid #4a525d;
padding: 0;
margin: 0 12px;
overflow: auto;
height: 250px;
`

const TrainVision = forwardRef(({ toggle, engineType }: any, ref) => {
  const [t] = useTranslation('translation')
  const [formErrors, setFormErrors] = useState<clsEmptyErrors | adEmptyErrors>({ hasError: false })
  const [directoryId, setDirectoryId] = useState<number | undefined>()

  const { projectContextValue } = useContext(ProjectContext)
  const { trialContextValue, updateTrialContextValue } = useContext(TrialContext)

  const { classes } = ApiFetchDatasets('image', directoryId, engineType)
  const createTrain = ApiCreateTrain()

  const [isDataModalOpened, toggleDataModal] = useToggle()
  const [isClassConfigModalOpened, toggleClassConfigModal] = useToggle()
  useImperativeHandle(ref, () => ({
    handleSubmit
  }))

  const requestData = trialContextValue.requestData ? trialContextValue.requestData : objDeepCopy(engineType === engine.vision_ad ? emptyVADTrial : engine.vision_cls_ml ? emptyMLClsTrial : emptySLClsTrial)
  const trainType = trialContextValue.trainMode || 'auto'

  useEffect(() => {
    if (!classes) return
    let class_list = objDeepCopy(classes.classes)
    let classList = new Array<any>()
    class_list.forEach((c: string) => {
      let classObj: any = {}
      classObj[c] = 1
      classList.push(classObj)
    })
    handleRequestData('class_list', classList)

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [classes])

  const handleDirectoryIdChange = (directoryId: number) => {
    setDirectoryId(directoryId)
  }
  
  const handleRequestData = (key: string, value: any) => {
      if (engineType === engine.vision_ad) {
        handleADRequestData(key, value)
      } else {
        handleClsRequestData(key, value)
      }
  }

  const handleClsRequestData = (key: string, value: any) => {
    let newRequestData = objDeepCopy(requestData)

    if (key === 'resolution') {
      if (value === '') return
      let size = value.split('x')
      newRequestData.train_config.width = Number(size[0]) || 1
      newRequestData.train_config.height = Number(size[1]) || 1
    } else if (['class_list', 'width', 'height', 'auto_stop'].includes(key)) {
      newRequestData.train_config[key] = value
    } else if (key === 'tiff_frame_number') {
      newRequestData.train_config[key] = Number(value)
    } else if (key.startsWith('target_metric')) {
      const ks = value.split('.')
      newRequestData.train_config.target_metric = engineType === engine.vision_cls_sl ? 
              {'wa': 0, 'uwa': 0, 'precision': 0, 'recall': 0, 'f1': 0} 
            : {'image_accuracy': 0, 'image_precision': 0, 'image_recall': 0, 'image_f1_score': 0, 'label_accuracy': 0, 'label_precision': 0, 'label_recall': 0, 'label_f1_score': 0}
      newRequestData.train_config.target_metric[ks[1]] = 100
    }  else {
      newRequestData[key] = value
    }

    updateTrialContextValue({requestData: newRequestData})
  }
  
  const handleADRequestData = (key: string, value: any) => {
      let newRequestData = objDeepCopy(requestData)

      if (key === 'resolution') {
          if (value === '') return
          let size = value.split('x')
          newRequestData.train_config.width = Number(size[0]) || 1
          newRequestData.train_config.height = Number(size[1]) || 1
      } else if (['default_config_file', 'class_list', 'width', 'height', 'auto_stop'].includes(key)) {
          newRequestData.train_config[key] = value
      } else if (key === 'tiff_frame_number') {
          newRequestData.train_config[key] = Number(value)
      } else if (key === 'base_lr') {
          let val = value
          if (val < 0.0000000001) val = 0.0000000001
          if (val > 1) val = 1
          newRequestData.train_config.base_lr = val
      } else if (['epochs'].includes(key)) {
          newRequestData.train_config[key] = checkOnlyNumber(value)
      } else if (['save_top_k', 'train_batch_size'].includes(key)) {
          if (value === '-') {
              newRequestData.train_config[key] = value
          } else {
              newRequestData.train_config[key] = checkOnlyNumber(value)
          }

      } else if (key.startsWith('target_metric')) {
          const ks = value.split('.')
          newRequestData.train_config.target_metric = { 'wa': 0, 'uwa': 0, 'precision': 0, 'recall': 0, 'f1': 0, 'auroc': 0, 'prauc': 0 }
          newRequestData.train_config.target_metric[ks[1]] = 100
      } else {
          newRequestData[key] = value
      }

      updateTrialContextValue({requestData: newRequestData})
  }

  const handleRequestDataset = (path: string, id: number) => {
    let newRequestData = objDeepCopy(requestData)

    newRequestData.data_path = path
    newRequestData.test_db = t('train.config.test_db', { value: path })
    newRequestData.trial_name = t(`train.config.trial_name.${trainType}`, { value: path })
    newRequestData['dataset_id'] = id

    updateTrialContextValue({requestData: newRequestData})
  }

  const setConfig = (): CreateTrainType => {
    let newRequestData = objDeepCopy(requestData)

    newRequestData.project_id = projectContextValue.selectedProject.project_id

    let class_list = newRequestData.train_config.class_list
    let labels = new Array<any>()
    class_list.map((c: any) => {
      for (var key in c) {
        labels.push(c[key])
      }
      return labels
    })
    newRequestData.train_config.class_list = labels
    newRequestData.train_config.width = Number(newRequestData.train_config.width)
    newRequestData.train_config.height = Number(newRequestData.train_config.height)
    
    if ( engineType === engine.vision_cls_ml ) {
      newRequestData.train_config.target_metric.image_accuracy = newRequestData.train_config.target_metric.image_accuracy /100
      newRequestData.train_config.target_metric.image_precision = newRequestData.train_config.target_metric.image_precision /100
      newRequestData.train_config.target_metric.image_recall = newRequestData.train_config.target_metric.image_recall /100
      newRequestData.train_config.target_metric.image_f1_score = newRequestData.train_config.target_metric.image_f1_score /100
      newRequestData.train_config.target_metric.label_accuracy = newRequestData.train_config.target_metric.label_accuracy /100
      newRequestData.train_config.target_metric.label_precision = newRequestData.train_config.target_metric.label_precision /100
      newRequestData.train_config.target_metric.label_recall = newRequestData.train_config.target_metric.label_recall /100
      newRequestData.train_config.target_metric.label_f1_score = newRequestData.train_config.target_metric.label_f1_score /100
    } else {
      newRequestData.train_config.target_metric.wa = newRequestData.train_config.target_metric.wa /100
      newRequestData.train_config.target_metric.uwa = newRequestData.train_config.target_metric.uwa /100
      newRequestData.train_config.target_metric.precision = newRequestData.train_config.target_metric.precision /100
      newRequestData.train_config.target_metric.recall = newRequestData.train_config.target_metric.recall /100
      newRequestData.train_config.target_metric.f1 = newRequestData.train_config.target_metric.f1 /100
    }

    if ( engineType === engine.vision_ad ) {
      newRequestData.train_config.target_metric.auroc = newRequestData.train_config.target_metric.auroc / 100
      newRequestData.train_config.target_metric.prauc = newRequestData.train_config.target_metric.prauc / 100
    }

    return {
      project_id: projectContextValue.selectedProject.project_id as number,
      trial_name: newRequestData.trial_name,
      test_db: newRequestData.test_db,
      dataset_id: newRequestData.dataset_id,
      train_type: trainType,
      data_type: 'image',
      engine_type: engineType,
      trial_id: newRequestData.trial_id,
      params: JSON.stringify(newRequestData)
    }
  }

  const handleSubmit = () => {
    const errors = engineType === engine.vision_ad ? validateAD(trainType, requestData, t) : validateCLS(trainType, requestData, t) 
    setFormErrors(errors)

    if (!errors.hasError) {
      createTrain.mutate(setConfig())
      return true
    }

    return false
  }

  const handleSelectData = (path: string, id: number) => {
    handleRequestDataset(path, id)
  }

  const handleClassConfig = (value: number[]) => {
    handleRequestData('class_list', value)
  }


  return (
    <Form noValidate validated={formErrors.hasError}>
      <Row>
        <Col sm={engineType === engine.vision_cls_ml ? 6 : 12} xs={`${trainType === 'auto' ? 12 : 6}`}>
          <Card>
            <Card.Header>{t('ui.train.title.common')}</Card.Header>
            <Card.Body>
              <Form.Group>
                <Row>
                  <Form.Label column='sm' sm={4}>
                    {t('ui.train.data_path')}
                  </Form.Label>
                  <Col>
                    <InputGroup className='mb-1'>
                      <Form.Control value={requestData.data_path} readOnly />
                      <Button variant='info' onClick={toggleDataModal}>
                        {t('button.select')}
                      </Button>
                    </InputGroup>
                  </Col>
                </Row>
                <Row>


                  {classes && classes.is_tiff && (
                    <LabelSelect
                      title={
                        <span>
                          {t('ui.train.tiffframe')}
                        </span>
                      }
                      name={'tiff_frame_number'}
                      onChange={(e: any) => handleRequestData('tiff_frame_number', e.target.value)}
                      value={engineType === engine.vision_ad ? 
                        requestData.train_config.tiff_frame_number 
                        : requestData.train_config.target_metric
                      }
                    >
                      <option key='frame-0' value={0}>
                        Frame 0
                      </option>
                      <option key='frame-1' value={1}>
                        Frame 1
                      </option>
                      <option key='frame-2' value={2}>
                        Frame 2
                      </option>
                    </LabelSelect>
                  )}
                </Row>
                { engineType !== engine.vision_ad && 
                  ( trainType === 'manual' && requestData.train_config.class_list.length > 0 ? (
                    <Row>
                      <Form.Label column='sm' sm={4}>
                        {t('ui.formatter.classratio')}
                      </Form.Label>
                      <Col>
                        <InputGroup className='mb-1'>
                          <Form.Control
                            value={requestData.train_config.class_list
                              .map((label: string) => label[(Object.keys(label) as (keyof typeof label)[])[0]])
                              .join(', ')}
                            readOnly
                          />
                          <Button variant='info' onClick={toggleClassConfigModal}>
                            {t('button.set')}
                          </Button>
                        </InputGroup>
                      </Col>
                    </Row>
                  ) : null)
                }
              </Form.Group>

              <RadioGPU selectGPU={handleRequestData} errors={formErrors} />

              { engineType === engine.vision_ad?  
                <TargetMetricAD requestData={requestData} handleRequestData={handleRequestData} />
                : engineType === engine.vision_cls_ml ?
                <TargetMetricMLCls requestData={requestData} handleRequestData={handleRequestData} />
                :<TargetMetric requestData={requestData} handleRequestData={handleRequestData} />
              }
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

              { engineType !== engine.vision_ad &&
                <SelectImageResolution
                  width={requestData.train_config.width}
                  height={requestData.train_config.height}
                  errors={formErrors}
                  selectResolution={handleRequestData}
                />
              }

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
        </Col>

        {engineType === engine.vision_cls_ml && 
          <Col xs={12} sm={6}>
              <Card>
                  <Card.Header>{t('ui.train.title.label.information')}</Card.Header>
                  <Card.Body>

                  <Form.Group className={'mb-1'}>
                    <Row>
                        <Form.Label column='sm' sm={4}>
                        {t('ui.train.title.label.count')}
                        </Form.Label>
                        <Container>
                          {classes?.classes.length}
                        </Container>
                    </Row>
                      <StyledRow>
                        <Form.Label column='sm' sm={4}>
                        {t('ui.train.title.label.list')}
                        </Form.Label>
                        <LabelContainer>
                            {classes?.classes?.sort().map((label: string) =>
                              <div>{label}</div>
                            )}
                        </LabelContainer>
                      </StyledRow>
                   </Form.Group>
                  </Card.Body>
              </Card>
          </Col>
        }
      </Row>

      <SelectDataModal show={isDataModalOpened} selectData={handleSelectData} toggle={toggleDataModal} isTest={false} dataType='image' directoryId={directoryId} onDirectoryIdChange={handleDirectoryIdChange} />
      <ClsClassConfigModal
        show={isClassConfigModalOpened}
        handleData={handleClassConfig}
        data={requestData.train_config.class_list}
        toggle={toggleClassConfigModal}
      />
    </Form>
  )
})

export default TrainVision