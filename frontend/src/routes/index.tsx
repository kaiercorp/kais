import React, { Suspense } from 'react'
import { useRoutes } from 'react-router-dom'

import { OnePageLayout, DefaultLayout } from 'layouts'
import Root from './Root'

const DashboardLanding = React.lazy(() => import('pages/Dashboard/index'))
const ProjectLanding = React.lazy(() => import('pages/Project/index'))
const ConfigurationLanding = React.lazy(() => import('pages/Configuration/Config'))
const MenuLanding = React.lazy(() => import('pages/Configuration/Menu'))
const SystemLanding = React.lazy(() => import('pages/Configuration/System'))
const UserLanding = React.lazy(() => import('pages/Auth/User'))
const HPOLanding = React.lazy(() => import('pages/HPO/HPOLanding'))
const DatasetLanding = React.lazy(() => import('pages/Dataset/DatasetLanding'))

const VisionCLSSLLanding = React.lazy(() => import('pages/Vision/Cls.SL/VisionCLSLanding'))
const VisionCLSSLTrainDetail = React.lazy(() => import('pages/Vision/Cls.SL/VisionCLSTrainDetail'))
const VisionCLSSLTestDetail = React.lazy(() => import('pages/Vision/Cls.SL/VisionCLSTestDetail'))

const VisionCLSMLLanding = React.lazy(() => import('pages/Vision/Cls.ML/VisionCLSLanding'))
const VisionCLSMLTrainDetail = React.lazy(() => import('pages/Vision/Cls.ML/VisionCLSTrainDetail'))
const VisionCLSMLTestDetail = React.lazy(() => import('pages/Vision/Cls.ML/VisionCLSTestDetail'))

const VisionADLanding = React.lazy(() => import('pages/Vision/AD/VisionADLanding'))
const VisionADTrainDetail = React.lazy(() => import('pages/Vision/AD/VisionADTrainDetail'))
const VisionADTestDetail = React.lazy(() => import('pages/Vision/AD/VisionADTestDetail'))

const TableCLSLanding = React.lazy(() => import('pages/TableCls/TableCLSLanding'))
const TableCLSTrainDetail = React.lazy(() => import('pages/TableCls/TableCLSTrainDetail'))
const TableCLSTestDetail = React.lazy(() => import('pages/TableCls/TableCLSTestDetail'))

const TableREGLanding = React.lazy(() => import('pages/TableReg/TableREGLanding'))
const TableREGTrainDetail = React.lazy(() => import('pages/TableReg/TableREGTrainDetail'))
const TableREGTestDetail = React.lazy(() => import('pages/TableReg/TableREGTestDetail'))

const TSADLanding = React.lazy(() => import('pages/TSAD/TSADLanding'))
const TSADTrainDetail = React.lazy(() => import('pages/TSAD/TSADTTrainDetail'))
const TSADTestDetail = React.lazy(() => import('pages/TSAD/TSADTestDetail'))

const Login = React.lazy(() => import('pages/Auth/Login'))
const PageNotFound = React.lazy(() => import('pages/error/PageNotFound'))
const ServerError = React.lazy(() => import('pages/error/ServerError'))

const loading = () => <div className='loading' />

type LoadComponentProps = {
  component: React.LazyExoticComponent<() => JSX.Element>
}

const LoadComponent = ({ component: Component }: LoadComponentProps) => (
  <Suspense fallback={loading()}>
    <Component />
  </Suspense>
)

