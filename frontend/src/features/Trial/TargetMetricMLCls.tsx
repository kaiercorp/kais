import { LabelSelect, PopoverLabel, StatusRow } from "components"
import { useEffect, useState } from "react"
import { useTranslation } from "react-i18next"

const TargetMetricMLCls = ({ requestData, handleRequestData}: any) => {
    const [t] = useTranslation('translation')

    const [target, setTarget] = useState('target_metric.image_f1_score')
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
                        <StatusRow><span>{t(`metric.image_accuracy.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.image_precision.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.image_recall.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.image_f1_score.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.label_accuracy.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.label_precision.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.label_recall.desc`)}</span></StatusRow>
                        <StatusRow><span>{t(`metric.label_f1_score.desc`)}</span></StatusRow>
                    </PopoverLabel>
                </span>
            }
            name={'target_metric'}
            onChange={(e: any) => handleRequestData('target_metric', e.target.value)}
            value={target}
        >
            <option key='target_metric_imageAccuracy' value='target_metric.image_accuracy'>
                {t(`metric.image_accuracy`)}
            </option>

            <option key='target_metric_imagePrecision' value='target_metric.image_precision'>
                {t(`metric.image_precision`)}
            </option>

            <option key='target_metric_imageRecall' value='target_metric.image_recall'>
                {t(`metric.image_recall`)}
            </option>

            <option key='target_metric_imageF1' value='target_metric.image_f1_score'>
                {t(`metric.image_f1_score`)}
            </option>

            <option key='target_metric_labelAccuracy' value='target_metric.label_accuracy'>
                {t(`metric.label_accuracy`)}
            </option>

            <option key='target_metric_labelPrecision' value='target_metric.label_precision'>
                {t(`metric.label_precision`)}
            </option>

            <option key='target_metric_labelRecall' value='target_metric.label_recall'>
                {t(`metric.label_recall`)}
            </option>

            <option key='target_metric_labelF1' value='target_metric.label_f1_score'>
                {t(`metric.label_f1_score`)}
            </option>
        </LabelSelect>
    )
}

export default TargetMetricMLCls