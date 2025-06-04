import { CustomToast } from 'components/Toasts'
import { useEffect, useState } from 'react'
import { Col, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import { ContentArea, ModelCard, ModelCell, ModelColumn, ModelColumnHeader, ModelsArea, TrainModelContainer } from 'features/TrainModel'
import CFMatrix from './CFMatrix'

type Props = {
  models?: any
}

const ModelTestResultTable = ({ models }: Props) => {
  const [t] = useTranslation('translation')

  const [resultColumn, setResultColumn] = useState<any>([])

  useEffect(() => {
    if (!models || models.length < 1) return

    let all_headers = Object.keys(models[0].all_result)
    all_headers = all_headers.filter((header: string) => {
      if (header === 'model') return false
      if (header === 'epoch') return false
      return true
    })

    all_headers.sort((a: string, b: string) => {
      if (a === 'wa') return -1
      if (a === 'uwa' && b === 'wa') return 1
      if (a === 'uwa') return -1
      return 0
    })

    let resultColumn: any[] = []

    // all_result
    all_headers.forEach(function (header: string) {
      resultColumn.push(
        <ModelColumn key={`column-${header}`}>
          <ModelColumnHeader key={`column-${header}-header`}>
            {t(`metric.${header}`)}
          </ModelColumnHeader>
          {
            models.map(function (model: any) {
              return (
                <ModelCell key={`all_column-${header}-row-${model.train_id}_${model.name}`}>
                  {Number(model.all_result[header] * 100 || 0).toFixed(2)}
                </ModelCell>
              )
            })
          }
        </ModelColumn>
      )
    })

    let class_headers = Object.keys(models[0].class_result)
    class_headers = class_headers.filter((header: string) => {
      if (header === 'model') return false
      return true
    })
    class_headers.sort()

    // class_result
    class_headers.forEach(function (header: string) {
      resultColumn.push(
        <ModelColumn key={`column-${header}`}>
          <ModelColumnHeader key={`column-${header}-header`} title={header} header={true}>
            {header.replace('_acc', t('train.result_model.acc')).replace('_count', t('train.result_model.count'))}
          </ModelColumnHeader>
          {
            models.map(function (model: any) {
              return (
                <ModelCell key={`class_column-${header}-row-${model.train_id}_${model.name}`}>
                  {model.class_result[header]}
                </ModelCell>
              )
            })
          }
        </ModelColumn>
      )
    })
    setResultColumn(resultColumn)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [models])

  const [toasts, setToasts] = useState<number[]>([])
  const handleCloseAlert = (index: number) => {
    const updatedList = [...toasts]
    updatedList.splice(index, 1)
    setToasts(updatedList)
  }

  const onAutoClose = (index: number) => {
    handleCloseAlert(index)
  }

  const [cf_matrix, setCFMatrix] = useState<any>()
  if (!models) return <></>

  return (
    <Row>
      {(toasts || []).map((color, index) => {
        return (
          <CustomToast key={`custom-toast-${index}`} index={index} message={t('ui.info.copy.clipboard')} onAutoClose={onAutoClose} />
        )
      })}
      <Col sm={6}>
        <TrainModelContainer>
          <ModelsArea>
            <ModelCard
              key={`empty-left-header`}
              model={null}
              isHeader={true}
            />
            {models && models.map((model: any) => {
              return (
                <ModelCard
                  key={`${model.train_id}_${model.name}`}
                  model={model}
                  onClick={()=>setCFMatrix(model.cf_matrix.String)}
                />
              )
            })}
          </ModelsArea>
          <ContentArea>{resultColumn}</ContentArea>
        </TrainModelContainer>
      </Col>

      <Col sm={6}>
            {cf_matrix && <CFMatrix cf_matrix={cf_matrix} />}
      </Col>
    </Row>
  )
}

export default ModelTestResultTable
