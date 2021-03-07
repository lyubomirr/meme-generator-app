class Endpoints {
    static Login = "/api/login";
    static Register = "/api/register";
    static Memes = "/api/auth/meme";
    static Templates = "/api/auth/template";

    static GetTemplateFileUrl(id) {
        return `/api/template/file/${id}`
    }
}

export default Endpoints