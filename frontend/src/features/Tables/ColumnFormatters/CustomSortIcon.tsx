import styled from 'styled-components'

const SortIcon = styled.span`
font-size: 15px;
margin-left: 10px;
& i {
  margin-left: -10px;
}
`

const CustomSortIcon = (order: any, column: any) => {
    return (
      <SortIcon>
        {order !== 'desc' && <i className='mdi mdi-arrow-up-thin' />}
        {order !== 'asc' && <i className='mdi mdi-arrow-down-thin' />}
      </SortIcon>
    )
  }

export default CustomSortIcon