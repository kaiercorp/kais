import { ReactNode } from 'react'
import { LocationContextProvider, ProjectContextProvider, LayoutContextProvider, DiskContextProvider, GPUContextProvider, TrialContextProvider, FilterContextProvider } from './' 

export const AppContextProvider = ({children}: {children: ReactNode}) => {
    return (
        <LocationContextProvider>
            <ProjectContextProvider>
                <LayoutContextProvider>
                    <DiskContextProvider>
                        <GPUContextProvider>
                            <TrialContextProvider>
                                <FilterContextProvider>
                                    {children}
                                </FilterContextProvider>
                            </TrialContextProvider>
                        </GPUContextProvider>
                    </DiskContextProvider>
                </LayoutContextProvider>
            </ProjectContextProvider>
        </LocationContextProvider>
    )
}