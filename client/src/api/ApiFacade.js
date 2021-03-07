import { getCurrentUser } from "../auth";
import Endpoints from "./Endpoints"

class ApiFacade {
    static get(url, auth = false) {
        let jwt = ""
        if(auth) {
            const user = getCurrentUser()
            jwt = user.jwt
        }

        return new Promise((resolve, reject) => {
            fetch(url, {
                    method: "GET",
                    headers: {
                        "Content-type": "application/json; charset=UTF-8",
                        "Authorization": "Bearer " + jwt
                    }
                })
                .then(response => response.json())
                .then(jsonResponse => {
                    if (jsonResponse.hasOwnProperty("errorMessage")) {
                        reject(jsonResponse.errorMessage);
                    } else {
                        resolve(jsonResponse);
                    }
                })
                .catch(err => reject(err));
        })
    }

    static postJson(url, body) {
        return new Promise((resolve, reject) => {
            fetch(url, {
                    method: "POST",
                    body: JSON.stringify(body),
                    headers: {
                        "Content-type": "application/json; charset=UTF-8"
                    }
                })
                .then(response => response.json())
                .then(jsonResponse => {
                    if (jsonResponse.hasOwnProperty("errorMessage")) {
                        reject(jsonResponse.errorMessage);
                    } else {
                        resolve(jsonResponse);
                    }
                })
                .catch(err => reject(err));
        })
    }

    static postFormData(url, formData, auth=false) {
        let jwt = ""
        if(auth) {
            const user = getCurrentUser()
            jwt = user.jwt
        }

        return new Promise((resolve, reject) => {
            fetch(url, {
                    method: "POST",
                    body: formData,
                    headers: {
                        "Authorization": "Bearer " + jwt
                    }
                })
                .then(response => response.json())
                .then(jsonResponse => {
                    if (jsonResponse.hasOwnProperty("errorMessage")) {
                        reject(jsonResponse.errorMessage);
                    } else {
                        resolve(jsonResponse);
                    }
                })
                .catch(err => reject(err));
        })
    }

    static login(loginData) {
        return this.postJson(Endpoints.Login, loginData)
    }

    static register(registerData) {
        return this.postJson(Endpoints.Register, registerData)
    } 

    static getMemes() { 
        return this.get(Endpoints.Memes, true)
    }

    static getTemplates() { 
        return this.get(Endpoints.Templates, true)
    }

    static getTemplate(id) {
        return this.get(`${Endpoints.Templates}/${id}`, true)
    }

    static createMeme(meme, fileBlob) {
        const formData = new FormData();
        formData.append("file", fileBlob)
        formData.append("meme", JSON.stringify(meme))

        return this.postFormData(Endpoints.Memes, formData, true)
    }
}

export default ApiFacade