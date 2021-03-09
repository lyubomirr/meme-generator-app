import { useState, useEffect } from 'react';
import ApiFacade from '../api/ApiFacade';
import { useToasts } from 'react-toast-notifications';
import Endpoints from '../api/Endpoints';
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faPlusSquare } from "@fortawesome/free-solid-svg-icons";


const Templates = (props) => {
    const { addToast } = useToasts();
    const [templates, setTemplates] = useState([]);
    const isAdmin = props.user.role === "Administrator";

    useEffect(() => {
        ApiFacade.getTemplates()
        .then(templates => {
            setTemplates(templates)
        }, err => {
            addToast(err, {appearance: 'error', autoDismiss: true})
        })
    }, [])
    
    return (
        <div>
            <div className="row">
                <h2 className="mb-4">Choose template</h2>
                {isAdmin && <div class="add-icon"> 
                    <Link to="/create-template">
                        <FontAwesomeIcon icon={faPlusSquare} />
                    </Link>
                </div>}

            </div>
            <div className="row">
            {templates.map((template) => {
                return <div className="col-md-3 col-sm-4 col-6 mb-4" key={template.id}>
                <div className="card">
                    <img className="card-img-top" src={Endpoints.GetTemplateFileUrl(template.id)} alt={template.name} />
                    <div className="card-body">
                        <h5 className="card-title font-weight-bold text-center">{template.name}</h5>
                        <Link to={`/create-meme/${template.id}`} className="btn btn-dark d-block mt-4">Create</Link>
                    </div>
                </div>
            </div>
            })}
            </div>            
        </div>
    )
}

export default Templates