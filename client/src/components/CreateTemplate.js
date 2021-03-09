import { useState } from 'react';

const CreateTemplate = (props) => {
    const [template, setTemplate] = useState({
        name: "",
        textPositions: []
    });

    const handleSubmit = e => {
        e.preventDefault();
    }

    return (
        <div>
            <h2 className="mb-4">Create template</h2>
            <div className="row">
                <div className="col-sm-6 col-12" id="canvas-wrapper">
                <canvas id="template-canvas" className="template-img"></canvas>
                </div>
                <div className="col-sm-6 col-12 new-meme-form-wrapper shadow">
                    <form onSubmit={handleSubmit}>
                        <div className="form-group">
                          <label htmlFor="name">Name</label>
                          <input type="text" className="form-control" id="name" 
                            placeholder="Name" value={template.name} required />
                        </div>
                        <div className="form-group text-center">
                            <button type="submit" className="btn btn-dark">Create</button>
                        </div>                  
                    </form>
                </div>
            </div>            
        </div>
    )
}

export default CreateTemplate