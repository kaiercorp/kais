import ConfigCol from "features/TrialDetail/ConfigCol"
import { useEffect } from "react"
import { Card } from "react-bootstrap"

const PredictionProbabilitiesREG = ({data}: any) => {
    useEffect(() => {
       if (!data)  return
    }, [data])
    return (
        <Card>
            <Card.Header>Prediction probabilities</Card.Header>
            <Card.Body>
                {
                    data.target&&(
                        <ConfigCol>
                            <table>
                                <thead>
                                    <tr>
                                        <td>Target Column</td>
                                        {/* <td>Target Value</td> */}
                                        <td>Prediction</td>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>{data.target.key}</td>
                                        {/* <td>{data.target.value}</td> */}
                                        <td>{data.prediction}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </ConfigCol>
                    )
                }
            </Card.Body>
        </Card>
    )
}

export default PredictionProbabilitiesREG