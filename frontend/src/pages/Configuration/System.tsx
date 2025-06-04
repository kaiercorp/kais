import { useEffect, useState, useContext } from 'react'
import { useTranslation } from 'react-i18next'

import { logger, objDeepCopy, QUERY_KEY } from 'helpers'
import { Button, Card, Col, Form, Row } from 'react-bootstrap'
import { LocationContext } from 'contexts'
import { useQueryClient } from '@tanstack/react-query'

import { ApiFetchConfigs, ApiUpdateConfigs } from 'helpers/api'
import { UserType } from 'common'
import { ButtonArea, CardHeaderLeft, ConfigLabel, ConfigRow } from '../../components/Containers'


const ConfigurationLanding = () => {
    const [t] = useTranslation('translation')
    const { updateLocationContextValue } = useContext(LocationContext)

    const { configs } = ApiFetchConfigs()
    const updateConfigs = ApiUpdateConfigs()

    const queryClient = useQueryClient()
    const user = queryClient.getQueryData<UserType>([QUERY_KEY.login])

    useEffect(() => {
        logger.log(`Change Location to ${t('title.system')}`)
        updateLocationContextValue({ location: 'system' })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const [inits, setInits] = useState<any[]>([])

    useEffect(() => {
        initializeConfigs()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [configs, user])

    const initializeConfigs = () => {
        if (!configs || configs.length < 1) return

        let initval = objDeepCopy(configs)
        initval = initval.filter((config: any) => {
            if (config.config_type === 'SYSTEM') return true
            return false
        })

        setInits(initval)
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

    const onSaveConfigs = () => {
        let newConfigs = objDeepCopy(inits)

        updateConfigs.mutate(newConfigs)
    }

    return (
        <>
            <Card>
                <Card.Header>
                    <CardHeaderLeft>System Configs</CardHeaderLeft>
                </Card.Header>
                <Card.Body>
                    <Row key={`config-init-000`} style={{ border: '1px solid #999999', marginBottom: '3px' }}>
                        <Col sm={2}>
                            <ConfigLabel>Type</ConfigLabel>
                        </Col>
                        <Col sm={4}>
                            <ConfigLabel>Key</ConfigLabel>
                        </Col>
                        <Col sm={4}>
                            <ConfigLabel>Value</ConfigLabel>
                        </Col>
                    </Row>
                    {
                        inits?.map((init: any) => {
                            return (
                                <ConfigRow key={`config-init-${init.id}`}>
                                    <Col sm={2}>
                                        <ConfigLabel>{init.config_type}</ConfigLabel>
                                    </Col>
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

            <ButtonArea>
                <Button variant={'danger'} onClick={initializeConfigs}>{t('button.cancel')}</Button>
                <Button onClick={onSaveConfigs}>{t('button.save')}</Button>
            </ButtonArea>
        </>
    )
}

export default ConfigurationLanding