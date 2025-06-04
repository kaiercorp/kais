import { LabelSelect, PopoverLabel, StatusRow } from "components"
import { useEffect, useState } from "react"
import { useTranslation } from "react-i18next"

const TargetMetricSLCls = ({ requestData, handleRequestData}: any) => {
    const [t] = useTranslation('translation')

    const [target, setTarget] = useState('target_metric.wa')
    useEffect(() => {
        if (!requestData || !requestData.train_config || !requestData.train_config.target_metric) return

        const target_metric = requestData.train_config.target_metric
        setTarget(`target_metric.${Object.keys(target_metric)
            .filter((metric: string) => {
                return target_metric[metric] === 100
            })[0]}`)
    }, [requestData])

    return (
        <LabelSelect
            title={
                <span>
                    {t('ui.train.targetmetric')}
                    <PopoverLabel>
                        <StatusRow><span>{t(`metric.uwa.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.wa.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.precision.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.recall.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.f1.desc`)}</span></StatusRow>
                    </PopoverLabel>
                </span>
            }
            name={'target_metric'}
            onChange={(e: any) => handleRequestData('target_metric', e.target.value)}
            value={target}
        >
            <option key='target_metric_wa' value='target_metric.wa'>
                {t(`metric.wa`)}
            </option>
            <option key='target_metric_uwa' value='target_metric.uwa'>
                {t(`metric.uwa`)}
            </option>
            <option key='target_metric_precision' value='target_metric.precision'>
                {t(`metric.precision`)}
            </option>
            <option key='target_metric_recall' value='target_metric.recall'>
                {t(`metric.recall`)}
            </option>
            <option key='target_metric_f1' value='target_metric.f1'>
                {t(`metric.f1`)}
            </option>
        </LabelSelect>
    )
}

export default TargetMetricSLCls