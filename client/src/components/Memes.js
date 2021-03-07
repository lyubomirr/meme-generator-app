import { useState, useEffect } from 'react';
import ApiFacade from '../api/ApiFacade';
import { useToasts } from 'react-toast-notifications';

const Memes = () => {
    const { addToast } = useToasts();
    const [memes, setMemes] = useState([])

    useEffect(() => {
        ApiFacade.getMemes()
        .then(memes => {
            setMemes(memes)
        }, err => {
            addToast(err, {appearance: 'error', autoDismiss: true})
        })
    }, [])



    return (
        <h1>Memes</h1>
    )
}

export default Memes