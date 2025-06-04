import { useState, useEffect } from 'react'
import { Col, FormGroup, Form, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'

import config from 'config'
import { engine } from 'appConstants/trial'

interface Props {
  max?: boolean
  position?: string
}

const ResultTable = styled.table`
  background-color: #dddddd;
  border: 1px solid grey;
`

const ResultHeader = styled.thead`
  background-color: #999999;
`

const ResultTH = styled.th`
  min-width: 100px;
  height: 30px;
  text-align: center;
  color: black;
`

const StyledTR = styled.tr<Props>`
  background-color: ${props=>props.max===true?'white':'inherit'}
`

const StyledTD = styled.td<Props>`
  font-weight: ${props=> props.position==='bottom'?600:400};
  color: black;
  text-align: center;
  ${props=>props.position==='bottom'?'border-top: 1px solid grey;':''}
`

const StyledNumTD = styled.td<Props>`
  font-weight: ${props=> props.position==='bottom'?600:400};
  color: black;
  text-align: right;
  padding: 0 5%;
  ${props=>props.position==='bottom'?'border-top: 1px solid grey;':''}
`
const ResultDiv = styled.div`
  flex: 1;
  background-color: white;
  border-radius: 5px;
  padding: 5px;
  align-items: center;
  display: flex;
  justify-content: center;

  & img {
    width: 100%
  }
`

const ResultTR = (props: any) => {
  return (
    <StyledTR max={props.max}>
      <StyledTD>{props.index}</StyledTD>
      <StyledTD>{props.label}</StyledTD>
      <StyledTD>{props.result && <Form.Check ><Form.Check.Input type='checkbox' checked={props.result} isValid={props.result}></Form.Check.Input></Form.Check>}</StyledTD>
      <StyledNumTD>{(props.value * 100).toFixed(2)}%</StyledNumTD>
    </StyledTR>
  )
}

const StyledCol = styled(Col)`
  display: flex;
  flex-direction: column;
  color: #ffffff;
  font-size: 14px;
  margin-bottom: 10px;
`

function FileTestResult({data, engineType}: any) {
  const [t] = useTranslation('translation')
  const [overlay, setOverlay] = useState(true)
  const onChange = () => setOverlay(!overlay)

  const imgOrigin = <img alt='' src={config.staticURL + data.origin_path} />
  const imgHeatmap = <img alt='' src={config.staticURL + data.heatmap_path} />
  const imgMerged = <img alt='' src={config.staticURL + data.overlay_path} />

  const [heatmap, setHeatmap] = useState(<img alt='' src={config.staticURL + data.heatmap} />)
  useEffect(() => {
    if (overlay) setHeatmap(imgMerged)
    else setHeatmap(imgHeatmap)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [overlay])
  
  const predictedLabelIndex = data?.proba?.map((p: number, index: number) => {
    return p >= 0.5 ? index : -1
  }).filter((p: number) => {return p !== -1})

  return (
    <>
      <FormGroup>
          <Row>
            <StyledCol xs='6'>
              <Row className="title-checkbox">
                <label>&nbsp;&nbsp;&nbsp;&nbsp;{t('ui.train.title.origin')}</label>
              </Row>
              <ResultDiv>
                {imgOrigin} 
              </ResultDiv>
            </StyledCol>

          {engineType === engine.vision_cls_sl && 
            <StyledCol xs='6'>
              <Row className="title-checkbox">
                <label>&nbsp;&nbsp;&nbsp;&nbsp;{t('ui.train.image.heatmap')}&nbsp;&nbsp;</label>
                <div >
                  <input type="checkbox" checked={overlay} onChange={onChange} />
                  <label>&nbsp;{t('ui.train.image.overlay')}</label>
                </div>
              </Row>
              <ResultDiv>
                {heatmap}
              </ResultDiv>
            </StyledCol>
          }
          </Row>
      </FormGroup>

      <FormGroup>
          <Row style={{marginBottom: '10px', color: 'white'}}>
            <Form.Label column='sm' sm={4}>{t('ui.result.sr.predicted_label')}</Form.Label>
            <Form.Label column='sm' sm={8}>
              {predictedLabelIndex.map((labelIndex: number, index: number) => {
                const label = data?.label[labelIndex]               
                
                return index === predictedLabelIndex.length - 1 ? label : label + ', '
              })}
            </Form.Label>
          </Row>
          <ResultTable>
            <ResultHeader>
              <tr>
                <ResultTH></ResultTH>
                <ResultTH>{t('ui.formatter.classname')}</ResultTH>
                <ResultTH>{t('ui.result.sr.prediction')}</ResultTH>
                <ResultTH>{t('ui.formatter.ratio')}</ResultTH>
              </tr>
            </ResultHeader>
            <tbody>
            {
              data?.label?.map((key: string, index: number) => {
                return <ResultTR key={`resultTr-${key}`} index={index + 1} label={key} value={data.proba[index]} result={predictedLabelIndex.includes(index)} max={true} />
              })
            }
            {engineType === engine.vision_cls_sl && 
              <tr>
                <StyledTD position='bottom'>{t('ui.label.total')}</StyledTD>
                <StyledTD position='bottom'></StyledTD>
                <StyledNumTD position='bottom'>{((data.proba[0] + data.proba[1]) * 100).toFixed(2)}%</StyledNumTD>
              </tr>
            }
            </tbody>
          </ResultTable>
      </FormGroup>
    </>
  )
}

export default FileTestResult
