import { useState, useEffect } from 'react';
import ApiFacade from '../api/ApiFacade';
import Endpoints from '../api/Endpoints';
import { useToasts } from 'react-toast-notifications';
import { fabric } from "fabric";
import { setCanvasBackgroundImage } from '../utils';
import { useHistory } from "react-router-dom";

const CreateMeme = (props) => {
    const history = useHistory();
    const { addToast } = useToasts();
    const [canvas, setCanvas] = useState({});

    const [title, setTitle] = useState("")
    const [fontSize, setFontSize] = useState(40)
    const [template, setTemplate] = useState({
        name: "",
        id: 0,
        mimeType: ""
    })

    const initializeCanvas = (textboxes) => {
        let wantedWith = document.getElementById("canvas-wrapper").offsetWidth * 0.9;
        const c = new fabric.Canvas("meme-canvas", {width: wantedWith});
        setCanvasBackgroundImage(c, Endpoints.GetTemplateFileUrl(props.match.params.id))
        
        let idx = 1;
        for(const box of textboxes) {
            const t = new fabric.Textbox('TEXT ' + idx,  
            {
                fontFamily: "Impact",
                top: box.topOffset,
                left: box.leftOffset,
                width: box.width,
                height: box.height,
                stroke: 'black',
                strokeWidth: 2,
                fill: 'white',
                textAlign: 'center',
                splitByGrapheme: true,
                onSelect: () => setFontSize(t.fontSize)
            });

            idx++;        
            c.add(t);
        }
        setCanvas(c);
    }

    useEffect(() => {
        ApiFacade.getTemplate(props.match.params.id)
            .then(t => {
                setTemplate(t)
                initializeCanvas(t.textboxes)
            }, err => {
                addToast(err.message, {appearance: 'error', autoDismiss: true})
            })
    }, [])

    const handleSubmit = e => {
        e.preventDefault();

        const canvasEl = document.getElementById("meme-canvas")
        canvasEl.toBlob(blob => {
            const meme = {
                title: title,
                templateId: template.id
            }

            ApiFacade.createMeme(meme, blob)
                .then(m => {
                    addToast("Meme added successfully!", {appearance: 'success', autoDismiss: true})
                    history.push("/")
                }, err => {
                    addToast(err.message, {appearance: 'error', autoDismiss: true})
                })

        }, template.mimeType)
    };

    const changeFontSize = e => {
        if(!canvas.getActiveObject()) {
            return;
        }

        setFontSize(e.target.value);
        canvas.getActiveObject().set('fontSize', e.target.value);
        canvas.renderAll();
    }

    return (
        <div>
            <h2 className="mb-4">Create meme</h2>
            <div className="row">
                <div className="col-sm-6 col-12" id="canvas-wrapper">
                <canvas id="meme-canvas" className="template-img"></canvas>
                </div>
                <div className="col-sm-6 col-12 new-meme-form-wrapper shadow">
                    <form onSubmit={handleSubmit}>
                    <div className="form-group">
                          <label htmlFor="template">Template</label>
                          <input type="text" className="form-control" id="template" value={template.name} disabled/>
                        </div>
                        <div className="form-group">
                          <label htmlFor="title">Title</label>
                          <input type="text" className="form-control" id="title" 
                            placeholder="Title" value={title} 
                            onChange={e => setTitle(e.target.value)} required />
                        </div>
                        <div className="form-group">
                          <label htmlFor="font-size">Font size</label>
                          <input type="number" className="form-control" id="font-size" min="10" max="70"                      
                            value={fontSize} onChange={changeFontSize} required />
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

export default CreateMeme