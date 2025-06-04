import { Row } from 'react-bootstrap'
import styled from 'styled-components'

const ImageArea = styled.div`
padding-right: 23px;
& img {
    width: 100%;
}
`

const TableErrorGraphImage = ({ error_graph }: any) => {
    return (
        <Row>
            {
                error_graph && <ImageArea>
                    <img src={`data:image/;base64,${error_graph}`} alt='origin' />
                </ImageArea>
            }
        </Row>
    )
}

export default TableErrorGraphImage