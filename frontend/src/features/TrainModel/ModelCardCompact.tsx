import { convertUtcTime } from "helpers"
import styled from "styled-components"

interface ICardContent {
    header?: boolean
}

const CardContent = styled.div<ICardContent>`
font-size: 12px;
color: black;
cursor: pointer;
border-bottom: 1px solid grey;
padding-top: 3px;
overflow: hidden;
text-overflow: ellipsis;
white-space: nowrap;
word-break: break-all;
`

const ModelCardCompact = ({ model, isHeader=false, onClick, children }: any) => {
    return (
        <CardContent header={isHeader} onClick={onClick}>
            <div className="d-flex justify-content-between mb-0">
                {model && <div style={{marginLeft:'10px', marginRight:'10px'}}>
                    <span className="text-muted mb-0">{`Trial${model.train_id} `}</span><b>{model.name}</b>
                    <p className="text-muted mb-0">
                        {`${convertUtcTime(model.updated_at)}`}
                    </p>
                </div>}
                {children}
            </div>
        </CardContent>
    )
}

export default ModelCardCompact