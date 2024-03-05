import _ from 'lodash';
import React, { PropsWithChildren, createContext, useContext } from 'react';
import { UrlData, Evidence } from 'src/global_types';

interface IEvidencesContext {
	imgDataSetter: (urlData: UrlData | null) => void
	imgData: UrlData | null
	activeEvidence?: Evidence
}

export const EvidencesContext = createContext<IEvidencesContext>({
	imgData: null,
	imgDataSetter: () => 0,
})

export const useEvidenceContext = () => {
	return useContext(EvidencesContext)
}

interface EvidencesContextProviderProps {
	activeEvidence: Evidence
}

const EvidencesContextProvider: React.FC<PropsWithChildren<EvidencesContextProviderProps>> = ({ children, activeEvidence }) => {
	const [currImageData, setCurrImageData] = React.useState<UrlData| null>(null)

	return (
		<EvidencesContext.Provider value={{
			imgData: currImageData,
			// This should prevent possible accidental state updates, since although
			// the object could be equal, the reference is different, which would trigger
			// a state update
			imgDataSetter: (newData) => setCurrImageData(currData => {
				if (_.isEqual(newData, currData)) {
					return currData
				}

				return newData
			}),
			activeEvidence
		}}>
			{children}
		</EvidencesContext.Provider>
	);
}

export default EvidencesContextProvider;
