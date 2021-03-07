import { useState } from 'react';
import ApiFacade from "../api/ApiFacade"
import { useToasts } from 'react-toast-notifications';
import { saveUser } from "../auth"
import { Link, useHistory } from "react-router-dom";

const Register = (props) => {
    let history = useHistory();
    const { addToast } = useToasts();

    const [registerData, setRegisterData] = useState({
        username: "",
        password: "",
        confirmPassword: ""
    })
    
    const handleSubmit = (e) => {
        e.preventDefault()
        ApiFacade.register(registerData)
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
                    placeholder="Username" value={registerData.username} 
                    onChange={e => setRegisterData({...registerData, username: e.target.value})} required />
                </div>
                <div className="form-group">
                  <label htmlFor="password">Password</label>
                  <input type="password" className="form-control" id="password" 
                    placeholder="Password" value={registerData.password}
                    onChange={e => setRegisterData({...registerData, password: e.target.value})} required />
                </div>
                <div className="form-group">
                  <label htmlFor="password">Confirm Password</label>
                  <input type="password" className="form-control" id="password" 
                    placeholder="Password" value={registerData.confirmPassword}
                    onChange={e => setRegisterData({...registerData, confirmPassword: e.target.value})} required />
                </div>
                <button type="submit" className="btn btn-dark">Submit</button>
            </form>
        </section>
    )
}

export default Register