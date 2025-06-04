import { useEffect, useState } from 'react'
import { Button, Form, Col, Modal, Row, Card } from 'react-bootstrap'
import { objDeepCopy } from 'helpers'
import { useTranslation } from 'react-i18next'

type ModalProps = {
  show: boolean
  toggle: () => void
  data: any
  handleData: (value: number[]) => void
}

const ClsClassConfigModal = ({ show, toggle, data, handleData }: ModalProps) => {
  const [t] = useTranslation('translation')
  const [classes, setClasses] = useState(objDeepCopy(data))

  useEffect(() => {
    if (show) {
      setClasses(objDeepCopy(data))
    }
  }, [show, data])

  const handleConfigButton = (key: string, value: number) => {
    let newClasses = classes.slice()
    newClasses.forEach((nc: any) => {
      const label = Object.keys(nc)[0]
      let val = nc[label]

      if (label === key) val = nc[label] + value
      if (val < 1) val = 1
      if (val > 100) val = 100

      nc[label] = val
    })
    setClasses(newClasses)
  }

  const handleConfigData = () => {
    handleData(classes)
    toggle()
  }

  return (
    <Modal show={show} keyboard={true} onHide={toggle}>
      <Modal.Header closeButton>{t('ui.train.classweight')}</Modal.Header>
      <Modal.Body>
        {classes &&
          classes.map((elem: string) => {
            const key = (Object.keys(elem) as (keyof typeof elem)[])[0]
            const keystr = key.toString()
            return (
              <Card key={`class-label-${keystr}`}>
                <Row>{key.toString()}</Row>
                <Row>
                  <Col xs={4}>
                    <Form.Control type='number' value={Number(elem[key])} readOnly />
                  </Col>
                  <Col xs={8}>
                    <Button
                      onClick={() => {
                        handleConfigButton(keystr, -10)
                      }}
                    >
                      -10
                    </Button>
                    <Button
                      onClick={() => {
                        handleConfigButton(keystr, -1)
                      }}
                    >
                      -1
                    </Button>
                    <Button
                      onClick={() => {
                        handleConfigButton(keystr, 1)
                      }}
                    >
                      1
                    </Button>
                    <Button
                      onClick={() => {
                        handleConfigButton(keystr, 10)
                      }}
                    >
                      10
                    </Button>
                  </Col>
                </Row>
              </Card>
            )
          })}
      </Modal.Body>
      <Modal.Footer>
        <Button variant='success' className='btn-rounded btn-sm' onClick={() => setClasses(objDeepCopy(data))}>
          {t('button.init')}
        </Button>
        <Button variant='info' className='btn-rounded btn-sm' onClick={handleConfigData}>
          {t('button.apply')}
        </Button>
        <Button variant='dark' className='btn-rounded btn-sm' onClick={toggle}>
          {t('button.cancel')}
        </Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ClsClassConfigModal
