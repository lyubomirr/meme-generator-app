class Endpoints {
    static Login = "/api/login";
    static Register = "/api/register";
    static Memes = "/api/auth/meme";
    static Templates = "/api/auth/template";
    static AdminTemplatePath = "/api/auth/admin/template"

    static GetTemplateFileUrl(id) {
        return `/api/template/file/${id}`
    }

    static GetMemeFileUrl(id) {
        return `/api/meme/file/${id}`
    }

    static GetCommentUrl(memeId) {
        return `/api/auth/meme/${memeId}/comment`
    }
}

export default Endpoints