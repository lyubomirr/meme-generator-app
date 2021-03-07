import Header from './components/Header';
import Login from './components/Login';
import Memes from './components/Memes';
import Register from './components/Register';
import Templates from './components/Templates';
import { Switch, Route } from "react-router-dom";
import PrivateRoute from './components/PrivateRoute';
import { ToastProvider } from 'react-toast-notifications';
import CustomToast from "./components/CustomToast"
import { getCurrentUser } from './auth';
import { useState } from 'react';
import CreateMeme from './components/CreateMeme';

function App() {
  const [user, setUser] = useState(getCurrentUser())

  return (
    <ToastProvider placement="top-center" components={{Toast: CustomToast}}>
      <div className="App">
        <Header user={user} setUserFunc={setUser} />
        <div className="container main-container">
          <Switch>
            <PrivateRoute component={Memes} path="/" exact />
            <PrivateRoute component={Templates} path="/templates" exact />
            <PrivateRoute component={CreateMeme} path="/create-meme/:id" />
            <Route exact path="/login">
              <Login setUserFunc={setUser} />
            </Route>
            <Route exact path="/register">
              <Register setUserFunc={setUser} />
            </Route>
          </Switch>
        </div>
      </div>
    </ToastProvider>
  );
}

export default App;
