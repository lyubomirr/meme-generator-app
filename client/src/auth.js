
function getCurrentUser() {
    let user = JSON.parse(localStorage.getItem("user"))
    return user
}

function saveUser(userObject) {
    localStorage.setItem("user", JSON.stringify(userObject))
}

function clearUser() {
    localStorage.removeItem("user")
}

export { saveUser, getCurrentUser, clearUser }