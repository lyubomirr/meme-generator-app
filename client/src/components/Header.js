import { Link } from "react-router-dom";
import { clearUser } from "../auth";
import { useHistory } from "react-router-dom";

const Header = (props) => {

    let history = useHistory()

    const handleLogout = (e) => {
        clearUser()
        props.setUserFunc(null)
        history.push("/login")
    }

    return (
        <nav className="navbar navbar-expand-lg navbar-dark bg-dark header-component">
            <Link className="navbar-brand" to="/">Meme Generator</Link>
            <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                <span className="navbar-toggler-icon"></span>
            </button>
            <div className="collapse navbar-collapse" id="navbarSupportedContent">
                {props.user && 
                    <ul className="navbar-nav mr-auto">
                        <li className="nav-item">
                            <Link to="/" className="nav-link">Browse Memes</Link>
                        </li>
                        <li className="nav-item">
                            <Link to="/templates" className="nav-link">Create Meme</Link>
                        </li>
                    </ul>            
                }

                {props.user && 
                    <ul className="navbar-nav">
                        <li className="nav-item mr-3">
                            <span className="navbar-text">Hello, {props.user.username}</span>
                        </li>
                        <li className="nav-item">
                            <a className="nav-link" href="#" onClick={handleLogout}>Logout</a> 
                        </li>
                    </ul>
                }
            </div>
        </nav>
    )
}

export default Header