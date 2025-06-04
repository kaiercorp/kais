import classNames from 'classnames'
import styled from 'styled-components'

const StyledLabel = styled.div`
  display: flex;
  justify-content: center;
  margin-right: 5px;

  & i {
    margin-right: 3px;
  }
`

const IconLabel = ({fontSize, icon, title}: any) => {
  return (
    <StyledLabel>
      <i className={classNames(fontSize ? fontSize : 'font-24', icon ? icon : '')}></i>
      <h5>
        <span>{title}</span>
      </h5>
    </StyledLabel>
  )
}

export default IconLabel
