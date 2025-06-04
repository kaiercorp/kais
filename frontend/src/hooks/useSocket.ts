import { useEffect, useRef, MutableRefObject, SetStateAction } from 'react'
import config from 'config'
import { logger } from 'helpers'

const useSocket = (
    apiUrl: string,
    content: string,
    onmessage: (e: MessageEvent<any>) => void,
    options: {
        wsAddress?: string 
        setSocketConnected?: (value: SetStateAction<boolean>) => void 
        shouldCleanup?: boolean
        shouldConnect?: boolean
    } = {}
): MutableRefObject<WebSocket | null>  => {
    const ws = useRef<WebSocket | null>(null)
    
    useEffect(() => {
        if (ws.current) {
            ws.current.onmessage = onmessage            
        }
    }, [onmessage])
    
    useEffect(() => {
        if (apiUrl.includes('undefined')) return
        if (!ws.current) {
            if (options.shouldConnect === undefined) return
            if (typeof options.shouldConnect === 'boolean' && !options.shouldConnect) return

            let wsAddress = options.wsAddress
            if (!wsAddress) {
                wsAddress = 'ws://' + window.location.host + config.wsURL
                if (window.location.port === '3000') {
                    wsAddress = 'ws://' + window.location.hostname + config.basePort + config.wsURL
                }
            }
            
            ws.current = new WebSocket(wsAddress + apiUrl)
            ws.current.onopen = () => {
                logger.log(`[${new Date()}] Connected WS - ${content}`)                
                if (options.setSocketConnected) {
                    options.setSocketConnected(true)
                } 
            }
            ws.current.onclose = (error) => {
                // logger.warn(error)
            }
            ws.current.onerror = (error) => {
                logger.warn(error)
            } 
            ws.current.onmessage = onmessage
        }
        
        if (options.shouldCleanup) {
            return () => {
                logger.log(`[${new Date()}] Disconnect WS - ${content}`)
                if (options.setSocketConnected) {
                    options.setSocketConnected(false)
                }
                if (ws.current) {
                    ws.current.close(1000)
                    ws.current = null
                } 
            }
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [apiUrl, options.shouldConnect])

    return ws
}

export default useSocket