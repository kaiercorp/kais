import i18n from 'i18next'
import detector from 'i18next-browser-languagedetector'
import { initReactI18next } from 'react-i18next'

import translationEn from 'locales/en/translation.json'
import translationKo from 'locales/ko/translation.json'

const resources = {
    en: {
        translation: translationEn
    },
    ko: {
        translation: translationKo
    }
}

i18n.use(detector)
    .use(initReactI18next)
    .init({
        resources,
        lng: 'ko',
        keySeparator: false,
        interpolation: {
            escapeValue: false
        }
    })

export default i18n