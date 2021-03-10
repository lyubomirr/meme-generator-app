import React, { useState, useEffect } from 'react';
import { fabric } from "fabric";
import {setCanvasBackgroundImage } from "../utils";
import { useToasts } from 'react-toast-notifications';
import ApiFacade from '../api/ApiFacade';
import { useHistory } from "react-router-dom";
import { text } from '@fortawesome/fontawesome-svg-core';

const CreateTemplate = (props) => {
    const history = useHistory();
    const { addToast } = useToasts();

    const [template, setTemplate] = useState({
        name: "",
        textboxes: []
    });
    const [canvas, setCanvas] = useState({});
    const [textboxes, setTextboxes] = useState([]);

    const fileInput = React.createRef();

    useEffect(() => {
        let canvasWidth = document.getElementById("canvas-wrapper").offsetWidth * 0.9;
        setCanvas(() => new fabric.Canvas("template-canvas", {
            width: canvasWidth
        }));
    }, [])

    const handleSubmit = e => {
        e.preventDefault();

        if(textboxes.length === 0) {
            addToast("You should add at least one textbox", {appearance: 'error', autoDismiss: true})
            return;
        }

        for(const t of textboxes) {
            template.textboxes.push({
                topOffset: t.top,
                leftOffset: t.left,
                width: t.width,
                height: t.height
            });
        }

        const blob = fileInput.current.files[0];
        ApiFacade.createTemplate(template, blob)
        .then(t => {
            addToast("Template added successfully!", {appearance: 'success', autoDismiss: true})
            history.push("/templates")                  
        }, err => {
            addToast(err.message, {appearance: 'error', autoDismiss: true})
        })
    }

    const handleFileUpload = e => {
        if (e.target.files.length === 0) {
            canvas.clear();
            return;
        }
        setCanvasBackgroundImage(canvas, URL.createObjectURL(e.target.files[0]));
    }

    const addTextbox = e => {
        e.preventDefault();

        const t = new fabric.Textbox('TEXT ' + (textboxes.length+1),  
        {
            fontFamily: "Impact",
            top: 0,
            left: 0,
            stroke: 'black',
            strokeWidth: 2,
            fill: 'white',
            textAlign: 'center',
            editable: false,
        });

        canvas.add(t);
        setTextboxes(old => [...old, t])
    }

    const copyLastTextbox = e => {
        e.preventDefault();
        if(textboxes.length === 0) {
            addToast("You should add at least one textbox!", {appearance: 'error'})
            return
        }

        const last = textboxes[textboxes.length - 1];
        const t = new fabric.Textbox('TEXT ' + (textboxes.length+1),  
        {
            fontFamily: "Impact",
            top: 0,
            left: 0,
            width: last.width,
            height: last.height,
            stroke: 'black',
            strokeWidth: 2,
            fill: 'white',
            textAlign: 'center',
            editable: false,
        });

        canvas.add(t);
        setTextboxes(old => [...old, t])
    }

    const removeTextbox = e => {
        e.preventDefault();

        const t = textboxes.pop();
        canvas.remove(t);
        setTextboxes(textboxes);
    }

    return (
        <div>
            <h2 className="mb-4">Create template</h2>
            <div className="row">
                <div className="col-sm-6 col-12" id="canvas-wrapper">
                <canvas id="template-canvas" className="template-img"></canvas>
                </div>
                <div className="col-sm-6 col-12 new-template-form-wrapper shadow">
                    <form onSubmit={handleSubmit}>
                        <div className="form-group">
                            <label htmlFor="file">Picture</label>
                            <input type="file" className="form-control" id="file" required 
                            onChange={handleFileUpload} accept="image/jpeg, image/png" ref={fileInput} />
                        </div>
                        <div className="form-group">
                            <label htmlFor="name">Name</label>
                            <input type="text" className="form-control" id="name" placeholder="Name"
                            value={template.name} onChange={e => setTemplate({...template, name: e.target.value})} required />
                        </div>
                        <div className="form-group textbox-btns-wrapper">
                            <button className="btn btn-success textbox-btn" onClick={addTextbox}>Add textbox</button>                   
                            <button className="btn btn-primary textbox-btn" onClick={copyLastTextbox} 
                                disabled={textboxes.length === 0}>Copy last textbox</button>
                            <button className="btn btn-danger textbox-btn" onClick={removeTextbox}>Remove textbox</button>
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