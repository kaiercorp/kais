import { useEffect, useState, useContext, SetStateAction } from 'react'
import { useTranslation } from 'react-i18next'

import { logger, objDeepCopy } from 'helpers'
import { Button, Card, Col, Form, Row } from 'react-bootstrap'
import { DefaultSwitch } from 'components'
import { LocationContext } from 'contexts'

import korFlag from 'assets/images/flags/korea.png'
import chinaFlag from 'assets/images/flags/china.png'
import engFlag from 'assets/images/flags/us.jpg'
import japanFlag from 'assets/images/flags/japan.jpg'
import { ApiFetchConfigs, ApiUpdateConfigs } from 'helpers/api'
import { ConfigType } from 'common'
import { ButtonArea, CardHeaderLeft, ConfigLabel, ConfigRow } from '../../components/Containers'
import { APICore } from 'helpers/api/apiCore'

const langs = [
    {
        name: '한국어',
        key: 'ko',
        flag: korFlag
    },
    {
        name: '中文',
        key: 'ch',
        flag: chinaFlag
    },
    {
        name: 'English',
        key: 'en',
        flag: engFlag
    },
    {
        name: '日本語',
        key: 'jp',
        flag: japanFlag
    },
]

const api = new APICore()

const Config = () => {
    const [t] = useTranslation('translation')
    const { updateLocationContextValue } = useContext(LocationContext)

    const { configs } = ApiFetchConfigs()
    const updateConfigs = ApiUpdateConfigs()

    const [language, setLanguage] = useState('ko')
    const [inits, setInits] = useState<any[]>([])
    const user = api.getLoggedInUser()

    useEffect(() => {
        logger.log(`Change Location to ${t('title.setting')}`)
        updateLocationContextValue({ location: 'setting' })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        initializeConfigs()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [configs])

    const initializeConfigs = () => {
        if (!configs || configs.length < 1) return

        initLang(configs, setLanguage)

        let initval = objDeepCopy(configs)
        initval = initval.filter((config: any) => {
            if (config.config_type === 'AUTO_STOP') return true
            if (config.config_type === 'FAIL_STOP') return true
            if (config.config_type === 'REMOVE_TOPN') return true
            return false
        })

        setInits(initval)
    }

    const initLang = (configs: ConfigType[], setConfig: (value: SetStateAction<any>, user?: any) => void) => {
        let config = configs.filter((config: any) => {
            return config.config_key === 'LANGUAGE'
        })[0]

        setConfig(config.config_val)
    }

    const onChangeInitConfigItem = (type: string, id: number, value: string) => {
        let newInits = inits.slice()
        newInits.forEach((init: any) => {
            if (init.id === id) {
                if (type === 'key') {
                    init.config_key = value
                } else if (type === 'val') {
                    init.config_val = value
                }
                return false
            }
        })
        setInits(newInits)
    }

    const onSave = () => {
        let newConfigs = objDeepCopy(inits)

        configs.forEach((config: any) => {
            if (config.config_key === 'LANGUAGE') {
                let newLang = objDeepCopy(config)
                newLang.config_val = language
                newConfigs.push(newLang)
            }
        })

        updateConfigs.mutate(newConfigs)
    }

    const onCancel = () => {
        initializeConfigs()
    }

    return (<>
        <Card>
            <Card.Header><CardHeaderLeft>Language selection</CardHeaderLeft></Card.Header>
            <Card.Body>
                <Row>
                    {langs.map((lang: any) => {
                        return (
                            <Col sm={2} key={`lang-${lang.key}`}>
                                <DefaultSwitch
                                    label={lang.name}
                                    checked={lang.key === language}
                                    image={lang.flag}
                                    onChange={() => setLanguage(lang.key)}
                                />
                            </Col>
                        )
                    })}
                </Row>
            </Card.Body>
        </Card>
        {
            user &&
            <Card>
                <Card.Header><CardHeaderLeft>Auto Stop Count</CardHeaderLeft></Card.Header>
                <Card.Body>
                    <Row key={`config-init-000`} style={{ border: '1px solid #999999', marginBottom: '3px' }}>
                        <Col sm={4}>
                            <ConfigLabel>Engine</ConfigLabel>
                        </Col>
                        <Col sm={4}>
                            <ConfigLabel>Value</ConfigLabel>
                        </Col>
                    </Row>
                    {
                        inits?.filter((config: any) => config.config_type === 'AUTO_STOP').map((init: any) => {
                            return (
                                <ConfigRow key={`config-init-${init.id}`}>
                                    <Col sm={4}>
                                        <ConfigLabel>{init.config_key}</ConfigLabel>
                                    </Col>
                                    <Col sm={4}>
                                        <Form.Control
                                            type='text'
                                            size='sm'
                                            value={init.config_val}
                                            onChange={(e: any) => onChangeInitConfigItem('val', init.id, e.target.value)}
                                        />
                                    </Col>
                                </ConfigRow>

                            )
                        })
                    }
                </Card.Body>
            </Card>
        }
        {
            user && user.token &&
            <Card>
                <Card.Header><CardHeaderLeft>Failed Stop Count</CardHeaderLeft></Card.Header>
                <Card.Body>
                    <Row key={`config-init-000`} style={{ border: '1px solid #999999', marginBottom: '3px' }}>
                        <Col sm={4}>
                            <ConfigLabel>Engine</ConfigLabel>
                        </Col>
                        <Col sm={4}>
                            <ConfigLabel>Value</ConfigLabel>
                        </Col>
                    </Row>
                    {
                        inits?.filter((config: any) => config.config_type === 'FAIL_STOP').map((init: any) => {
                            return (
                                <ConfigRow key={`config-init-${init.id}`}>
                                    <Col sm={4}>
                                        <ConfigLabel>{init.config_key}</ConfigLabel>
                                    </Col>
                                    <Col sm={4}>
                                        <Form.Control
                                            type='text'
                                            size='sm'
                                            value={init.config_val}
                                            onChange={(e: any) => onChangeInitConfigItem('val', init.id, e.target.value)}
                                        />
                                    </Col>
                                </ConfigRow>

                            )
                        })
                    }
                </Card.Body>
            </Card>
        }
        {
            user && user.token &&
            <Card>
                <Card.Header><CardHeaderLeft>TopN model count</CardHeaderLeft></Card.Header>
                <Card.Body>
                    <Row key={`config-init-000`} style={{ border: '1px solid #999999', marginBottom: '3px' }}>
                        <Col sm={4}>
                            <ConfigLabel>Engine</ConfigLabel>
                        </Col>
                        <Col sm={4}>
                            <ConfigLabel>Value</ConfigLabel>
                        </Col>
                    </Row>
                    {
                        inits?.filter((config: any) => config.config_type === 'REMOVE_TOPN').map((init: any) => {
                            return (
                                <ConfigRow key={`config-init-${init.id}`}>
                                    <Col sm={4}>
                                        <ConfigLabel>{init.config_key}</ConfigLabel>
                                    </Col>
                                    <Col sm={4}>
                                        <Form.Control
                                            type='text'
                                            size='sm'
                                            value={init.config_val}
                                            onChange={(e: any) => onChangeInitConfigItem('val', init.id, e.target.value)}
                                        />
                                    </Col>
                                </ConfigRow>

                            )
                        })
                    }
                </Card.Body>
            </Card>
        }
        <ButtonArea>
            <Button variant={'danger'} onClick={onCancel}>{t('button.cancel')}</Button>
            <Button onClick={onSave}>{t('button.save')}</Button>
        </ButtonArea>
    </>)
}

export default Config