import classNames from 'classnames'
import { HPOModelType, HPOParamsDistType, HPOParamsType } from 'common'
import { ButtonArea, CardHeaderLeft, CardHeaderRight } from 'components/Containers'
import { Button, Card, Col, Form, Row } from 'react-bootstrap'
import styled from 'styled-components'
import { useId } from 'react'

type HPModelType = {
    hps: HPOModelType[]
    engine: string
    onChange: any
}

type HPModelCardType = {
    engine: string
    model_idx: number
    model: HPOModelType
    onChange: any
}

type HPParamCardType = {
    engine: string
    model_idx: number
    param_idx: number
    param: HPOParamsType
    disabled: boolean
    onChange: any
}

type HPDistCardType = {
    engine: string
    model_idx: number
    param_idx: number
    dist_idx: number
    dist: HPOParamsDistType
    disabled: boolean
    onChange: any
}

const StyledContaier = styled(Row)`
    width: 100%;
`

const StyledModelContainer = styled(Card)`
    border: 1px solid #999999;
    margin-top: 5px;
    & > .card-header {
        background-color: #242c39;
    }
`

const StyledParamContainer = styled(Card)`
    border: 1px dotted #111111;
    padding: 0;
`

const StyledDistContainer = styled(Card)`
    padding: 0;
    & > .card-header {
        background-color: #353d4a;
    }
    
`

const ModelCard = ({ model, onChange, engine, model_idx }: HPModelCardType) => {
    const id = useId()
    return (
        <Col sm={6}>
            <StyledModelContainer>
                <Card.Header>
                    <CardHeaderLeft>Model : {model.model_name}</CardHeaderLeft>
                    <CardHeaderRight>
                        <Button className="btn-icon" onClick={(e) => onChange(e, engine, "add", "param", model_idx)}>
                            <i className={classNames('mdi', 'mdi-plus')} /> <span>Param</span>
                        </Button>
                        <Button variant={'danger'} className="btn-icon" onClick={(e) => onChange(e, engine, "delete", "model", model_idx)}>
                            <i className={classNames('mdi', 'mdi-trash-can')} />
                        </Button>
                    </CardHeaderRight>
                </Card.Header>
                <Card.Body>
                    <Row>
                        <Col>
                            <Form.Group className={'mb-3'}>
                                <Form.Control
                                    type={'text'}
                                    name={'model_name'}
                                    value={model.model_name}
                                    onChange={(e) => onChange(e, engine, "edit", "model", model_idx, -1, -1, e.target.name, e.target.value)}
                                />
                            </Form.Group>
                        </Col>
                        <Col>
                            <Form.Group className={'mb-3'}>
                                <Form.Select
                                    name="class_type"
                                    value={model.class_type}
                                    onChange={(e) => onChange(e, engine, "edit", "model", model_idx, -1, -1, e.target.name, e.target.value)}
                                >
                                    <option>bin</option>
                                    <option>both</option>
                                    <option>multi</option>
                                </Form.Select>
                            </Form.Group>
                        </Col>
                        <Col>
                            <Form.Group className={'mb-3'}>
                                <Form.Check
                                    type={'checkbox'}
                                    label={"Is Use"}
                                    name={"is_use"}
                                    checked={model.is_use}
                                    onChange={(e) => onChange(e, engine, "edit", "model", model_idx, -1, -1, e.target.name, e.target.checked)}
                                />
                            </Form.Group>
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            {model?.params?.map((param: HPOParamsType, index: number) => {
                                return <ParamCard
                                    key={`param-${id}-${index}`}
                                    engine={engine}
                                    model_idx={model_idx}
                                    param_idx={index}
                                    disabled={!model.is_use}
                                    param={param}
                                    onChange={onChange}
                                />
                            })}
                        </Col>
                    </Row>
                </Card.Body>
            </StyledModelContainer>
        </Col>
    )
}

