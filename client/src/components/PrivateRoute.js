import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { getCurrentUser } from '../auth';
import NotFound from './NotFound';

const PrivateRoute = ({component: Component, adminOnly, ...rest}) => {
    const currentUser = getCurrentUser()
    let shouldRender = !!currentUser
    if(adminOnly) {
        shouldRender = shouldRender && currentUser.role === "Administrator";
    }

    return (
        <Route {...rest} render={props => (
            shouldRender 
            ? <Component {...props} user={currentUser} />
            : (adminOnly 
                ? <NotFound />
                :<Redirect to="/login" />)
        )} />
    );
};

export default PrivateRoute;