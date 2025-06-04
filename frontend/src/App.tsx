import Routes from 'routes/Routes'

import 'assets/scss/Saas.scss'
import { ApiFetchConfigs } from 'helpers'
import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'


const App = () => {
  const { i18n } = useTranslation()
  const { configs } = ApiFetchConfigs()

  useEffect(() => {
    if (!configs || configs.length < 1) return

    const lang = configs.filter((config: any) => {
      return config.config_key === 'LANGUAGE'
    })[0]

    i18n.changeLanguage(lang.config_val)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [configs])

  return (
    <>
      <Routes />
    </>
  )
}

export default App
