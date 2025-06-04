import classNames from 'classnames'
import CustomAlert from 'components/CustomAlert'
import { IconLabel } from 'components/Label'
import { useEffect, useState, useContext } from 'react'
import { Button, Card, Col } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'
import { DiskContext } from 'contexts'

const StyledCard = styled(Card)`
  margin-right: ${(props) => (props.last === 'true' ? '0' : '10px')} !important;

  & button:hover {
    background-color: #464f5b !important;
    color: #f27a7a !important;
    border: 1px solid #ffffff
  }
`

type OneClickButtonProps = {
  icon: string
  title: string
  subTitle: string
  marginLast?: boolean
  onClick: () => void
}

const OneClickButton = ({ icon, title, subTitle, marginLast, onClick }: OneClickButtonProps) => {
  const [t] = useTranslation('translation')
  const {diskContextValue} = useContext(DiskContext)

  const [alertMsg, setAlertMsg] = useState('Workspace 디스크의 용량이 부족합니다.')
  const [limit, setLimit] = useState(0)

  useEffect(() => {
    if (title === t('oneclick.autotrain')) {
      setLimit(20)
      setAlertMsg(t('ui.alert.disklimit', {limit: 20}))
    } else if (title === t('oneclick.foldertest')) {
      setLimit(1)
      setAlertMsg(t('ui.alert.disklimit', {limit: 10}))
    } else if (title === t('oneclick.multitest')) {
      setLimit(1)
      setAlertMsg(t('ui.alert.disklimit', {limit: 10}))
    } else if (title === t('oneclick.singletest')) {
      setLimit(-1)
    } else if (title === t('oneclick.filetest')) {
      setLimit(-1)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [title])
  
  const disableByDisk = () => {
    if (!diskContextValue.disks || diskContextValue.disks.length < 1) return true
    if (limit < 0) return

    let returnVal = false
    diskContextValue.disks.forEach((disk:any) => {
      if (disk.path === "ROOT_PATH") {
        if (Number(disk.free) < limit) {
          returnVal = true
          return true
        }
      }
    })

    return returnVal
  }

  const handleClick = (e: any) => {
    const limited = disableByDisk()
    if (limited) {
      CustomAlert({
        message: alertMsg
      })
      return
    }

    onClick()
  }
  
  return (
    <Col>
      <StyledCard className={classNames('m-0')} last={marginLast ? 'true' : 'false'}>
        <Button className={'p-0'} variant='light' onClick={handleClick}>
          <Card.Body className='text-center' style={{padding: '15px 10px'}}>
            <IconLabel icon={icon} title={title} />
            <p className='text-muted font-15 mb-0'>{subTitle}</p>
          </Card.Body>
        </Button>
      </StyledCard>
    </Col>
  )
}

export default OneClickButton
