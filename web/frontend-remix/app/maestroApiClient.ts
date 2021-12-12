import {Thing} from "~/model/thing";
import {Error as ErrorModel} from "~/model/error"
import {Service} from "~/model/service";

class MaestroApiClient {
    readonly baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    async searchFromLink(link: string): Promise<Thing[]> {
        const urlSafeLink = encodeURIComponent(link);
        const apiUrl = `${this.baseUrl}/link?link=${urlSafeLink}`

        const res = await fetch(apiUrl)
        const json = await res.json()

        if (res.status != 200) {
            const resErr = json as ErrorModel
            throw new Error(resErr.Error)
        }

        return json as Thing[];
    }

    async getGroup(groupId: string): Promise<Thing[]> {
        const apiUrl = `${this.baseUrl}/${groupId}`

        const res = await fetch(apiUrl)
        const json = await res.json()

        if (res.status != 200) {
            const resErr = json as ErrorModel
            throw new Error(resErr.Error)
        }

        return json as Thing[];
    }

    async getServices(): Promise<Service[]> {
        const apiUrl = `${this.baseUrl}/services`

        const res = await fetch(apiUrl)
        const json = await res.json()

        if (res.status != 200) {
            const resErr = json as ErrorModel
            throw new Error(resErr.Error)
        }

        return json as Service[];
    }
}

export default MaestroApiClient;