import { convertUtcTime, logger } from 'helpers'
import { useState } from 'react'
import styled from 'styled-components'
import { useSocket } from 'hooks'


const Container = styled.div`
display: flex;
`

const ListArea = styled.div`
min-width: 100px;
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
${(props) => props.wsize === 'small' && 'width: 100px;'}
${(props) => props.wsize === 'mid' && 'width: 150px;'}
${(props) => props.wsize === 'normal' && 'width: 250px;'}
`

const DataList = styled.div`
overflow-x: auto;
`

interface IDataRow {
    isIncorrect?: boolean
}
const DataRow = styled.div<IDataRow>`
display: flex;
width: max-content;
border-bottom: 1px dotted #999999;
`

interface IDataItem {
    wsize?: string
}
const DataItem = styled.div<IDataItem>`
font-size: 12px;
padding: 5px;
${(props) => props.wsize === 'small' && 'width: 100px;'}
${(props) => props.wsize === 'mid' && 'width: 150px;'}
${(props) => props.wsize === 'normal' && 'width: 250px;'}
`

const getParamItems = (params: string, index: number) => {
    if (!params || params === '') return null

    const p = JSON.parse(params)
    const header = Object.keys(p)
    return header.map((h: any, idx: number) => {
        return <DataItem key={`dataitem-${index}-${idx}`}>{h}: <b>{p[h]}</b></DataItem>
    })
}

const HyperParameterHistory = ({ trial }: any) => {
    const [hphistory, setHphistory] = useState<any[]>()

    const handleSocketMessage = (e: MessageEvent<any>) => {
        try {
            let msg = JSON.parse(e.data)
            setHphistory(msg)
        } catch (e) {
            logger.error(e)
        }       
    }
    useSocket(`/trials/hphistory/${trial?.trial_id}`, 'Trains', handleSocketMessage, {shouldCleanup: true, shouldConnect: !!trial && trial.trial_id && trial.trial_id !== 0})

    return (
        <Container>
            <ListArea>
                <Header>
                    <HeaderItem wsize={'mid'}>Date</HeaderItem>
                    <HeaderItem wsize={'small'}>Train No.</HeaderItem>
                    <HeaderItem wsize={'normal'}>Model</HeaderItem>
                    <HeaderItem>Params</HeaderItem>
                </Header>
                <DataList>
                    {
                        hphistory&&hphistory.sort(function(x, y) {
                            return x.train_local_id - y.train_local_id
                          }).map((result: any, index: any) => {
                            return (
                                <DataRow key={`datarow-${result.train_uuid}`}>
                                    <DataItem wsize={'mid'}>{convertUtcTime(result.created_at)}</DataItem>
                                    <DataItem wsize={'small'}>{`Trial${result.train_local_id}`}</DataItem>
                                    <DataItem wsize={'normal'}>({result.model_num}) {result.model}</DataItem>
                                    {
                                        result.params && getParamItems(result.params, index)
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

export default HyperParameterHistory