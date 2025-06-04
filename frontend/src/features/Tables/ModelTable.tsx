import { Card } from 'react-bootstrap'
import BootstrapTable from 'react-bootstrap-table-next'
import paginationFactory from 'react-bootstrap-table2-paginator'
import filterFactory from 'react-bootstrap-table2-filter';
import ToolkitProvider from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit'
import styled from 'styled-components'

const TableArea = styled.div`
  & table th {
    font-weight: 600;
    font-size: 13px;
    padding: 10px;
  }

  & table td {
    padding: 5px 10px 5px 5px;
    vertical-align: middle;
    font-size: 12px;
  }

  & table tr:hover {
    font-weight: 800;
  }

  & table td.selection-cell {
    padding: 10px 12px;
  }

  & table .sortable {
    cursor: pointer;
  }
`

const ModelTable = ({ CustomColumn, filteredModels, onSelect, openMultiTest, downloadModel }: any) => {
    const sample=filteredModels[0]

    const selectRow = {
        mode: 'radio',
        clickToSelect: true,
        style: { backgroundColor: '#c8e6c9' },
        onSelect: (row: any, isSelect: any) => {
            onSelect(row)
        }
    }

    return (
        <ToolkitProvider
            keyField='model_id'
            data={filteredModels}
            columns={CustomColumn({
                sample,
                openMultiTest,
                downloadModel
            })}
            search={{ searchFormatted: true }}
        >
            {(props: any) => (
                <>
                    <Card>
                        <TableArea>
                            <BootstrapTable
                                {...props.baseProps}
                                hover
                                selectRow={selectRow}
                                bordered={false}
                                pagination={paginationFactory({ custom: false, sizePerPage: 10, hideSizePerPage: true })}
                                wrapperClasses='table-responsive'
                                headerClasses='trialtable-header'
                                rowClasses='text-nowrap'
                                filter={filterFactory()}
                            />
                        </TableArea>
                    </Card>
                </>
            )}
        </ToolkitProvider>
    )
}

export default ModelTable