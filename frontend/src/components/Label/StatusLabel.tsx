import { stateColors } from 'appConstants'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'

interface Props {
  state: string
}

const StyledDiv = styled.div<Props>`
  color: ${(props) => stateColors[props.state] || 'black'};
  & > i {
    color: ${(props) => stateColors[props.state] || 'black'};
    margin-right: 2px;
  }
`

const StatusLabel = (props: any) => {
  const [t] = useTranslation('translation')
  return (
    <StyledDiv {...props}>
      {props.state === 'finish-fail' ? (
        <i className='mdi mdi-alert-circle-outline'></i>
      ) : (
        <i className='mdi mdi-checkbox-blank-circle'></i>
      )}
      {t(`state.${props.state}`)}
    </StyledDiv>
  )
}

export default StatusLabel
