import { Row } from "react-bootstrap"
import styled from "styled-components"

const ImageArea = styled.div`
min-width: 100px;
max-width: 400px;
min-height: 200px;
padding: 10px;

& img {
    width: 100%;
}
`

const VisionCLSResultImage = ({ origin, overlay }: any) => {
    return (
        <Row>
            {
                origin && <ImageArea>
                    <img src={`data:image/;base64,${origin}`} alt='origin' />
                </ImageArea>
            }
            {
                overlay && <ImageArea><img src={`data:image/;base64,${overlay}`} alt='overlay' /></ImageArea>
            }
        </Row>
    )
}

export default VisionCLSResultImage