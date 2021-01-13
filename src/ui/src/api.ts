import { Cluster } from "./types/cluster"
import { Service } from "./types/service"

export class Api {
    static isLocalhost = Boolean(
        window.location.hostname === 'localhost' ||
        // [::1] is the IPv6 localhost address.
        window.location.hostname === '[::1]' ||
        // 127.0.0.1/8 is considered localhost for IPv4.
        window.location.hostname.match(
            /^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/
        )
    )

    static baseUrl = Api.isLocalhost ? "https://dd0l0av14k.execute-api.us-east-2.amazonaws.com/api/" : "/"

    static async getClusters(): Promise<Cluster[]> {
        return fetch(this.baseUrl + "clusters").then(response => response.json())
    }

    static async getClusterDetails(clusterArn: string): Promise<Service[]> {
        return fetch(this.baseUrl + "clusters?arn=" + encodeURIComponent(clusterArn)).then(response => response.json())
    }
}