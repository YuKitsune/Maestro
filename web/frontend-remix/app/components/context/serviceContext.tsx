import {createContext, PropsWithChildren, useCallback, useContext, useMemo, useState} from "react";
import {Service} from "~/model/service";

type ServiceContextState = {
    getService: (key: string) => Promise<Service | undefined>;
}

const ServiceContext = createContext<ServiceContextState>(null!);

const useServiceContext = () => {
    return useContext(ServiceContext);
}

type ServiceContextProviderProps = PropsWithChildren<{}>;

export const ServiceContextProvider = (props: ServiceContextProviderProps) => {
    const [stale, setStale] = useState(false)

    const getServices = useCallback(async () => {

        const apiUrl = `${process.env.API_URL}/services`
        const res = await fetch(apiUrl)
        const json = await res.json()

        setStale(false)

        return json as Service[];
    }, []);

    let services = useMemo(async () => await getServices(), [stale])

    const findService = useCallback(async (key: string) => {
        let availableServices = await services

        let service = availableServices.find(s => s.Key === key)
        if (!service && stale) {

            // Mark as stale and try again
            setStale(true)
            service = availableServices.find(s => s.Key === key)
        }

        if (!service) {
            return undefined
        }

        return service;
    }, [services, stale])

    return <ServiceContext.Provider value={{ getService: findService }}>
        {props.children}
    </ServiceContext.Provider>
}

export default useServiceContext;
