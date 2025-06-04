import { useEffect, useState } from "react"
import styled from "styled-components"

interface ICardContent {
    isSelected?:boolean
}
const CardContent = styled.div<ICardContent>`
font-size: 12px;
color: black;
cursor: pointer;
border: 1px solid grey;
padding-top: 5px;
width: 160px;
height: 50px;
overflow: hidden;
text-overflow: ellipsis;
white-space: nowrap;
word-break: break-all;
background-color: ${(props) => (props.isSelected?'#cccccc':'#ffffff')};
`

type Props = {
    model?:any
    isSelected?:boolean
    onClick?:()=>void
}

const TestModelCard = ({ model, isSelected=false, onClick }: Props) => {

    const [acc, setAcc] = useState('0')
    const [mse, setMse] = useState('0')

    useEffect(() => {
        if (model.perf) {
            const perf = JSON.parse(model.perf.String)
            if (perf.wa) setAcc(perf.wa)
            if (perf.mse) setMse(perf.mse)
        }
    }, [model])

    return (
        <CardContent isSelected={isSelected} onClick={onClick}>
            <div className="d-flex justify-content-between mb-1">
                {model && <div style={{marginLeft:'10px', marginRight:'10px'}}>
                    <span className="text-muted mb-0">{`Trial${model.train_local_id} `}</span><b>{model.name}</b>
                    {(acc !== '0') && <p className="text-muted mb-0">
                        {`Accuracy ${(Number(acc)*100||0).toFixed(2)}%`}
                    </p>}
                    {(mse !== '0') && <p className="text-muted mb-0">
                        {`MSE ${mse}`}
                    </p>}
                </div>}
            </div>
        </CardContent>
    )
}

export default TestModelCard