import { useEffect, useState } from 'react'
import styled from 'styled-components'

const Container = styled.div`
display: flex;
`

const ListArea = styled.div`
min-width: 100px;
max-width: 1000px;
overflow-x: auto;
margin: 0 10px;
`

const Header = styled.div`
display: flex;
border-bottom: 1px solid grey;
`

interface IHeaderItem {
    wsize?: string
}
const HeaderItem = styled.div<IHeaderItem>`
font-size: 12px;
padding: 5px;
${(props) => props.wsize === 'small' && 'width: 50px;'}
${(props) => props.wsize === 'mid' && 'width: 100px;'}
${(props) => props.wsize === 'normal' && 'width: 150px;'}
${(props) => props.wsize === 'big' && 'width: 300px;'}
`

const DataList = styled.div`
max-height: 400px;
`

interface IDataRow {
    isIncorrect?: boolean
}
const DataRow = styled.div<IDataRow>`
display: flex;
cursor: pointer;
${(props) => props.isIncorrect && 'border: 1px solid red;'}
`

interface IDataItem {
    wsize?: string
}
const DataItem = styled.div<IDataItem>`
font-size: 12px;
padding: 5px;
${(props) => props.wsize === 'small' && 'width: 50px;'}
${(props) => props.wsize === 'mid' && 'width: 100px;'}
${(props) => props.wsize === 'normal' && 'width: 150px;'}
${(props) => props.wsize === 'big' && 'width: 300px;'}
overflow-x: hidden;
`

const cfg = ['tp', 'tn', 'fp', 'fn']
const perf = ['accuracy', 'precision', 'recall', 'f1_score']

const TestDataResult = ({ resultList, onClick }: any) => {
    const [header, setHeader] = useState<any>()
    useEffect(() => {
        if (!resultList || !resultList[0] || !resultList[0].data) return

        resultList.forEach((result:any) => {
            result.props = JSON.parse(result.data)
        })
        let _header = Object.keys(resultList[0].props)
        _header = _header.filter((h:string) => {
            if (['kaier_id', 'sr.true', 'sr.prediction', 'sr.true_label', 'sr.predicted_label'].includes(h)) return false
            return true
        })
        
        if (resultList[0].save_path?.includes('vision-cls-ml')) {
            _header = [..._header.filter((header: string) => !(cfg.includes(header) || perf.includes(header))).sort()]           
        }

        setHeader(_header)
    }, [resultList])

    return (
        <Container>
            <ListArea>
                <Header>
                    <HeaderItem wsize={'small'}>No.</HeaderItem>
                    <HeaderItem wsize={'big'}>Data</HeaderItem>
                    <HeaderItem wsize={'mid'}>Ground Truth</HeaderItem>
                    <HeaderItem wsize={'mid'}>Prediction</HeaderItem>
                    <HeaderItem wsize={'small'}>Result</HeaderItem>
                    {
                        header && header.map((h: any) => {
                            return <HeaderItem key={`headeritem-${h}`} wsize={'mid'}>{h}</HeaderItem>
                        })
                    }
                </Header>
                <DataList>
                    {
                        resultList&&resultList.map((result: any, index: any) => {
                            let isIncorrect = false
                            if (result.props['sr.true_label'] && result.props['sr.predicted_label'] && (result.props['sr.true_label'] !== result.props['sr.predicted_label'])) isIncorrect = true
                            return (
                                <DataRow key={`datarow-${index}-${result.props.kaier_id}`} onClick={() => onClick(result.id)} isIncorrect={isIncorrect}>
                                    <DataItem wsize={'small'}>{index + 1}</DataItem>
                                    <DataItem wsize={'big'}>{result.props.kaier_id}</DataItem>
                                    <DataItem wsize={'mid'}>{result.props['sr.true_label']}</DataItem>
                                    <DataItem wsize={'mid'}>{result.props['sr.predicted_label']}</DataItem>
                                    <DataItem wsize={'small'}>{isIncorrect ? 'X' : 'O'}</DataItem>
                                    {
                                        header&&header.map((h: any) => {
                                            return (<DataItem key={`dataitem-${result.props.kaier_id}-${h}-${index}`} wsize={'mid'}>{(result.props[h]*100||0).toFixed(2)+'%'}</DataItem>)
                                        })
                                    }
                                </DataRow>
                            )
                        })
                    }
                </DataList>
            </ListArea>
        </Container>
    )
}

export default TestDataResult