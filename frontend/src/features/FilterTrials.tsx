import { forwardRef, useImperativeHandle, useState, useContext } from 'react'
import { Button, ButtonGroup, Card, Col, Form, Row } from 'react-bootstrap'
import Nouislider from 'nouislider-react'
import 'nouislider/distribute/nouislider.css'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

import { initialFilter, FilterType } from 'common'
import { HyperDatepicker, PopoverLabel, StatusLabel } from 'components'
import { stateColors } from 'appConstants'
import { objDeepCopy } from 'helpers'
import { FilterContext } from 'contexts'


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

const FilterTrials = forwardRef(({ toggle }: any, ref) => {
    const [t] = useTranslation('translation')
    const { filterContextValue, updateFilterContextValue } = useContext(FilterContext)
    const [filter, setFilter] = useState<any>(filterContextValue.filter)

    const onChangeFilter = (key: string, value: any) => {
        setFilter((filter: FilterType) => ({
                ...filter,
                [key] : value 
            }))
    }

    const onChangeRange = (key: string, value: any) => {
        setFilter((filter: FilterType) => ({
                ...filter,
                [key] : {
                    min: Number(value[0].toFixed(0)),
                    max: Number(value[1].toFixed(0))
                }
            }))
    }

    const onChangeStatus = (key: keyof FilterType['state']) => {
        let newFilter = objDeepCopy(filter)
        newFilter.state[key] = !filter.state[key]
        newFilter.state.total =
            Object.entries(newFilter.state).filter((state) => {
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
        setFilter((filter: FilterType) => ({
            ...filter,
            startDate: new Date(filter.endDate.getTime() - 24 * 60 * 60 * 1000)            
        }))
    }

    const onStartDate7Days = () => {
        setFilter((filter: FilterType) => ({
            ...filter,
            startDate: new Date(filter.endDate.getTime() - 7 * 24 * 60 * 60 * 1000)
        }))
    }

    const onStartDateMonth = (value: number) => {
        setFilter((filter: FilterType) => { 
            const startDate = new Date(filter.endDate.getTime())
            startDate.setDate(1)
            startDate.setMonth(startDate.getMonth() - value)
            startDate.setDate(
                Math.min(filter.endDate.getDate(), new Date(startDate.getFullYear(), startDate.getMonth() + 1, 0).getDate())
            )           

            return {
                ...filter,
                startDate 
            }
        })
    }

    useImperativeHandle(ref, () => ({
        handleSubmit
    }))

    const handleSubmit = () => {
        updateFilterContextValue({filter})
        return true
    }

    const onInitFilter = () => {
        setFilter(initialFilter)
    }

    return (
        <Form>

            <Card>
                <Row>
                    <Col>
                        <Button variant='success' className='btn-rounded btn-sm' onClick={onInitFilter}>
                            {t('button.init')}
                        </Button>
                    </Col>
                </Row>
            </Card>
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
                                    start={filter ? [filter.accuracy.min, filter.accuracy.max] : [0, 100]}
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
                <br />
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
                                    start={filter ? [filter.inference_time.min, filter.inference_time.max] : [0, 1000]}
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
        </Form>
    )
})

export default FilterTrials