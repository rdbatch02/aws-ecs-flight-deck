export class Api {
    static baseUrl = "https://dd0l0av14k.execute-api.us-east-2.amazonaws.com"

    static async getClusters(): Promise<any> {
        return fetch(this.baseUrl + "/clusters").then(response => response.json())
    }
}