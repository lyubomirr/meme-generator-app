import { getCurrentUser } from "../auth";
import Endpoints from "./Endpoints"

class ApiFacade {
    static get(url, auth = false) {
        let jwt = "";
        if(auth) {
            const user = getCurrentUser();
            jwt = user.jwt;
        }

        return new Promise((resolve, reject) => {
            fetch(url, {
                    method: "GET",
                    headers: {
                        "Content-type": "application/json; charset=UTF-8",
                        "Authorization": "Bearer " + jwt
                    }
                })
                .then(response => response.json().catch(() => {
                    return {}
                }))
                .then(jsonResponse => {
                    if (jsonResponse.hasOwnProperty("errorMessage")) {
                        reject({message: jsonResponse.errorMessage});
                    } else {
                        resolve(jsonResponse);
                    }
                })
                .catch(err => reject(err));
        })
    }

    static postJson(url, body, auth=false) {
        let jwt = "";
        if(auth) {
            const user = getCurrentUser();
            jwt = user.jwt;
        }

        return new Promise((resolve, reject) => {
            fetch(url, {
                    method: "POST",
                    body: JSON.stringify(body),
                    headers: {
                        "Content-type": "application/json; charset=UTF-8",
                        "Authorization": "Bearer " + jwt
                    }
                })
                .then(response => response.json().catch(() => {
                    return {}
                }))
                .then(jsonResponse => {
                    if (jsonResponse.hasOwnProperty("errorMessage")) {
                        reject({message: jsonResponse.errorMessage});
                    } else {
                        resolve(jsonResponse);
                    }
                })
                .catch(err => reject(err));
        })
    }

    static postFormData(url, formData, auth=false) {
        let jwt = "";
        if(auth) {
            const user = getCurrentUser();
            jwt = user.jwt;
        }

        return new Promise((resolve, reject) => {
            fetch(url, {
                    method: "POST",
                    body: formData,
                    headers: {
                        "Authorization": "Bearer " + jwt
                    }
                })
                .then(response => response.json().catch(() => {
                    return {}
                }))
                .then(jsonResponse => {
                    if (jsonResponse.hasOwnProperty("errorMessage")) {
                        reject({message: jsonResponse.errorMessage});
                    } else {
                        resolve(jsonResponse);
                    }
                })
                .catch(err => reject(err));
        })
    }

    static delete(url, auth = false) {
        let jwt = "";
        if(auth) {
            const user = getCurrentUser();
            jwt = user.jwt;
        }

        return new Promise((resolve, reject) => {
            fetch(url, {
                    method: "DELETE",
                    headers: {
                        "Content-type": "application/json; charset=UTF-8",
                        "Authorization": "Bearer " + jwt
                    }
                })
                .then(response => response.json().catch(() => {
                    return {}
                }))
                .then(jsonResponse => {
                    if (jsonResponse.hasOwnProperty("errorMessage")) {
                        reject({message: jsonResponse.errorMessage});
                    } else {
                        resolve(jsonResponse);
                    }
                })
                .catch(err => reject(err));
        })
    }


    static login(loginData) {
        return this.postJson(Endpoints.Login, loginData);
    }

    static register(registerData) {
        return this.postJson(Endpoints.Register, registerData);
    } 

    static getMemes() { 
        return this.get(Endpoints.Memes, true);
    }

    static getTemplates() { 
        return this.get(Endpoints.Templates, true);
    }

    static getTemplate(id) {
        return this.get(`${Endpoints.Templates}/${id}`, true);
    }

    static createTemplate(template, fileBlob) {    
        const formData = new FormData();
        formData.append("file", fileBlob)
        formData.append("template", JSON.stringify(template))

        return this.postFormData(Endpoints.AdminTemplatePath, formData, true);    
    }

    static deleteTemplate(id) {
        return this.delete(`${Endpoints.AdminTemplatePath}/${id}`, true);
    }

    static createMeme(meme, fileBlob) {
        const formData = new FormData();
        formData.append("file", fileBlob)
        formData.append("meme", JSON.stringify(meme))

        return this.postFormData(Endpoints.Memes, formData, true);
    }

    static addComment(memeId, comment) {
        return this.postJson(Endpoints.GetCommentUrl(memeId), comment, true);
    }

    static deleteComment(memeId, commentId) {
        let url = Endpoints.GetCommentUrl(memeId);
        url = url + "/" + commentId;
        return this.delete(url, true)
    }

    static deleteMeme(memeId) {
        let url = `${Endpoints.Memes}/${memeId}`;
        return this.delete(url, true)
    }
}

export default ApiFacade