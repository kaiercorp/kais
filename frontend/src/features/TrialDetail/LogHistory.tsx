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
${(props) => props.wsize === 'normal' && 'max-width: 800px;'}
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
${(props) => props.wsize === 'normal' && 'max-width: 800px;'}
`

const LogHistory = ({ trial }: any) => {
    const [logs, setLogs] = useState<any[]>()

    const handleSocketMessage = (e: MessageEvent<any>) => {
        try {
            let msg = JSON.parse(e.data)
            if (!logs || logs.length < 1) setLogs(msg)
            else {
                let newlogs = msg.concat(logs)

                if (newlogs.length > 200) {
                    const oldcount = newlogs.length - 200
                    newlogs.splice(200, oldcount)
                }

                setLogs(newlogs)
            }
            
        } catch (e) {
            logger.error(e)
        }       
    }
    useSocket(`/sys/log/train/${trial?.trial_id}`, 'TrainLog', handleSocketMessage, {shouldCleanup: true, shouldConnect: !!trial && trial.trial_id && trial.trial_id !== 0})

    return (
        <Container>
            <ListArea>
                <Header>
                    <HeaderItem wsize={'mid'}>Date</HeaderItem>
                    <HeaderItem wsize={'small'}>Level</HeaderItem>
                    <HeaderItem wsize={'small'}>Train</HeaderItem>
                    <HeaderItem>Log</HeaderItem>
                </Header>
                <DataList>
                    {
                        logs&&logs.sort(function(x, y) {
                            return x.id - y.id
                          }).map((result: any, index: any) => {
                            return (
                                <DataRow key={`log-datarow-${result.id}`}>
                                    <DataItem wsize={'mid'}>{convertUtcTime(result.created_at)}</DataItem>
                                    <DataItem wsize={'small'}>{result.level}</DataItem>
                                    <DataItem wsize={'small'}>{`Trial${result.train_local_id}`}</DataItem>
                                    <DataItem wsize={'normal'}>{result.log}</DataItem>
                                </DataRow>
                            )
                        })
                    }
                </DataList>
            </ListArea>
        </Container>
    )
}

export default LogHistory