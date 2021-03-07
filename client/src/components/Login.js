import { useState } from 'react';
import ApiFacade from "../api/ApiFacade"
import { useToasts } from 'react-toast-notifications';
import { saveUser } from "../auth"
import { Link, useHistory } from "react-router-dom";


const Login = (props) => {
    let history = useHistory();
    const { addToast } = useToasts();

    const [loginData, setloginData] = useState({
        username: "",
        password: ""
    })
    
    const handleSubmit = (e) => {
        e.preventDefault()
        ApiFacade.login(loginData)
            .then(resp => {
                saveUser(resp)
                props.setUserFunc(resp)
                history.push("/")
            }, err => {
                addToast(err, {appearance: 'error', autoDismiss: true})
            })
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
                <Link to="/Register" className="float-right register-link">Don't have an account? Click here.</Link>
            </form>
        </section>
    )
}

export default Login