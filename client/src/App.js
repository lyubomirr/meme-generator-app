import Header from './components/Header';
import Login from './components/Login';
import Memes from './components/Memes';
import Register from './components/Register';
import Templates from './components/Templates';
import { Switch, Route } from "react-router-dom";
import PrivateRoute from './components/PrivateRoute';
import { useEffect } from 'react'

function App() {
  return (
    <div className="App">
      <Header />
      <div className="container">
        <Switch>
          <PrivateRoute component={Memes} path="/" exact />
          <PrivateRoute component={Templates} path="/templates" exact />
          <Route exact path="/login">          
            <Login />
          </Route>    
          <Route exact path="/register">
            <Register />
          </Route>
        </Switch>
      </div>
    </div>
  );
}

export default App;
