import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

import { convertUtcTime, customSortFixedStatus, ApiStopTrain, ApiStopTest, getDurationRealtime, customSortFixedStatusAcc } from 'helpers'
import { CustomConfirm, IconButtonWithPopover, StatusLabel } from 'components'
import { InferenceTimeFormatter, StatusColumnFormatter, TargetMetricFormatter } from 'features/Tables'
import CustomSortIcon from 'features/Tables/ColumnFormatters/CustomSortIcon'
import { getMetricKey } from 'appConstants/trial'

const AccuracySpan = styled.span`
  background-color: rgba(63, 195, 128, 0);
  color: '#ffffff';
`

const ButtonRow = styled.div`
  display: flex;

  & button {
    margin-right: 3px;
    margin-left: 3px;
  }
`

const ColumnModeling = ({
    sort,
    handleSort,
    openAdditionalTrainModal,
    hideAcc,
    accCol
}: any) => {
    const [t] = useTranslation('translation')

    const stopTrain = ApiStopTrain()
    const stopTest = ApiStopTest()

    const handleStopTrain = (row: any) => {
        CustomConfirm({
            onConfirm: () => {
                stopTrain.mutate(row.trial_id)
            },
            onCancel: () => { },
            message: t('ui.confirm.stop.train'),
        })
    }

    const handleStopTest = (row: any) => {
        CustomConfirm({
            onConfirm: () => {
                stopTest.mutate(row.trial_id)
            },
            onCancel: () => { },
            message: t('ui.confirm.stop.test'),
        })
    }

    const columnsA = [
        {
            dataField: 'trial_id',
            text: 'ID',
            headerStyle: { minWidth: '100px', textAlign: 'center' },
            align: () => {
                return 'center'
            },
            sort: true,
            onSort: (field: any, order: any) => handleSort(field, order),
            sortCaret: CustomSortIcon,
            sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
                customSortFixedStatus(a, b, order, dataField, rowA, rowB),
            formatter: (cell: any, row: any) => {
                if (row.parent_id > 0) {
                    return cell + ' (' + row.parent_id + ')'
                }
                return cell
            }
        },
        {
            dataField: 'state-search',
            text: 'state',
            hidden: true,
        },
        {
            dataField: 'created_at',
            text: t('ui.formatter.created'),
            headerStyle: { width: '100px', textAlign: 'center' },
            style: { cursor: 'pointer' },
            formatter: (cell: any, row: any) => {
                return convertUtcTime(cell)
            },
            sort: true,
            onSort: (field: any, order: any) => handleSort(field, order),
            sortCaret: CustomSortIcon,
            sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
                customSortFixedStatus(a, b, order, dataField, rowA, rowB),
        },
        {
            dataField: 'trial_name',
            text: t('ui.formatter.name'),
            sort: true,
            style: { cursor: 'pointer' },
            onSort: (field: any, order: any) => handleSort(field, order),
            sortCaret: CustomSortIcon,
            sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
                customSortFixedStatus(a, b, order, dataField, rowA, rowB),
        },
        {
            dataField: 'test_db',
            text: t('ui.formatter.testdb'),
            headerStyle: { minWidth: '120px', textAlign: 'center' },
            style: { cursor: 'pointer' },
            sort: true,
            onSort: (field: any, order: any) => handleSort(field, order),
            sortCaret: CustomSortIcon,
            sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
                customSortFixedStatus(a, b, order, dataField, rowA, rowB),
        },
    ]

    const columnsB = [
        {
            dataField: 'perf',
            text: t('ui.formatter.perf'),
            headerStyle: { minWidth: '160px', textAlign: 'center' },
            style: { cursor: 'pointer' },
            headerFormatter: (column: any, colIndex: any) => TargetMetricFormatter(column, colIndex, sort, t),
            formatter: (cell: any, row: any) => {
                if (!cell || !cell.target_metric) return '-'
                if (typeof (cell.target_metric) === 'string') {
                    if (!cell[cell.target_metric]) return "-"
                    if (['mse', 'mae', 'rmse'].includes(cell.target_metric)) return Number(cell[cell.target_metric] || 0).toFixed(3) + " (" + t(`metric.${cell.target_metric}`) + ")"
                    return Number(cell[getMetricKey[cell.target_metric]] * 100 || 0).toFixed(2) + "% (" + t(`metric.${cell.target_metric}`) + ")"
                } else if (typeof (cell.target_metric) === 'object') {
                    let key = ''
                    let val = 0
                    Object.keys(cell.target_metric).forEach((k: string) => {
                        if (cell.target_metric[k] > val) {
                            key = k
                            val = cell[k]
                        }
                    })
                    if (['mse', 'mae', 'rmse'].includes(key)) return Number(val || 0).toFixed(3) + " (" + t(`metric.${key}`) + ")"
                    return Number(val * 100 || 0).toFixed(2) + "% (" + t(`metric.${key}`) + ")"
                }
            }
        },
        {
            dataField: 'inference_time',
            text: t('ui.formatter.inftime'),
            headerStyle: { minWidth: '120px', textAlign: 'center' },
            style: { cursor: 'pointer' },
            headerFormatter: (column: any, colIndex: any) => InferenceTimeFormatter(column, colIndex, sort, t),
            align: () => {
                return 'center'
            },
            formatter: (cell: any, row: any) => {
                return cell === null || cell === undefined || cell === '0' ? t('ui.formatter.nonevalue') : (Number(cell) * 1000).toFixed(2) + 'ms'
            },
            sort: true,
            onSort: (field: any, order: any) => handleSort(field, order),
            sortCaret: CustomSortIcon,
            sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
                customSortFixedStatus(a, b, order, dataField, rowA, rowB),
        },
        {
            dataField: 'duration',
            text: t('ui.formatter.duration'),
            headerStyle: { minWidth: '100px', textAlign: 'center' },
            style: { cursor: 'pointer' },
            align: () => {
                return 'center'
            },
            formatter: (cell: any, row: any) => {
                return getDurationRealtime(row.created_at, row.updated_at, row.state)
            },
            sort: true,
            onSort: (field: any, order: any) => handleSort(field, order),
            sortCaret: CustomSortIcon,
            sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
                customSortFixedStatus(getDurationRealtime(rowA.created_at, rowA.updated_at, rowA.state), getDurationRealtime(rowB.created_at, rowB.updated_at, rowB.state), order, dataField, rowA, rowB),
        },
        {
            dataField: 'progress',
            text: t('ui.formatter.progress'),
            headerStyle: { minWidth: '90px', textAlign: 'center' },
            style: { cursor: 'pointer' },
            align: () => {
                return 'center'
            },
            formatter: (cell: any, row: any) => {
                let label = Number(cell || 0).toFixed(1) + '%'
                if (cell === '100') label = '100%'
                if (row.state === 'additional_train') {
                    label = t('ui.formatter.traincount', { count: row.train_count })
                }
                if (row.state === "finish_test" || cell === 100.0) label = '100%'
                return label
            },
            sort: true,
            onSort: (field: any, order: any) => handleSort(field, order),
            sortCaret: CustomSortIcon,
            sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
                customSortFixedStatus(a, b, order, dataField, rowA, rowB),
        },
        {
            dataField: 'state',
            text: '상태',
            headerStyle: { minWidth: '80px', textAlign: 'center' },
            style: { cursor: 'pointer' },
            events: {
                onClick: (e: any, column: any, columnIndex: any, row: any, rowIndex: any) => {
                    e.stopPropagation()
                },
            },
            headerFormatter: (column: any, colIndex: any) => StatusColumnFormatter(column, colIndex, sort, t),
            formatter: (cell: any, row: any, rowIndex: any) => {
                return (
                    <ButtonRow>
                        <StatusLabel state={cell} row={row}>
                            {cell}
                        </StatusLabel>
                        {(row.state === 'train' || row.state === 'additional_train') && (
                            <IconButtonWithPopover
                                name={`btn-stop-train-${rowIndex}`}
                                variant='danger'
                                onClick={() => handleStopTrain(row)}
                                popTitle={t('button.stop.train')}
                                icon='mdi-stop-circle-outline'
                            />
                        )}
                        {row.state === 'test' && (
                            <IconButtonWithPopover
                                name={`btn-stop-test-${rowIndex}`}
                                variant='danger'
                                onClick={() => handleStopTest(row)}
                                popTitle={t('button.stop.test')}
                                icon='mdi-stop-circle-outline'
                            />
                        )}
                        {/* {(row.state === 'finish' || row.state === 'finish-fail') && (row.train_type === 'auto') && (
                            <IconButtonWithPopover
                                name={`btn-finish-${rowIndex}`}
                                variant='info'
                                onClick={() => openAdditionalTrainModal(row)}
                                popTitle={t('button.additional_train')}
                                icon='mdi-sticker-plus-outline'
                            />
                        )} */}
                    </ButtonRow>
                )
            },
            sort: true,
            onSort: (field: any, order: any) => handleSort(field, order),
            sortCaret: CustomSortIcon,
            sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
                customSortFixedStatus(a, b, order, dataField, rowA, rowB),
        },
    ]

    const colAcc = {
        dataField: 'accuracy',
        text: (accCol ? accCol : t('ui.formatter.accuracy')) + '(%)',
        headerStyle: { minWidth: '90px', textAlign: 'center' },
        style: { cursor: 'pointer' },
        align: () => {
            return 'center'
        },
        formatter: (cell: any, row: any) => {
            return (
                <AccuracySpan>{row.perf === null || row.perf === undefined ? t('ui.formatter.nonevalue') : row.perf.wa ? Number(row.perf.wa * 100).toFixed(2) : Number((row.perf.label_accuracy || 0) * 100).toFixed(2)}</AccuracySpan>
            )
        },
        sort: true,
        onSort: (field: any, order: any) => handleSort(field, order),
        sortCaret: CustomSortIcon,
        sortFunc: (a: any, b: any, order: any, dataField: any, rowA: any, rowB: any) =>
            customSortFixedStatusAcc(a, b, order, dataField, rowA, rowB),
    }

    let columns = []
    if (hideAcc !== undefined && hideAcc === true) {
        columns = [...columnsA, ...columnsB]
    } else {
        columns = [...columnsA, colAcc, ...columnsB]
    }

    return columns
}

export default ColumnModeling