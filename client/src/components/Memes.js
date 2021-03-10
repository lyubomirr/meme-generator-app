import { useState, useEffect } from 'react';
import ApiFacade from '../api/ApiFacade';
import { useToasts } from 'react-toast-notifications';
import Endpoints from '../api/Endpoints';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";

const Memes = (props) => {
    const isAdmin = props.user.role === "Administrator";
    const { addToast } = useToasts();
    const [memes, setMemes] = useState([])

    useEffect(() => {
        ApiFacade.getMemes()
        .then(memes => {
            for (let m of memes) {
                parseMeme(m)
            }
            setMemes(memes)
        }, err => {
            addToast(err.message, {appearance: 'error', autoDismiss: true})
        })
    }, [])

    const parseMeme = meme => {
        meme.newCommentText = ""
        for(let c of meme.comments) {
            c.createdAt = new Date(c.createdAt).toLocaleString()
        }
    }

    const updateMeme = newMeme => {
        setMemes(
            memes.map(meme => 
                meme.id === newMeme.id 
                ? newMeme
                : meme 
        ))
    }

    const addComment = (e, meme) => {
        e.preventDefault();
        ApiFacade.addComment(meme.id, {
            memeId: meme.id,
            content: meme.newCommentText
        }).then(m => {
            parseMeme(m)
            updateMeme(m)
        }, err => {
            addToast(err.message, {appearance: 'error', autoDismiss: true})
        })
    }

    const deleteComment = (meme, commentId) => {
        ApiFacade.deleteComment(meme.id, commentId).then(m => {
            parseMeme(m)
            updateMeme(m)
            addToast("Sucesfully deleted comment.", {appearance: 'success', autoDismiss: true})
        }, err => {
            addToast(err.message, {appearance: 'error', autoDismiss: true})
        })
    }

    const deleteMeme = meme => {
        ApiFacade.deleteMeme(meme.id).then(() => {
            removeMemeFromState(meme)
            addToast("Sucesfully deleted meme.", {appearance: 'success', autoDismiss: true})
        }, err => {
            addToast(err.message, {appearance: 'error', autoDismiss: true})
        })
    }

    const removeMemeFromState = meme => {
        const idx = memes.indexOf(meme)
        if (idx === -1) {
            return
        }
        memes.splice(idx, 1)
        setMemes(memes)
    }

    return (
        <div>
            <div className="row">
                <h2 className="mb-4">Memes</h2>
            </div>
            <div className="row">
            {memes.length === 0 &&
                <div className="col-12 text-center mt-5">
                    <h3 className="text-center">No memes added :(</h3>
                </div> 
            }
            {memes.map(meme => {
                return (
                    <div className="col-sm-10 col-12 mb-4" key={meme.id}>
                        <div className="card">
                            <div className="row">
                                <div className="col-sm-6 col-12 meme-img-container">
                                    <img className="card-img-top" src={Endpoints.GetMemeFileUrl(meme.id)} alt={meme.title} />
                                    <div className="card-body text-center">
                                        <h5 className="card-title meme-title">{meme.title}</h5>
                                    </div> 
                                </div>
                                <div className="col-sm-6 col-12 comments-container">
                                    {(meme.author.username === props.user.username || isAdmin) && 
                                        <span className="delete-meme-icon" onClick={() => deleteMeme(meme)}>
                                            <FontAwesomeIcon icon={faTrashAlt} />
                                        </span>                                                                
                                    }
    
                                    <div className="d-flex justify-content-center pt-3 pb-2"> 
                                        <form onSubmit={(e) => addComment(e, meme)}>
                                            <input type="text" name="text" placeholder="+ Add comment" className="form-control addtxt" required
                                            value={meme.newCommentText} onChange={e => updateMeme({...meme, newCommentText : e.target.value})}/>
                                        </form>
                                    </div>
                                    {meme.comments.map(comment => {
                                        return (
                                        <div className="d-flex justify-content-center py-2" key={comment.id}> 
                                            <div className="second py-2 px-2">
                                                <div className="d-flex justify-content-between py-1 pt-2">
                                                    <div><span className="text2">{comment.content}</span></div>
                                                    {(comment.author.username === props.user.username || isAdmin) && 
                                                        <div onClick={() => deleteComment(meme, comment.id)}>
                                                            <span className="text3 delete-comment-icon">
                                                                <FontAwesomeIcon icon={faTrashAlt} />
                                                            </span>
                                                        </div>                                                                                                        
                                                    }
                                                </div>
                                                <div className="d-flex justify-content-between py-1 pt-2">
                                                    <div><span className="text2">{comment.author.username}</span></div>
                                                    <div><span className="text3">{comment.createdAt}</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        )
                                    })}
                                </div>
                            </div>
                        </div>
                    </div>                
                )
            })}
            </div>
        </div>
    )
}

export default Memes