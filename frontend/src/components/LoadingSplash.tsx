import { Spinner } from 'components'
import styled from 'styled-components'

type Props = {
  loading: boolean
}

interface IWrapper {
  loading: string
}

const Wrapper = styled.div<IWrapper>`
  width: 100%;
  height: 100%;
  z-index: 1200;
  position: fixed;
  padding-top: 20%;
  top: 0;
  text-align: center;
  background-color: rgba(170, 184, 197, 0.4);
  display: ${(props) => (props.loading === 'true' ? 'block' : 'none')};
`

const LoadingSplash = ({ loading }: Props) => {
  return (
    <Wrapper loading={loading.toString()}>
      <Spinner type='grow' />
    </Wrapper>
  )
}

export default LoadingSplash
