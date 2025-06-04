import { LabelSelect, PopoverLabel, StatusRow } from "components"
import { useEffect, useState } from "react"
import { useTranslation } from "react-i18next"

const TargetMetricReg = ({ requestData, handleRequestData }: any) => {
    const [t] = useTranslation('translation')

    const [target, setTarget] = useState('target_metric.mse')
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
                        <StatusRow><span>{t(`metric.mse.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.rmse.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.mae.desc`)}</span></StatusRow>
                    </PopoverLabel>
                </span>
            }
            name={'target_metric'}
            onChange={(e: any) => handleRequestData('target_metric', e.target.value)}
            value={target}
        >
            <option key='target_metric_mse' value='target_metric.mse'>
                {t(`metric.mse`)}
            </option>
            <option key='target_metric_rmse' value='target_metric.rmse'>
                {t(`metric.rmse`)}
            </option>
            <option key='target_metric_mae' value='target_metric.mae'>
                {t(`metric.mae`)}
            </option>
        </LabelSelect>
    )
}

export default TargetMetricReg