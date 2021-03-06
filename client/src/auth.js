
function getCurrentUser() {
    let user = JSON.parse(localStorage.getItem("user"))
    return user
}

export { getCurrentUser }