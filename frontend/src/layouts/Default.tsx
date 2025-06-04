import { Suspense, useEffect, useContext } from 'react'
import { Outlet } from 'react-router-dom'
import { changeBodyAttribute } from 'helpers'
import { LayoutContext } from 'contexts'

const loading = () => <div className=''></div>

type DefaultLayoutProps = {}

const DefaultLayout = (props: DefaultLayoutProps) => {
  const { layoutContextValue } = useContext(LayoutContext)

  useEffect(() => {
    changeBodyAttribute('data-layout-color', layoutContextValue.layoutColor)
  }, [layoutContextValue.layoutColor])

  return (
    <Suspense fallback={loading()}>
      <Outlet />
    </Suspense>
  )
}
export default DefaultLayout
