import { useState, useEffect } from 'react';
import ApiFacade from '../api/ApiFacade';
import Endpoints from '../api/Endpoints';
import { useToasts } from 'react-toast-notifications';
import { fabric } from "fabric";

const CreateMeme = (props) => {
    const setBackgroundImage = (canvas) => {
        fabric.Image.fromURL(Endpoints.GetTemplateFileUrl(props.match.params.id), img => {
            const imgRatio = img.width/img.height;
            canvas.setHeight(canvas.width/imgRatio)

            canvas.setBackgroundImage(img, canvas.renderAll.bind(canvas), {
                scaleX: canvas.width / img.width,
                scaleY: canvas.height / img.height
             })
        })
    }

    const initializeCanvas = (textPositions) => {
        let wantedWith = document.getElementById("canvas-wrapper").offsetWidth * 0.9;
        const canvas = new fabric.Canvas("meme-canvas", {
            width: wantedWith
        });
        setBackgroundImage(canvas)
        
        let idx = 1;
        for(const pos of textPositions) {
            const t = new fabric.Textbox('TEXT ' + idx,  
            {
                fontFamily: "Impact",
                top: pos.topOffset,
                left: pos.leftOffset,
                stroke: 'black',
                strokeWidth: 2,
                fill: 'white',
                textAlign: 'center'
            });
            idx++;
            canvas.add(t);
        }
    }
    
    const { addToast } = useToasts();
    const [template, setTemplate] = useState({
        name: "",
        id: 0,
        mimeType: ""
    })

    useEffect(() => {
        ApiFacade.getTemplate(props.match.params.id)
            .then(t => {
                setTemplate(t)
                initializeCanvas(t.textPositions)
            }, err => {
                addToast(err.message, {appearance: 'error', autoDismiss: true})
            })
    }, [])

    const [title, setTitle] = useState("")
    const handleSubmit = (e) => {
        e.preventDefault();

        const canvasEl = document.getElementById("meme-canvas")
        canvasEl.toBlob(blob => {
            const meme = {
                title: title,
                templateId: template.id
            }

            ApiFacade.createMeme(meme, blob)
                .then(m => {
                    addToast("Meme added successfully!", {appearance: 'success'})
                    setTitle("")
                }, err => {
                    addToast(err.message, {appearance: 'error', autoDismiss: true})
                })

        }, template.mimeType)
    };

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