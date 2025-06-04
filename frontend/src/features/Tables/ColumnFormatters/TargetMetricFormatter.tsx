import { PopoverLabel } from "components"

const TargetMetricFormatter = (column: any, colIndex: any, sort: any, t: any) => {
    return (
        <span>
            {column.text}
            <PopoverLabel name='perf'>
                {t('ui.formatter.perf.desc')}
            </PopoverLabel>
        </span>
    )
}

export default TargetMetricFormatter