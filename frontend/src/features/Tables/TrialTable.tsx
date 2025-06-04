import { useState, useContext } from 'react'
import { Card, Col } from 'react-bootstrap'
import BootstrapTable from 'react-bootstrap-table-next'
import paginationFactory from 'react-bootstrap-table2-paginator'
import ToolkitProvider from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit'
import { useLocation, useNavigate } from 'react-router-dom'
import styled from 'styled-components'
import TableButtons from './TableButtons'
import { TrialContext } from 'contexts'

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

const TrialTable = ({ CustomColumn, filteredTrials, openCompareModal, openFilterModal, openTrainDetailModal, openTestDetailModal, openAdditionalTrainModal, openFailModal, hideAcc, accCol }: any) => {
    const { trialContextValue, updateTrialContextValue } = useContext(TrialContext)

    const navigate = useNavigate()
    const location = useLocation()

    const selectRow = {
        mode: 'checkbox',
        clickToSelect: true,
        style: { backgroundColor: '#c8e6c9' },
        onSelect: (row: any, isSelect: any) => {
            let newSelected: any = trialContextValue.selectedRows? trialContextValue.selectedRows.slice() : []
            if (isSelect) {
                newSelected.push(row)
            } else {
                newSelected = newSelected.filter((r: any) => r.trial_id !== row.trial_id)
            }
            updateTrialContextValue({selectedRows: newSelected})
        },
        onSelectAll: (isSelect: any, rows: any) => {
            if (isSelect) {
                updateTrialContextValue({selectedRows: rows})
            } else {
                updateTrialContextValue({selectedRows: []})
            }
        },
    }

    const [sort, setSort] = useState({ dataField: 'customOrder', order: 'desc' })
    const handleSort = (field: string, order: string) => {
        setSort({ dataField: field, order: order })
    }

    const rowEvents = {
        onDoubleClick: (e: any, row:any, rowIndex:any) => {
            if (e.target.localName === 'input') {
                return
            }

            if (['train', 'additional_train', 'finish', 'cancel', 'fail'].includes(row.state)) {
                navigate(`${location.pathname}/${row.trial_id}/train`)
            } else if (['test', 'finish_test', 'test_cancel', 'test_fail'].includes(row.state)) {
                navigate(`${location.pathname}/${row.trial_id}/test`)
            }
        },
    }

    return (
       <Col>
        <ToolkitProvider
            keyField='trial_local_id'
            data={filteredTrials}
            columns={CustomColumn({
                sort,
                handleSort,
                openTrainDetailModal,
                openTestDetailModal,
                openAdditionalTrainModal,
                openFailModal,
                hideAcc,
                accCol
            })}
            search={{ searchFormatted: true }}
        >
            {(props: any) => (
                <>
                    <TableButtons {...props.searchProps} openCompareModal={openCompareModal} openFilterModal={openFilterModal} />
                    <Card>
                        <TableArea>
                            <BootstrapTable
                                {...props.baseProps}
                                hover
                                selectRow={selectRow}
                                sort={sort}
                                bordered={false}
                                pagination={paginationFactory({ custom: false, sizePerPage: 15, hideSizePerPage: true })}
                                wrapperClasses='table-responsive'
                                headerClasses='trialtable-header'
                                rowClasses='text-nowrap'
                                rowEvents={rowEvents}
                            />
                        </TableArea>
                    </Card>
                </>
            )}
        </ToolkitProvider>
       </Col>
    )
}

export default TrialTable