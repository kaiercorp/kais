import { createRoot } from 'react-dom/client'
import { QueryClientProvider, QueryClient, QueryCache } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'

import './i18n'
import App from './App'
import { AppContextProvider } from 'contexts'
import reportWebVitals from './reportWebVitals'
import './kais.css'

const domNode = document.getElementById('root')

if (domNode) {
  const root = createRoot(domNode)
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
        refetchOnMount: false,
        refetchOnReconnect: false,
      },
    },
    queryCache: new QueryCache({
      onError: (error, query) => {
        // if (query.meta.errorMessage) toast(() => query.meta.errorMessage);
      },
    }),
  })

  root.render(
    <QueryClientProvider client={queryClient}>
      <ReactQueryDevtools initialIsOpen={false} />
      <AppContextProvider>
        <App />
      </AppContextProvider>
    </QueryClientProvider>
  )
}

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
