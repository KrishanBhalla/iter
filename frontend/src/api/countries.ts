import { HOST } from "./constants"


export class Countries {


    // /countries
    public static async getCountries(): Promise<CountryResponse> {
        // We can use the `Headers` constructor to create headers
        // and assign it as the type of the `headers` variable
        const headers: Headers = new Headers()
        // Add a few headers
        headers.set('Content-Type', 'application/x-www-form-urlencoded')
        headers.set('Accept', 'application/json')

        // Create the request object, which will be a RequestInfo type. 
        // Here, we will pass in the URL as well as the options object as parameters.
        const request: RequestInfo = new Request("http://" + HOST + "/countries", {
            method: 'GET',
            headers: headers
        })

        return await fetch(request)
        // the JSON body is taken from the response
        .then(res => res.json())
        .then(res => {
          return (res as CountryResponse)
        })
    }
}

type CountryResponse = string[]