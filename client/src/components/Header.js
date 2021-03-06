import { Link } from "react-router-dom";
  
const Header = () => {
    return (
        <nav className="navbar navbar-expand-lg navbar-dark bg-dark header-component">
            <a className="navbar-brand" href="#">Meme Generator</a>
            <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                <span className="navbar-toggler-icon"></span>
            </button>
            <div className="collapse navbar-collapse" id="navbarSupportedContent">
                <ul className="navbar-nav mr-auto">
                    <li className="nav-item">
                        <Link to="/" className="nav-link">Memes</Link>
                    </li>
                    <li className="nav-item">
                        <Link to="/templates" className="nav-link">Create Meme</Link>
                    </li>                
                </ul>
            </div>
        </nav>
    )
}

export default Header