const ParamCard = ({ param, disabled, onChange, engine, model_idx, param_idx }: HPParamCardType) => {
    const id = useId()
    return (
        <StyledParamContainer>
            <Card.Header>
                <CardHeaderLeft>Parameter : {param.name}</CardHeaderLeft>
                <CardHeaderRight>
                    <Button className="btn-icon" onClick={(e) => onChange(e, engine, "add", "dist", model_idx, param_idx)}>
                        <i className={classNames('mdi', 'mdi-plus')} /> <span>Dist</span>
                    </Button>
                    <Button variant={'danger'} className="btn-icon" onClick={(e) => onChange(e, engine, "delete", "param", model_idx, param_idx)}>
                        <i className={classNames('mdi', 'mdi-trash-can')} />
                    </Button>
                </CardHeaderRight>
            </Card.Header>
            <Card.Body>
                <Row>
                    <Col>
                        <Form.Group className={'mb-3'}>
                            <Form.Control
                                type={'text'}
                                name={'name'}
                                value={param.name}
                                disabled={disabled}
                                onChange={(e) => onChange(e, engine, "edit", "param", model_idx, param_idx, -1, e.target.name, e.target.value)}
                            />
                        </Form.Group>
                    </Col>
                    <Col>
                        <Form.Group className={'mb-3'}>
                            <Form.Select
                                name="suggest_type"
                                value={param.suggest_type}
                                disabled={disabled}
                                onChange={(e) => onChange(e, engine, "edit", "param", model_idx, param_idx, -1, e.target.name, e.target.value)}
                            >
                                <option>category</option>
                                <option>continuous</option>
                            </Form.Select>
                        </Form.Group>
                    </Col>
                    <Col>
                        <Form.Group className={'mb-3'}>
                            <Form.Select
                                name="data_type"
                                value={param.data_type}
                                disabled={disabled}
                                onChange={(e) => onChange(e, engine, "edit", "param", model_idx, param_idx, -1, e.target.name, e.target.value)}
                            >
                                <option>bool</option>
                                <option>float</option>
                                <option>int</option>
                                <option>str</option>
                            </Form.Select>
                        </Form.Group>
                    </Col>
                    <Col>
                        <Form.Group className={'mb-3'}>
                            <Form.Check
                                type={'checkbox'}
                                label={"Is Use"}
                                name={"is_use"}
                                checked={param.is_use}
                                disabled={disabled}
                                onChange={(e) => onChange(e, engine, "edit", "param", model_idx, param_idx, -1, e.target.name, e.target.checked)}
                            />
                        </Form.Group>
                    </Col>
                </Row>
                <Row>
                    <Col>
                        {param?.dists?.map((dist: HPOParamsDistType, index: number) => {
                            return <DistCard
                                key={`dist-${id}-${index}`}
                                disabled={!param.is_use || disabled}
                                dist={dist}
                                onChange={onChange}
                                engine={engine}
                                model_idx={model_idx}
                                param_idx={param_idx}
                                dist_idx={index}
                            />
                        })}
                    </Col>
                </Row>
            </Card.Body>
        </StyledParamContainer>
    )
}

const DistCard = ({ dist, disabled, onChange, engine, model_idx, param_idx, dist_idx }: HPDistCardType) => {
    return (
        <StyledDistContainer>
            <Card.Header>
                <CardHeaderLeft>Dist {dist.id}</CardHeaderLeft>
                <CardHeaderRight>
                    <Button variant={'danger'} className="btn-icon" onClick={(e) => onChange(e, engine, "delete", "dist", model_idx, param_idx, dist_idx)}>
                        <i className={classNames('mdi', 'mdi-trash-can')} />
                    </Button>
                </CardHeaderRight>
            </Card.Header>
            <Card.Body>
                <Row>
                    <Col>
                        <Form.Group className={'mb-3'}>
                            <Form.Control
                                type={'text'}
                                name={'dist'}
                                value={dist.dist}
                                disabled={disabled}
                                onChange={(e) => onChange(e, engine, "edit", "dist", model_idx, param_idx, dist_idx, e.target.name, e.target.value)}
                            />
                        </Form.Group>
                    </Col>
                    <Col>
                        <Form.Group className={'mb-3'}>
                            <Form.Check
                                type={'checkbox'}
                                label={"Is Use"}
                                name={"is_use"}
                                checked={dist.is_use}
                                disabled={disabled}
                                onChange={(e) => onChange(e, engine, "edit", "dist", model_idx, param_idx, dist_idx, e.target.name, e.target.checked)}
                            />
                        </Form.Group>
                    </Col>
                </Row>
                <Row>
                    <Col>
                        <Form.Control
                            type={'text'}
                            name={'cond.key'}
                            value={dist.cond.key}
                            disabled={disabled}
                            onChange={(e) => onChange(e, engine, "edit", "dist", model_idx, param_idx, dist_idx, e.target.name, e.target.value)}
                        />
                    </Col>
                    <Col>
                        <Form.Control
                            type={'text'}
                            name={'cond.operator'}
                            value={dist.cond.operator}
                            disabled={disabled}
                            onChange={(e) => onChange(e, engine, "edit", "dist", model_idx, param_idx, dist_idx, e.target.name, e.target.value)}
                        />
                    </Col>
                    <Col>
                        <Form.Control
                            type={'text'}
                            name={'cond.value'}
                            value={dist.cond.value}
                            disabled={disabled}
                            onChange={(e) => onChange(e, engine, "edit", "dist", model_idx, param_idx, dist_idx, e.target.name, e.target.value)}
                        />
                    </Col>
                </Row>
            </Card.Body>
        </StyledDistContainer>
    )
}

const HPCardList = ({ hps, engine, onChange }: HPModelType) => {
    const id = useId()
    return (
        <StyledContaier>
            <ButtonArea>
                <Button onClick={(e) => onChange(e, engine, "add", "model")}>Add New Model</Button>
            </ButtonArea>
            {hps.map((hp: HPOModelType, index: number) => {
                return (
                    <ModelCard key={`models-${id}-${index}`} engine={engine} model_idx={index} model={hp} onChange={onChange} />
                )
            })}

        </StyledContaier>
    )
}

export default HPCardList