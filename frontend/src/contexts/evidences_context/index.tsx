import _ from 'lodash';
import React, { PropsWithChildren, createContext, useContext } from 'react';
import { UrlData, Evidence } from 'src/global_types';

interface IEvidencesContext {
	imgDataSetter: (key: string, urlData: UrlData) => void
	cachedUrls: Map<string, UrlData>
}

export const EvidencesContext = createContext<IEvidencesContext>({
	imgDataSetter: () => 0,
	cachedUrls: new Map()
})

export const useEvidenceContext = () => {
	return useContext(EvidencesContext)
}

interface EvidencesContextProviderProps {
	activeEvidence: Evidence
}

const EvidencesContextProvider: React.FC<PropsWithChildren<EvidencesContextProviderProps>> = ({ children, activeEvidence }) => {
	const [cachedUrls, setCachedUrls] = React.useState<Map<string, UrlData>>(new Map())

	return (
		<EvidencesContext.Provider value={{
			imgDataSetter: (key, urlData) => {
				const newCachedUrls = new Map(cachedUrls)

				newCachedUrls.set(key, urlData)

				setCachedUrls(newCachedUrls)
			},
			cachedUrls
		}}>
			{children}
		</EvidencesContext.Provider>
	);
}

export default EvidencesContextProvider;
