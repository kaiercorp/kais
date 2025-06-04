import { Card, Col, OverlayTrigger, ProgressBar, Row, Tooltip } from "react-bootstrap"
import styled from "styled-components"

const DiskContainer = styled(Card)`
    margin-top: 10px;
    margin-right: 10px;
    background-color: #464f5b;
    padding-right: 5px;
    padding-left: 5px;
    min-width: 200px;
`

const StyledHeader = styled(Row)`
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
`

const StyledColLeft = styled(Col)`
    width: 30px;
    max-width: 30px;
`

const StyledI = styled.i`
    line-height: 35px;
    font-size: 25px;
`

const getPathname = (full_path: string) => {
    if (full_path.includes("/")) {
        return full_path.split("/").slice(-1)
    } else if (full_path.includes("\\")) {
        return full_path.split("\\").slice(-1)
    }

    return full_path
}

const DiskInfo = ({ diskInfo }: any) => {
    const used = diskInfo.usedPercent.toFixed(2)
    const renderTooltip = (props: any) => {
        return (
            <Tooltip id={`PopoverFocus-${diskInfo.path}`} {...props}>
                {diskInfo.path}
            </Tooltip>
        )
    }
    const renderPercent = (props: any) => {
        return (
            <Tooltip id={`PopoverFocus-${diskInfo.path}`} {...props}>
                USED: {used}%
            </Tooltip>
        )
    }
    return (
        <DiskContainer>
            <Row>
                <StyledColLeft>
                    <StyledI className='mdi mdi-harddisk'></StyledI>
                </StyledColLeft>
                <Col>
                    <StyledHeader>
                        <OverlayTrigger placement='bottom' overlay={renderTooltip}>
                            <h6 style={{ marginTop: '3px', marginBottom: '2px' }}>{getPathname(diskInfo.path)}</h6>
                        </OverlayTrigger>
                    </StyledHeader>
                    <Row>
                        <Col>
                            <OverlayTrigger placement='bottom' overlay={renderPercent}>
                                <ProgressBar now={used} style={{backgroundColor: '#829ec2'}} className="progress-md" />
                            </OverlayTrigger>
                        </Col>
                    </Row>
                </Col>
                <Col>
                    <Row>free</Row>
                    <Row style={{ fontSize: '12px' }}>{diskInfo.free} GB</Row>
                </Col>
            </Row>
        </DiskContainer>
    )
}

export default DiskInfo