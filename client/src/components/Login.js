import { useState } from 'react';

const Login = () => {
    const [loginData, setloginData] = useState({
        username: "",
        password: ""
    })
    
    var handleSubmit = (e) => {
        e.preventDefault()
        console.log(loginData)
    }

    return (
        <section className="col-12 col-sm-6 offset-sm-3 login-component shadow">
            <h3 className="mb-4">Login</h3>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                  <label htmlFor="username">Username</label>
                  <input type="text" className="form-control" id="username" 
                    placeholder="Username" value={loginData.username} 
                    onChange={e => setloginData({...loginData, username: e.target.value})} required />
                </div>
                <div className="form-group">
                  <label htmlFor="password">Password</label>
                  <input type="password" className="form-control" id="password" 
                    placeholder="Password" value={loginData.password}
                    onChange={e => setloginData({...loginData, password: e.target.value})} required />
                </div>
                <button type="submit" className="btn btn-dark">Submit</button>
            </form>
        </section>
    )
}

export default Login