const AllRoutes = () => {
  return useRoutes([
    {
      path: '/',
      element: <Root />,
    },
    {
      path: '/',
      element: <DefaultLayout />,
      children: [
        {
          path: '/auth/login',
          element: <LoadComponent component={Login} />,
        },
        {
          path: '/auth/logout',
          element: <LoadComponent component={Login} />,
        }
      ]
    },
    {
      path: '/',
      element: <OnePageLayout />,
      children: [
        {
          path: 'dashboard',
          element: <LoadComponent component={DashboardLanding} />,
        },
        {
          path: 'error-404',
          element: <LoadComponent component={PageNotFound} />,
        },
        {
          path: 'error-500',
          element: <LoadComponent component={ServerError} />,
        },
        {
          path: 'dataset',
          element: <LoadComponent component={DatasetLanding} />,
        },
        {
          path: 'hpo',
          element: <LoadComponent component={HPOLanding} />,
        }
      ]
    },
    {
      path: '/vision',
      element: <OnePageLayout />,
      children: [
        {
          path: 'cls-sl',
          element: <LoadComponent component={ProjectLanding} />,
        },
        {
          path: 'cls-sl/:project_id',
          element: <LoadComponent component={VisionCLSSLLanding} />,
        },
        {
          path: 'cls-sl/:project_id/:id/train',
          element: <LoadComponent component={VisionCLSSLTrainDetail} />,
        },
        {
          path: 'cls-sl/:project_id/:id/test',
          element: <LoadComponent component={VisionCLSSLTestDetail} />,
        },
        {
          path: 'cls-ml',
          element: <LoadComponent component={ProjectLanding} />,
        },
        {
          path: 'cls-ml/:project_id',
          element: <LoadComponent component={VisionCLSMLLanding} />,
        },
        {
          path: 'cls-ml/:project_id/:id/train',
          element: <LoadComponent component={VisionCLSMLTrainDetail} />,
        },
        {
          path: 'cls-ml/:project_id/:id/test',
          element: <LoadComponent component={VisionCLSMLTestDetail} />,
        },
        {
          path: 'ad',
          element: <LoadComponent component={ProjectLanding} />,
        },
        {
          path: 'ad/:project_id',
          element: <LoadComponent component={VisionADLanding} />,
        },
        {
          path: 'ad/:project_id/:id/train',
          element: <LoadComponent component={VisionADTrainDetail} />,
        },
        {
          path: 'ad/:project_id/:id/test',
          element: <LoadComponent component={VisionADTestDetail} />,
        },
      ],
    },
    {
      path: '/table',
      element: <OnePageLayout />,
      children: [
        {
          path: 'cls',
          element: <LoadComponent component={ProjectLanding} />,
        },
        {
          path: 'cls/:project_id',
          element: <LoadComponent component={TableCLSLanding} />,
        },
        {
          path: 'cls/:project_id/:id/train',
          element: <LoadComponent component={TableCLSTrainDetail} />,
        },
        {
          path: 'cls/:project_id/:id/test',
          element: <LoadComponent component={TableCLSTestDetail} />,
        },
        {
          path: 'reg',
          element: <LoadComponent component={ProjectLanding} />,
        },
        {
          path: 'reg/:project_id',
          element: <LoadComponent component={TableREGLanding} />,
        },
        {
          path: 'reg/:project_id/:id/train',
          element: <LoadComponent component={TableREGTrainDetail} />,
        },
        {
          path: 'reg/:project_id/:id/test',
          element: <LoadComponent component={TableREGTestDetail} />,
        }
      ],
    },
    {
      path: '/ts',
      element: <OnePageLayout />,
      children: [
        {
          path: 'ad',
          element: <LoadComponent component={ProjectLanding} />,
        },
        {
          path: 'ad/:project_id',
          element: <LoadComponent component={TSADLanding} />,
        },
        {
          path: 'ad/:project_id/:id/train',
          element: <LoadComponent component={TSADTrainDetail} />,
        },
        {
          path: 'ad/:project_id/:id/test',
          element: <LoadComponent component={TSADTestDetail} />,
        },
      ],
    },
    {
      path: '/config',
      element: <OnePageLayout />,
      children: [
        {
          path: 'setting',
          element: <LoadComponent component={ConfigurationLanding} />,
        },
        {
          path: 'menu',
          element: <LoadComponent component={MenuLanding} />,
        },
        {
          path: 'user',
          element: <LoadComponent component={UserLanding} />,
        },
        {
          path: 'system',
          element: <LoadComponent component={SystemLanding} />,
        },
      ]
    }
  ])
}

export { AllRoutes }
