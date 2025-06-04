import { useState } from 'react'
import { Button, ButtonGroup, Card, Col, Form, Modal, Row } from 'react-bootstrap'
import Nouislider from 'nouislider-react'
import 'nouislider/distribute/nouislider.css'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

import { initialFilter } from 'common'
import { HyperDatepicker, PopoverLabel, StatusLabel } from 'components'
import { stateColors } from 'appConstants'
import { objDeepCopy } from 'helpers'

type ModalProps = {
  show: boolean
  data: any
  onConfirm: (data: any) => void
  toggle: () => void
}

const StyledCol = styled(Col)`
  & button {
    margin-right: 5px;
  }
`

const StyledSwitch = styled(Form.Check)`
  & :checked {
    background-color: ${(props) => props.bgcolor && props.bgcolor};
    border-color: ${(props) => props.bgcolor && props.bgcolor};
  }
`

const StatusSwitch = ({ state, checked, onChange }: any) => {
  return (
    <Col xs={2}>
      <StyledSwitch
        bgcolor={stateColors[state]}
        type='switch'
        id={`switch-${state}`}
        checked={checked}
        label={<StatusLabel htmlFor={`switch-${state}`} state={state} />}
        onChange={onChange}
      />
    </Col>
  )
}

const ClsFilterModal = ({ show, data, onConfirm, toggle }: ModalProps) => {
  const [t] = useTranslation('translation')
  const [filter, setFilter] = useState(data)

  const onChangeFilter = (key: string, value: any) => {
    let newFilter = objDeepCopy(filter)
    newFilter[key] = value
    setFilter(newFilter)
  }

  const onChangeRange = (key: string, value: any) => {
    let newFilter = objDeepCopy(filter)
    newFilter[key].min = Number(value[0].toFixed(0))
    newFilter[key].max = Number(value[1].toFixed(0))
    setFilter(newFilter)
  }

  const onChangeStatus = (key: string) => {
    let newFilter = objDeepCopy(filter)
    newFilter.state[key] = !filter.state[key]
    newFilter.state.total =
      Object.entries(filter.state).filter((state) => {
        return state[0] !== 'total' && state[0] !== 'finish-fail' && state[1] === false
      }).length < 1
    setFilter(newFilter)
  }

  const onChangeTotal = (e: any) => {
    const checked = e.target.checked
    let newFilter = objDeepCopy(filter)
    for (const [key] of Object.entries(newFilter.state)) {
      newFilter.state[key] = checked
    }

    setFilter(newFilter)
  }

  const onStartDate24Hour = () => {
    let newFilter = objDeepCopy(filter)
    let startDate = new Date(filter.endDate.getTime() - 24 * 60 * 60 * 1000)
    newFilter.startDate = startDate
    setFilter(newFilter)
  }

  const onStartDate7Days = () => {
    let newFilter = objDeepCopy(filter)
    let startDate = new Date(filter.endDate.getTime() - 7 * 24 * 60 * 60 * 1000)
    newFilter.startDate = startDate
    setFilter(newFilter)
  }

  const onStartDateMonth = (value: number) => {
    let newFilter = objDeepCopy(filter)
    const endDate = filter.endDate
    let startDate = new Date(endDate.getTime())
    startDate.setDate(1)
    startDate.setMonth(startDate.getMonth() - value)
    startDate.setDate(
      Math.min(endDate.getDate(), new Date(startDate.getFullYear(), startDate.getMonth() + 1, 0).getDate())
    )
    newFilter.startDate = startDate
    setFilter(newFilter)
  }

  const onCloseModal = () => {
    onConfirm(filter)
    toggle()
  }

  const onInitFilter = () => {
    setFilter(initialFilter)
  }

  return (
    <Modal show={show} size='xl' keyboard={true} onHide={toggle}>
      <Modal.Header closeButton>{t('ui.filter.title')}</Modal.Header>
      <Modal.Body>
        <Card>
          <h5>상태</h5>
          <Row>
            <StatusSwitch state='total' checked={filter.state.total} onChange={onChangeTotal} />
            <Col xs={10}>
              <Row>
                <StatusSwitch state='train' checked={filter.state.train} onChange={() => onChangeStatus('train')} />
                <StatusSwitch
                  state='additional_train'
                  checked={filter.state.additional_train}
                  onChange={() => onChangeStatus('additional_train')}
                />
                <StatusSwitch
                  state='finish'
                  checked={filter.state.finish}
                  onChange={() => onChangeStatus('finish')}
                />
                <StatusSwitch
                  state='cancel'
                  checked={filter.state.cancel}
                  onChange={() => onChangeStatus('cancel')}
                />
                <StatusSwitch state='fail' checked={filter.state.fail} onChange={() => onChangeStatus('fail')} />
              </Row>
              <Row>
                <StatusSwitch state='test' checked={filter.state.test} onChange={() => onChangeStatus('test')} />
                <StatusSwitch
                  state='finish_test'
                  checked={filter.state.finish_test}
                  onChange={() => onChangeStatus('finish_test')}
                />
                <StatusSwitch
                  state='test_cancel'
                  checked={filter.state.test_cancel}
                  onChange={() => onChangeStatus('test_cancel')}
                />
                <StatusSwitch
                  state='test_fail'
                  checked={filter.state.test_fail}
                  onChange={() => onChangeStatus('test_fail')}
                />
                <StatusSwitch
                  state='idle'
                  checked={filter.state.idle}
                  onChange={() => onChangeStatus('idle')}
                />
              </Row>
            </Col>
          </Row>
        </Card>

        <Card>
          <h5>{t('ui.filter.date.title')}</h5>
          <Row>
            <Col xs={3}>
              <HyperDatepicker
                hideAddon={false}
                value={new Date(filter.startDate)}
                dateFormat={t('ui.filter.date.format')}
                onChange={function (date: Date): void {
                  onChangeFilter('startDate', date)
                }}
                maxDate={new Date(filter.endDate)}
              />
            </Col>
            <Col xs={3}>
              <HyperDatepicker
                hideAddon={false}
                value={new Date(filter.endDate)}
                dateFormat={t('ui.filter.date.format')}
                onChange={function (date: Date): void {
                  onChangeFilter('endDate', date)
                }}
                minDate={new Date(filter.startDate)}
              />
            </Col>
            <StyledCol>
              <ButtonGroup>
                <Button variant='secondary' className='btn-rounded btn-sm' onClick={() => onStartDate24Hour()}>
                  {t('ui.filter.date.24h')}
                </Button>
                <Button variant='secondary' className='btn-rounded btn-sm' onClick={() => onStartDate7Days()}>
                {t('ui.filter.date.1w')}
                </Button>
                <Button variant='secondary' className='btn-rounded btn-sm' onClick={() => onStartDateMonth(1)}>
                {t('ui.filter.date.1m')}
                </Button>
                <Button variant='secondary' className='btn-rounded btn-sm' onClick={() => onStartDateMonth(6)}>
                {t('ui.filter.date.6m')}
                </Button>
                <PopoverLabel name='btn-filter-date' size='22' marginLeft='10'>
                  {t('ui.filter.date.description')}
                </PopoverLabel>
              </ButtonGroup>
            </StyledCol>
          </Row>
        </Card>

        <Card>
          <h5>{t('ui.filter.perf.title')}</h5>
          <Row>
            <Col xs={2}>
              <label>{t('ui.filter.perf.accuracy')}</label>
            </Col>
            <Col xs={6}>
              <Row>
                <Col xs={1}>
                  <label>{filter.accuracy.min}</label>
                </Col>
                <Col style={{ paddingTop: '5px' }}>
                  <Nouislider
                    range={{ min: 0, max: 100 }}
                    start={[0, 100]}
                    connect
                    onSlide={(render, handle, value, un, percent) => onChangeRange('accuracy', value)}
                  />
                </Col>
                <Col xs={1}>
                  <label>{filter.accuracy.max}</label>
                </Col>
              </Row>
            </Col>
          </Row>
          <br/>
          <Row>
            <Col xs={2}>
              <label>{t('ui.filter.perf.inftime')}</label>
            </Col>
            <Col xs={6}>
              <Row>
                <Col xs={1}>
                  <label>{filter.inference_time.min}</label>
                </Col>
                <Col style={{ paddingTop: '5px' }}>
                  <Nouislider
                    range={{ min: 0, max: 1000 }}
                    start={[0, 1000]}
                    connect
                    onSlide={(render, handle, value, un, percent) => onChangeRange('inference_time', value)}
                  />
                </Col>
                <Col xs={1}>
                  <label>{filter.inference_time.max}</label>
                </Col>
              </Row>
            </Col>
          </Row>
        </Card>
      </Modal.Body>
      <Modal.Footer>
        <Button variant='success' className='btn-rounded btn-sm' onClick={onInitFilter}>
        {t('button.init')}
        </Button>
        <Button variant='info' className='btn-rounded btn-sm' onClick={onCloseModal}>
        {t('button.apply')}
        </Button>
        <Button variant='dark' className='btn-rounded btn-sm' onClick={toggle}>
        {t('button.cancel')}
        </Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ClsFilterModal
